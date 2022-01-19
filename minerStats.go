package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"gorm.io/gorm"

	"github.com/GregoryUnderscore/Mining-Automation-Shared/database"
	. "github.com/GregoryUnderscore/Mining-Automation-Shared/models"
)

// ====================================
// Configuration File (ZergPoolData.hcl)
// ====================================
type Config struct {
	// Database Connectivity
	Host     string `hcl:"host"`     // The server hosting the database
	Port     string `hcl:"port"`     // The port of the database server
	Database string `hcl:"database"` // The database name
	User     string `hcl:"user"`     // The user to use for login to the database server
	Password string `hcl:"password"` // The user's password for login
	TimeZone string `hcl:"timezone"` // The time zone where the program is run

	// Miner Specific Settings
	MinerName string `hcl:"minerName"` // The name of the mining hardware
	MinerPath string `hcl:"minerPath"` // The path to the mining executable
	// This is used to find the hash rate in the mining program's screen output (which is saved to a file).
	StatSearchPhrase string `hcl:"statSearchPhrase"`
	// 1=Yes, 0=No - Whether to include GPU(s) when calculating statistics.
	AllowGPU uint8 `hcl:"allowGPU"`
	// The amount of time to wait before checking output for statistics, in seconds.
	// It can be helpful to give the program a few minutes sometimes, as it often calculates an average
	// hash rate instead of a current hash rate.
	StatWaitTime uint16 `hcl:"statWaitTime"`
}

func main() {
	const configFileName = "MinerStats.hcl"
	var config Config
	// Used to pull down all the mining software and compare it against the planned mining program in
	// the config.
	var minerSoftware []MinerSoftware
	// The matched mining software
	var minerProggy MinerSoftware
	// This will have all the algos support by the software that have a pool.
	var minerSoftwareAlgos []MinerSoftwareAlgos

	// Grab the configuration details for the database connection. These are stored in ZergPoolData.hcl.
	err := hclsimple.DecodeFile(configFileName, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load config file "+configFileName+".\n", err)
	}

	// Connect to the database and create/validate the schema.
	db := database.Connect(config.Host, config.Port, config.Database, config.User, config.Password,
		config.TimeZone)
	database.VerifyAndUpdateSchema(db)

	// Open the new database transaction and get all the coins from CoinGecko along with the BTC price.
	tx := db.Begin()

	defer func() { // Ensure transaction rollback on panic
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	log.Println("Creating records required for calculations...")
	minerID := verifyMiner(tx, config.MinerName)
	verifyMinerSoftware(tx)

	// Determine the mining software.
	result := tx.Find(&minerSoftware)
	if result.Error != nil {
		log.Fatalf("Failure to find miner software in the database. This probably did not store "+
			"properly. Please try re-running the program.\n", result.Error)
	}
	// Cycle over all the software in the database and check for a match on the config setting.
	// The mining software must be identified to store statistics.
	for _, software := range minerSoftware {
		// The executable prefix should be in the miner path. NOTE: This can be a problem for
		// mining programs that share a prefix! TODO Create a config parameter to override this.
		if strings.Contains(config.MinerPath, software.ExecutablePrefix) {
			minerProggy = software
			break
		}
	}
	if (MinerSoftware{}) == minerProggy {
		log.Fatalf("Failure to locate the mining program in the config file. All miner software " +
			"executable prefixes were checked, and none of them matched. Please update the " +
			"config file to a matching prefix or update miner_software table.")
	}

	// Get all the algorithms for the miner that have a pool.
	tx.Joins("INNER JOIN pools on pools.algorithm_id = miner_software_algos.algorithm_id").
		Distinct().
		Where("miner_software_id = ?", minerProggy.ID).
		Find(&minerSoftwareAlgos)

	log.Println("Beginning mining statistic calculations...")

	// Cycle over the miner software algos, execute the miner with the algo, and grab statistics via file
	// output. Once statistics are obtained, store them into the database.
	for _, algo := range minerSoftwareAlgos {
		log.Println("Starting statistics for " + algo.Name + "...")
		// Create the core parameter structure for the miner software.
		// This includes the algorithm parameter requirements and any other requirements for
		// benchmarking an algorithm.
		params := strings.Split(minerProggy.OtherParams, " ")
		// Create the full parameter list
		params = append([]string{minerProggy.Name, minerProggy.AlgoParam, algo.Name},
			params...)
		// Some algorithms have parameters specific to them.
		if len(algo.ExtraParams) > 0 {
			extraParams := strings.Split(algo.ExtraParams, " ")
			params = append(params, extraParams...)
		}

		// Output should go to a file in the run folder.
		outputFile := minerProggy.Name + "-" + algo.Name + "-" + time.Now().Format("20060102150405")
		stdout, _ := os.Create(outputFile) // Used to pipe the output from stdout to the file

		// Open the miner program in a child process.
		attr := &os.ProcAttr{
			"",
			nil,
			[]*os.File{os.Stdin, stdout, stdout},
			&syscall.SysProcAttr{},
		}
		proc, error := os.StartProcess(config.MinerPath, params, attr)
		if error != nil {
			log.Fatalf("Unable to start mining software.\n", error)
		}

		// Give the process enough time to produce statistics.
		log.Println("Waiting for statistics...")
		time.Sleep(time.Duration(config.StatWaitTime) * time.Second)

		// The wait time has finished. Force the process to stop.
		error = proc.Kill()

		// Cycle over the file and look for the 5th match on the search phrase.
		// This ensures the first statistic output is not used (usually invalid).
		stdout.Seek(0, 0) // Start of file
		scanner := bufio.NewScanner(stdout)
		linesFound := 0 // This must get to 5
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, config.StatSearchPhrase) {
				linesFound++
				if linesFound < 5 { // TODO Make this configurable.
					continue
				}
				// Process the hash statistic and store into the database.
				processHashLine(tx, algo, minerID, line, config.StatSearchPhrase)
			}
		}
		stdout.Close()
	}

	err = tx.Commit().Error // Finalize data storage
	if err != nil {
		log.Fatalf("Issue committing changes.\n", err)
	}
	log.Println("Statistics stored.\nOperations complete.\n")
}

// Get the megahash factor for a string. The string will be converted to lowercase.
// @param unit - The string unit to convert. E.g. "mh/s" = 1
// @returns The numeric equivalent
func getMhFactor(unit string) float64 {
	lowerCaseUnit := strings.ToLower(unit)

	switch {
	case strings.Contains(lowerCaseUnit, "ph/s"):
		return 1000000000
	case strings.Contains(lowerCaseUnit, "th/s"):
		return 1000000
	case strings.Contains(lowerCaseUnit, "gh/s"):
		return 1000
	case strings.Contains(lowerCaseUnit, "mh/s"):
		return 1
	case strings.Contains(lowerCaseUnit, "kh/s"):
		return 0.001
	case strings.Contains(lowerCaseUnit, "h/s"):
		return 0.000001
	default:
		log.Fatalf("Invalid mega-hash factor: \"" + lowerCaseUnit + "\"")
	}
	return -1 // This should never get hit, but if it does, it is an error.
}

// Pass in a line of mining output to parse out the actual work done and store the statistics in the database.
// @params tx - The active database session
// @params algo - The algorithm in use by the miner
// @params minerID - The miner's ID in the database
// @params line - The line to process
// @params searchPhrase - The phrase that indicates where the work statistic is in the line
func processHashLine(tx *gorm.DB, algo MinerSoftwareAlgos, minerID uint64, line string, searchPhrase string) {
	// Split into tokens for matching on search phrase.
	pieces := strings.Split(line, " ")
	phraseFound := false

	// Cycle over all the tokens and look for a search phrase match. Then, parse out the work calculated
	// and store into the database for future reporting.
	for index, piece := range pieces {
		if !phraseFound && piece == searchPhrase {
			phraseFound = true
			workPerSecondIndex := index + 1
			unitIndex := index + 2
			// The format for the output must be NNNNN UNIT.
			// Examples:
			//    45.9 h/s
			//    4000 Kh/s
			// TODO This may require configurability if some miners
			// output the unit without a space, e.g. 45.9h/s.
			if len(pieces) > unitIndex {
				var stat MinerStats
				stat.WorkPerSecond, _ = strconv.ParseFloat(
					pieces[workPerSecondIndex], 64)
				stat.MhFactor = getMhFactor(pieces[unitIndex])
				log.Println("Calculated " +
					strconv.FormatFloat(stat.WorkPerSecond,
						'f', 3, 64) +
					" " + pieces[unitIndex] + " (" +
					strconv.FormatFloat(stat.MhFactor,
						'f', 6, 64) + "). Storing...")

				stat.AlgorithmID = algo.AlgorithmID
				stat.Instant = time.Now()
				stat.MinerID = minerID
				stat.MinerSoftwareID = algo.MinerSoftwareID
				result := tx.Create(&stat)
				if result.Error != nil {

					log.Fatalf("Issue creating miner statistic for"+
						algo.Name+".\n", result.Error)
				}
				// Stop processing
				break
			}
		}
	}
}

// Verify the miner exists in the database. If not, create it.
// @param tx - The active database session
// @param minerName - The name of the mining hardware
// @returns The ID associated with the miner.
func verifyMiner(tx *gorm.DB, minerName string) uint64 {
	var miner Miner
	result := tx.Where("name = ?", minerName).Limit(1).Find(&miner)
	if result.RowsAffected == 0 {
		log.Println("Creating miner...")
		miner.Name = minerName
		result = tx.Create(&miner)
		if result.Error != nil {
			log.Fatalf("Issue creating miner.\n", result.Error)
		}
	} else if result.Error != nil {
		log.Fatalf("Unknown issue storing miner.\n", result.Error)
	} else {
		log.Println("Found existing miner.")
	}
	return miner.ID
}

// Verify mining software exists and if not add it.
// This is mining software that is supported out-of-the-box.
// @param tx - The active database session
func verifyMinerSoftware(tx *gorm.DB) {
	verifyCPUMinerOpt(tx) // cpuminer-opt
}

// CPU-Miner-Opt is supported out of the box. Verify it is in the system, and if not, add it.
// This handles the software record and the algorithm maps.
// @param tx - The active database session
func verifyCPUMinerOpt(tx *gorm.DB) {
	var minerSoftware MinerSoftware
	softwareName := "cpuminer-opt"
	website := "https://github.com/JayDDee/cpuminer-opt/releases"
	prefix := "cpuminer-"
	algoParam := "--algo"
	otherParams := "--benchmark" // Runs in benchmark mode without URL target

	// Check if it exists, and if not, create.
	result := tx.Where("name = ?", softwareName).Limit(1).Find(&minerSoftware)
	if (MinerSoftware{}) == minerSoftware {
		log.Println("Creating miner software record for cpuminer-opt...")
		minerSoftware.Name = softwareName
		minerSoftware.Website = website
		minerSoftware.ExecutablePrefix = prefix
		minerSoftware.AlgoParam = algoParam
		minerSoftware.OtherParams = otherParams
		result = tx.Create(&minerSoftware)
		if result.Error != nil {
			log.Fatalf("Issue creating miner software.\n", result.Error)
		}

		// Create the map to Zergpool algos.
		mapCPUMinerOptToZergAlgos(tx, minerSoftware.ID)

	} else if result.Error != nil {
		log.Fatalf("Unknown issue storing miner software.\n", result.Error)
	}
}

// This creates the bridge records in MinerSoftwareAlgos. These map the algorithms supported
// by cpuminer-opt to the algorithms stored from Zergpool (and potentially others).
// @param tx - The active database session
// @param minerSoftwareID - The ID for the mining software.
func mapCPUMinerOptToZergAlgos(tx *gorm.DB, minerSoftwareID uint64) {
	// The key is the name in cpuminer-opt and the value is the name from the pool.
	// If they are equivalent (or missing from the pool), the value is an empty string.
	supportedAlgos := map[string]string{
		"allium":      "", //        Garlicoin (GRLC)
		"anime":       "", //         Animecoin (ANI)
		"argon2":      "", //        Argon2 Coin (AR2)
		"argon2d250":  "",
		"argon2d500":  "argon2d-dyn", //    argon2d-dyn, Dynamic (DYN)
		"argon2d4096": "",            //   argon2d-uis, Unitus (UIS)
		"axiom":       "",            //         Shabal-256 MemoHash
		"blake":       "",            //         blake256r14 (SFR)
		"blake2b":     "",            //       Blake2b 256
		"blake2s":     "",            //       Blake-2 S
		"blakecoin":   "",            //     blake256r8
		"bmw":         "",            //           BMW 256
		"bmw512":      "",            //        BMW 512
		"c11":         "",            //           Chaincoin
		"decred":      "",            //        Blake256r14dcr
		"deep":        "",            //          Deepcoin (DCN)
		"dmd-gr":      "",            //        Diamond
		"groestl":     "",            //       Groestl coin
		"hex":         "",            //           x16r-hex
		"hmq1725":     "",            //       Espers
		"hodl":        "",            //          Hodlcoin
		"jha":         "",            //           jackppot (Jackpotcoin)
		"keccak":      "",            //        Maxcoin
		"keccakc":     "",            //       Creative Coin
		"lbry":        "",            //          LBC, LBRY Credits
		"lyra2h":      "",            //        Hppcoin
		"lyra2re":     "",            //       lyra2
		"lyra2rev2":   "lyra2v2",     //     lyrav2
		"lyra2rev3":   "",            //     lyrav2v3
		"lyra2z":      "",
		"lyra2z330":   "",    //     Lyra2 330 rows
		"m7m":         "",    //           Magi (XMG)
		"myr-gr":      "",    //        Myriad-Groestl
		"minotaur":    "",    //      Ringcoin (RNG)
		"neoscrypt":   "",    //     NeoScrypt(128, 2, 1)
		"nist5":       "",    //         Nist5
		"pentablake":  "",    //    5 x blake512
		"phi1612":     "phi", //       phi
		"phi2":        "",
		"polytimos":   "",
		"power2b":     "", //       MicroBitcoin (MBC)
		"quark":       "", //         Quark
		"qubit":       "", //         Qubit
		"scrypt":      "", //        scrypt(1024, 1, 1) (default)
		// scryptn2 is handled separately from the others due to the requirement of a special parameter.
		// --param-n 1048576
		//"scrypt":     "scryptn2",
		"scrypt:N":     "", //      scrypt(N, 1, 1)
		"sha256d":      "", //       Double SHA-256
		"sha256q":      "", //       Quad SHA-256, Pyrite (PYE)
		"sha256t":      "", //       Triple SHA-256, Onecoin (OC)
		"sha3d":        "", //         Double Keccak256 (BSHA3)
		"shavite3":     "", //      Shavite3
		"skein":        "", //         Skein+Sha (Skeincoin)
		"skein2":       "", //        Double Skein (Woodcoin)
		"skunk":        "", //         Signatum (SIGT)
		"sonoa":        "", //         Sono
		"timetravel":   "", //    timeravel8, Machinecoin (MAC)
		"timetravel10": "", //  Bitcore (BTX)
		"tribus":       "", //        Denarius (DNR)
		"vanilla":      "", //       blake256r8vnl (VCash)
		"veltor":       "",
		"verthash":     "",
		"whirlpool":    "",
		"whirlpoolx":   "",
		"x11":          "",    //           Dash
		"x11evo":       "",    //        Revolvercoin (XRE)
		"x11gost":      "sib", //       sib (SibCoin)
		"x12":          "",    //           Galaxie Cash (GCH)
		"x13":          "",    //           X13
		"x13bcd":       "",    //        bcd
		"x13sm3":       "",    //        hsr (Hshare)
		"x14":          "",    //           X14
		"x15":          "",    //           X15
		"x16r":         "",
		"x16rv2":       "",
		"x16rt":        "", //         Gincoin (GIN)
		"x16rt-veil":   "", //    Veil (VEIL)
		"x16s":         "",
		"x17":          "",
		"x21s":         "",
		"x22i":         "",
		"x25x":         "",
		"xevan":        "",            //         Bitsend (BSD)
		"yescrypt":     "",            //      Globalboost-Y (BSTY)
		"yescryptr8":   "",            //    BitZeny (ZNY)
		"yescryptr8g":  "",            //   Koto (KOTO)
		"yescryptr16":  "yescryptR16", //   Eli
		"yescryptr32":  "yescryptR32", //   WAVI
		"yespower":     "",            //      Cryply
		"yespowerr16":  "yespowerR16", //   Yenten (YTN)
		"yespower-b2b": "",            //  generic yespower + blake2b
		"zr5":          "",            //           Ziftr
	}

	// Verify there are records in Algorithms.
	// If not, error out and let the user know they need to run the Zergpool statistics program
	// first (or another pool provider statistics program).
	var algo Algorithm
	result := tx.First(&algo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("No algorithms exist in the algorithm table. Run a pool statistics collection " +
			"before calculating miner statistics. Only pool algorithms are calculated to avoid " +
			"wasting time calculating useless statistics. To run a pool statistics collection, " +
			"see Mining-Automation-Zergpool.com as an example.")
	}

	// Handle scryptn2 first, as it requires special handling.
	var minerAlgo MinerSoftwareAlgos
	algo = Algorithm{} // Reset
	result = tx.Where("LOWER(name) = ?", "scryptn2").Limit(1).Find(&algo)
	if (Algorithm{}) == algo || result.Error != nil {
		log.Println("Skipping scryptn2 due to unexpected issue.\n", result.Error)
	} else {
		// Create the MinerSoftwareAlgos record with the special parameter.
		minerAlgo.MinerSoftwareID = minerSoftwareID
		minerAlgo.AlgorithmID = algo.ID
		minerAlgo.Name = "scrypt" // Uses scrypt with special parameter below.
		minerAlgo.ExtraParams = "--param-n 1048576"
		result = tx.Create(&minerAlgo)
		if result.Error != nil {
			log.Fatalf("Issue creating miner software algo map for "+minerAlgo.Name+".\n",
				result.Error)
		}
	}

	// Cycle over the map. Check if the algorithm exists in the algorithms table.
	// Do that by using the value. If the value is blank, use the key.
	for softwareAlgo, poolAlgo := range supportedAlgos {
		algoToFind := poolAlgo
		if poolAlgo == "" {
			algoToFind = softwareAlgo
		}
		algo = Algorithm{} // Reset to avoid any collisions.
		result = tx.Where("LOWER(name) = ?", strings.ToLower(algoToFind)).Limit(1).Find(&algo)
		// Skip anything not in the database as it likely is not in use by a pool.
		if algo.ID > 0 {
			// Reset for creation.
			minerAlgo = MinerSoftwareAlgos{}
			minerAlgo.MinerSoftwareID = minerSoftwareID
			minerAlgo.AlgorithmID = algo.ID
			minerAlgo.Name = softwareAlgo // This is what the software will require in params.
			result = tx.Create(&minerAlgo)
			if result.Error != nil {
				log.Fatalf("Issue creating miner software algo map for "+minerAlgo.Name+".\n",
					result.Error)
			}
		}
	}
}

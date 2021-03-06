package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"gorm.io/gorm"

	. "github.com/GregoryUnderscore/Mining-Automation-Shared/database"
	. "github.com/GregoryUnderscore/Mining-Automation-Shared/models"
	. "github.com/GregoryUnderscore/Mining-Automation-Shared/utils/pools"
)

// ====================================
// Configuration File (MinerStats.hcl)
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
	MinerName string `hcl:"minerName"`       // The name of the mining hardware
	Wallet    string `hcl:"wallet,optional"` // A valid wallet (some software requires a pool connection)

	// Mining software configuration
	Software []SoftwareConfig `hcl:"software,block"`
}

// Represents mining software
type SoftwareConfig struct {
	Name string `hcl:"name,label"` // Used to determine if the miner is already in the system
	// Optional - may be used in the future to handle new releases
	ReleaseWebsite string `hcl:"releaseWebsite"`
	FilePath       string `hcl:"filePath"`  // The path to the executable
	AlgoParam      string `hcl:"algoParam"` // The parameter used for algorithms
	// Whether the software should connect for the actual assessment of statistics. Sometimes required.
	ConnectForAssessment bool   `hcl:"connectForAssessment"`
	PoolParam            string `hcl:"poolParam,optional"`      // Used for pool URL specification
	PasswordParam        string `hcl:"passwordParam,optional"`  // Necessary for pool operations sometimes.
	WalletParam          string `hcl:"walletParam,optional"`    // Passes a wallet to the connected pool.
	FileParam            string `hcl:"fileParam,optional"`      // Use if software can log to a file.
	BenchmarkParam       string `hcl:"benchmarkParam,optional"` // Some support a benchmark parameter
	OtherParams          string `hcl:"otherParams"`             // Other important parameters
	// This is used to find the hash rate in the mining program's screen output (which is saved to file).
	StatSearchPhrase string `hcl:"statSearchPhrase"`
	// The amount of time to wait before checking output for statistics, in seconds.
	// It can be helpful to give the program a few minutes sometimes, as it often calculates an average
	// hash rate instead of a current hash rate.
	StatWaitTime uint16 `hcl:"statWaitTime"`
	// The number of hashrate lines to skip. Can be useful if the software outputs low hashrate initially.
	// 1 will skip 1 line of hashrate output.
	SkipLines uint8 `hcl:"skipLines"`
	// The number of string tokens to skip to get to the hash unit. For instance, if the software outputs:
	// 1234 H/s  - this is set to 0, nothing to skip
	// 1234 1000 444 H/s - this is set to 2, as there are 2 numbers to skip between 1234 and H/s.
	// This setting can vary depending on the number preferred.
	SkipTokens  uint8        `hcl:"skipTokens,optional"`
	AlgoConfigs []AlgoConfig `hcl:"algo,block"`
}

// Maps the algo name for the miner to the algo name for a pool provider.
type AlgoConfig struct {
	MinerName   string `hcl:"minerName,label"`      // The miner's name for the algo
	PoolName    string `hcl:"poolName"`             // The pool's name for the algo, in the algorithm table.
	ExtraParams string `hcl:"extraParams,optional"` // Any additional parameters unique to the algorithm
}

func main() {
	const configFileName = "MinerStats.hcl"
	var config Config

	// Grab the configuration details for the database connection. These are stored in ZergPoolData.hcl.
	err := hclsimple.DecodeFile(configFileName, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load config file "+configFileName+".\n", err)
	}

	// Connect to the database and create/validate the schema.
	db := Connect(config.Host, config.Port, config.Database, config.User, config.Password,
		config.TimeZone)
	VerifyAndUpdateSchema(db)

	// Open the new database transaction.
	tx := db.Begin()

	defer func() { // Ensure transaction rollback on panic
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	log.Println("Creating records required for calculations...")
	minerID := VerifyMiner(tx, config.MinerName)
	err = tx.Commit().Error // Commit changes to the database
	if err != nil {
		log.Fatalf("Issue committing changes.\n", err)
	}
	// Cycle over all the software in the config and check for a match in the database. Create/update
	// accordingly. This handles all the mapping of software to pool algos.
	// If a file path was specified for the miner, run calculations and store them into the database.
	for _, minerSoft := range config.Software {
		// The matched mining software
		var minerProggy MinerSoftware
		// This will have all the algos support by the software that have a pool.
		var minerSoftwareAlgos []MinerSoftwareAlgos
		tx = db.Begin() // Start anew
		// Verify the miner software exists and if not create it.
		minerProggy = verifyMinerSoftware(tx, minerSoft)
		err = tx.Commit().Error // Commit changes to the database
		if err != nil {
			log.Fatalf("Issue committing changes.\n", err)
		}

		if len(minerSoft.FilePath) > 0 { // If a file path was specified, run calculations.
			if (MinerSoftware{}) == minerProggy {
				log.Fatalf("Unexpected failure to locate the mining program in the database.")
			}
			// Store the file path for the mining software on the miner.
			storeMinerSoftwareDetails(db, minerID, minerProggy.ID, minerSoft.FilePath)
			// Get all the algorithms for the miner that have a pool.
			db.Order("name").
				Joins("INNER JOIN pools on pools.algorithm_id = "+
					"miner_software_algos.algorithm_id").
				Distinct().
				Where("miner_software_id = ?", minerProggy.ID).
				Find(&minerSoftwareAlgos)

			log.Println("Beginning mining statistic calculations...")

			// Cycle over the miner software algos, execute the miner with the algo, and
			// grab statistics via file output. Once statistics are obtained, store them into the
			// database.
			for _, algo := range minerSoftwareAlgos {
				log.Println("Starting statistics for " + algo.Name + "...")
				// Create the core parameter structure for the miner software.
				// This includes the algorithm parameter requirements and any other
				// requirements for benchmarking an algorithm.
				// Create the full parameter list
				params := []string{minerProggy.Name, minerProggy.AlgoParam, algo.Name}
				// Add in the benchmarking parameters, if any.
				if len(minerProggy.BenchmarkParams) > 0 {
					benchmarkParams := strings.Split(minerProggy.BenchmarkParams, " ")
					params = append(params, benchmarkParams...)
				}
				// Process any additional parameters in the catch-all other parameters.
				if len(minerProggy.OtherParams) > 0 {
					otherParams := strings.Split(minerProggy.OtherParams, " ")
					params = append(params, otherParams...)
				}
				// A pool connection is required. Generate a URL and append to params.
				if minerSoft.ConnectForAssessment {
					params = append(params, minerSoft.PoolParam)
					url := GeneratePoolURL(db, algo.AlgorithmID)
					params = append(params, url)
					// If connecting to a pool, a wallet is sometimes required.
					if len(minerSoft.WalletParam) > 0 {
						params = append(params, minerSoft.WalletParam, config.Wallet)
					}
				}
				// Some algorithms have parameters specific to them.
				if len(algo.ExtraParams) > 0 {
					extraParams := strings.Split(algo.ExtraParams, " ")
					params = append(params, extraParams...)
				}

				var stdout *os.File
				// Output should go to a file in the run folder.
				// XMRRig has forward slashes in the algo names that must be removed.
				outputFile := minerProggy.Name + "-" +
					strings.ReplaceAll(algo.Name, "/", "") + "-" +
					time.Now().Format("20060102150405.txt")
				output := []*os.File{os.Stdin, os.Stdout, os.Stderr}
				// If software does not support logging, try to force its output to a
				// file via stdout.
				if len(minerSoft.FileParam) == 0 {
					// Used to pipe the output from stdout to the file
					stdout, _ = os.Create(outputFile)
					output = []*os.File{os.Stdin, stdout, os.Stderr}
				} else { // Software supports logging. Use its native ability.
					params = append(params, minerSoft.FileParam, outputFile)
				}

				// Open the miner program in a child process.
				attr := &os.ProcAttr{
					"",
					nil,
					output,
					&syscall.SysProcAttr{},
				}
				proc, error := os.StartProcess(minerSoft.FilePath, params, attr)
				if error != nil {
					log.Fatalf("Unable to start mining software.\n", error)
				}

				// Give the process enough time to produce statistics.
				log.Println("Waiting for statistics (" +
					fmt.Sprint(minerSoft.StatWaitTime) + "s)...")
				time.Sleep(time.Duration(minerSoft.StatWaitTime) * time.Second)

				// The wait time has finished. Force the process to stop.
				error = proc.Kill()
				proc.Wait() // Wait for everything to wrap up.

				if len(minerSoft.FileParam) > 0 {
					stdout, error = os.Open(outputFile)
					if error != nil {
						log.Fatalf("Unable to open log file for stats.\n", error)
					}
				}

				// Cycle over the file and look for the 5th match on the search phrase.
				// This ensures the first statistic output is not used (usually invalid).
				stdout.Seek(0, 0) // Start of file
				scanner := bufio.NewScanner(stdout)
				linesFound := 0 // This must get to 5
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, minerSoft.StatSearchPhrase) {
						linesFound++
						// Skip hashrate output according to the settings.
						if linesFound < int(minerSoft.SkipLines) {
							continue
						}
						// Process the hash statistic and store into the database.
						tx = db.Begin()
						processHashLine(tx, algo, minerID, line,
							minerSoft.StatSearchPhrase, minerSoft.SkipTokens)
						err = tx.Commit().Error // Commit changes to the database
						if err != nil {
							log.Fatalf("Issue committing changes.\n", err)
						}
					}
				}
				stdout.Close()
			}
		}
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
// @params unitSkipCount - In some cases, the unit is not adjacent and tokens must be skipped to get it.
func processHashLine(tx *gorm.DB, algo MinerSoftwareAlgos, minerID uint64, line string, searchPhrase string,
	unitSkipCount uint8) {
	// Split into tokens for matching on search phrase.
	pieces := strings.Split(line, " ")
	phraseFound := false

	// Cycle over all the tokens and look for a search phrase match. Then, parse out the work calculated
	// and store into the database for future reporting.
	for index, piece := range pieces {
		if !phraseFound && piece == searchPhrase {
			phraseFound = true
			workPerSecondIndex := index + 1
			unitIndex := index + 2 + int(unitSkipCount)
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

					log.Fatalf("Issue creating miner statistic for "+
						algo.Name+".\n", result.Error)
				}
				// Stop processing
				break
			}
		}
	}
}

// Store into the database the file path for a miner/software.
// @param db - The active database connection
// @param minerID - The database ID for the miner
// @param minerSoftwareID - The database ID for the mining software
// @param filePath - The file path to the mining software on the miner
func storeMinerSoftwareDetails(db *gorm.DB, minerID uint64, minerSoftwareID uint64, filePath string) {
	var minerSoftwareDetails MinerMinerSoftware
	tx := db.Begin() // Start anew
	tx.Where("miner_id = ? AND miner_software_id = ?", minerID, minerSoftwareID).
		Find(&minerSoftwareDetails)
	if (MinerMinerSoftware{}) == minerSoftwareDetails {
		minerSoftwareDetails.FilePath = filePath
		minerSoftwareDetails.MinerID = minerID
		minerSoftwareDetails.MinerSoftwareID = minerSoftwareID
		result := tx.Create(&minerSoftwareDetails)
		if result.Error != nil {
			log.Fatalf("Issue storing file path for miner software.\n",
				result.Error)
		}
	} else { // Save in case the file path has changed.
		minerSoftwareDetails.FilePath = filePath
		result := tx.Save(&minerSoftwareDetails)
		if result.Error != nil {
			log.Fatalf("Issue updating file path for miner software.\n",
				result.Error)
		}
	}
	err := tx.Commit().Error // Commit changes to the database
	if err != nil {
		log.Fatalf("Issue committing changes.\n", err)
	}
}

// Verify mining software exists and if not add it.
// This is mining software that is supported out-of-the-box.
// @param tx - The active database session
// @param software - The configuration for a mining program
// @returns - The miner software record from the database
func verifyMinerSoftware(tx *gorm.DB, software SoftwareConfig) MinerSoftware {
	var minerSoftware MinerSoftware

	// Check if it exists, and if not, create.
	result := tx.Where("name = ?", software.Name).Limit(1).Find(&minerSoftware)
	if (MinerSoftware{}) == minerSoftware {
		log.Println("Creating miner software record for cpuminer-opt...")
		minerSoftware.Name = software.Name
		minerSoftware.Website = software.ReleaseWebsite
		minerSoftware.AlgoParam = software.AlgoParam
		minerSoftware.PoolParam = software.PoolParam
		minerSoftware.PasswordParam = software.PasswordParam
		minerSoftware.WalletParam = software.WalletParam
		minerSoftware.FileParam = software.FileParam
		minerSoftware.BenchmarkParams = software.BenchmarkParam
		minerSoftware.OtherParams = software.OtherParams
		minerSoftware.SkipLines = software.SkipLines
		result = tx.Create(&minerSoftware)
		if result.Error != nil {
			log.Fatalf("Issue creating miner software.\n", result.Error)
		}
	} else if minerSoftware.ID > 0 {
		// Update it, in case the settings have changed.
		minerSoftware.Name = software.Name
		minerSoftware.Website = software.ReleaseWebsite
		minerSoftware.AlgoParam = software.AlgoParam
		minerSoftware.PoolParam = software.PoolParam
		minerSoftware.PasswordParam = software.PasswordParam
		minerSoftware.WalletParam = software.WalletParam
		minerSoftware.FileParam = software.FileParam
		minerSoftware.BenchmarkParams = software.BenchmarkParam
		minerSoftware.OtherParams = software.OtherParams
		minerSoftware.SkipLines = software.SkipLines
		result = tx.Save(&minerSoftware)
		if result.Error != nil {
			log.Fatalf("Issue updating miner software, "+software.Name+".\n", result.Error)
		}
	} else if result.Error != nil {
		log.Fatalf("Unknown issue storing miner software.\n", result.Error)
	}

	// Create the map to Zergpool algos.

	// Verify there are records in Algorithms.
	// If not, error out and let the user know they need to run the Zergpool statistics program
	// first (or another pool provider statistics program).
	var algo Algorithm
	result = tx.First(&algo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalf("No algorithms exist in the algorithm table. Run a pool statistics collection " +
			"before calculating miner statistics. Only pool algorithms are calculated to avoid " +
			"wasting time calculating useless statistics. To run a pool statistics collection, " +
			"see Mining-Automation-Zergpool.com as an example.")
	}

	// Cycle over the software algorithms and map them. Check if the algorithm exists in the
	// algorithms table. Do that by using the pool name. If the value is blank, use the miner name.
	for _, softwareAlgo := range software.AlgoConfigs {
		algoToFind := softwareAlgo.PoolName
		if algoToFind == "" {
			algoToFind = softwareAlgo.MinerName
		}
		algo = Algorithm{} // Reset to avoid any collisions.
		result = tx.Where("LOWER(name) = ?", strings.ToLower(algoToFind)).Limit(1).Find(&algo)
		// Skip anything not in the database as it likely is not in use by a pool.
		if algo.ID > 0 {
			var minerAlgos []MinerSoftwareAlgos
			var minerAlgo MinerSoftwareAlgos

			// Check if the algorithm is already mapped. If so, just update it to ensure accuracy.
			result = tx.Where("miner_software_id = ? AND algorithm_id = ?",
				minerSoftware.ID, algo.ID).Find(&minerAlgos)
			if len(minerAlgos) > 0 { // Update all the miner algorithms in case something changed.
				for _, mAlgo := range minerAlgos {
					mAlgo.MinerSoftwareID = minerSoftware.ID
					mAlgo.AlgorithmID = algo.ID
					// This is what the software will require in params.
					mAlgo.Name = softwareAlgo.MinerName
					// Store any extra parameters required for the algorithm.
					mAlgo.ExtraParams = softwareAlgo.ExtraParams
					result = tx.Save(&mAlgo)
					if result.Error != nil {
						log.Fatalf("Issue updating miner software algo map for "+
							mAlgo.Name+".\n", result.Error)
					}
				}
			} else { // Create for the first time.
				minerAlgo.MinerSoftwareID = minerSoftware.ID
				minerAlgo.AlgorithmID = algo.ID
				// This is what the software will require in params.
				minerAlgo.Name = softwareAlgo.MinerName
				// Store any extra parameters required for the algorithm.
				minerAlgo.ExtraParams = softwareAlgo.ExtraParams
				result = tx.Create(&minerAlgo)
				if result.Error != nil {
					log.Fatalf("Issue creating miner software algo map for "+
						minerAlgo.Name+".\n", result.Error)
				}
			}
		}
	}

	return minerSoftware
}

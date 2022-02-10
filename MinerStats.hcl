// Database Connectivity
host="localhost"
port="5432"
database="mining"
user="postgres"
password="whateves"
timezone="America/Chicago"

// Miner Stats Configuration
minerName="Your Miner"  // An identifier for the mining hardware.
// Optional, as some software requires to connect to a pool which may require a wallet
wallet=""  

// Mining software definitions
software "cpuminer-opt" {
  releaseWebsite = "https://github.com/JayDDee/cpuminer-opt/releases"
  // The path to the mining software. If this is blank, it is skipped.
  filePath = ""
  algoParam = "--algo"
  connectForAssessment = false  // Supports benchmarking without connection
  poolParam = "--url"  // Used for the pool connection
  passwordParam = "--pass"
  walletParam = "--user"
  benchmarkParam = "--benchmark" // Runs in benchmark mode without URL target
  otherParams = ""  // Any other important params
  // This is used to find the hash rate in the mining program's screen output (which is saved to a file).
  statSearchPhrase = "Total:"
  // The amount of time to wait before checking output for statistics, in seconds.
  // It can be helpful to give the program a few minutes sometimes, as it often calculates an average
  // hash rate instead of a current hash rate.
  statWaitTime=60
  // How many lines to skip on the output. If the software outputs low hashrate initially, use this to
  // skip those values. 1 will skip 1 line of hashrate output.
  skipLines=3

  // Algorithm maps - The mining software may not use the pool's algo name (stored in the algorithm table).
  // If so, this can be used to map the mining name to the pool name.
  algo "allium" { //        Garlicoin (GRLC)
    // A blank pool name means either the pool does not support it or the pool name matches the
    // miner's name for the algo.
    poolName = "" 
  }
  algo "anime" {
    poolName = "" //              Animecoin (ANI)
  }
  algo "argon2" {
    poolName = "" //              Argon2 Coin (AR2)
  }
  algo "argon2d250" {
    poolName = ""
  }
  algo "argon2d500" {
    poolName = "argon2d-dyn" //   argon2d-dyn, Dynamic (DYN)
  }
  algo "argon2d4096" {
    poolName = ""            //   argon2d-uis, Unitus (UIS)
  }
  algo "axiom" {
    poolName = ""            //   Shabal-256 MemoHash
  }
  algo "blake" {
    poolName = ""            //   blake256r14 (SFR)
  }
  algo "blake2b" {
    poolName = ""            //   Blake2b 256
  }
  algo "blake2s" {
    poolName = ""            //   Blake-2 S
  }
  algo "blakecoin" {
    poolName = ""            //   blake256r8
  }
  algo "bmw" {
    poolName = ""            //   BMW 256
  }
  algo "bmw512" {
    poolName = ""            //   BMW 512
  }
  algo "c11" {
    poolName = ""            //   Chaincoin
  }
  algo "decred" {
    poolName = ""            //   Blake256r14dcr
  }
  algo "deep" {
    poolName = ""            //   Deepcoin (DCN)
  }
  algo "dmd-gr" {
    poolName = ""            //   Diamond
  }
  algo "groestl" {
    poolName = ""            //   Groestl coin
  }
  algo "hex" {
    poolName = ""            //   x16r-hex
  }
  algo "hmq1725" {
    poolName = ""            //   Espers
  }
  algo "hodl" {
    poolName = ""            //   Hodlcoin
  }
  algo "jha" {
    poolName = ""            //   jackppot (Jackpotcoin)
  }
  algo "keccak" {
    poolName = ""            //   Maxcoin
  }
  algo "keccakc" {
    poolName = ""            //   Creative Coin
  }
  algo "lbry" {
    poolName = ""            //   LBC, LBRY Credits
  }
  algo "lyra2h" {
    poolName = ""            //   Hppcoin
  }
  algo "lyra2re" {
    poolName = ""            //   lyra2
  }
  algo "lyra2rev2" {
    poolName = "lyra2v2"     //   lyrav2
  }
  algo "lyra2rev3" {
    poolName = ""            //   lyrav2v3
  }
  algo "lyra2z" {
    poolName = ""
  }
  algo "lyra2z330" {
    poolName = ""    //           Lyra2 330 rows
  }
  algo "m7m" {
    poolName = ""    //           Magi (XMG)
  }
  algo "myr-gr" {
    poolName = ""    //           Myriad-Groestl
  }
  algo "minotaur" {
    poolName = ""    //           Ringcoin (RNG)
  }
  algo "neoscrypt" {
    poolName = ""    //           NeoScrypt(128, 2, 1)
  }
  algo "nist5" {
    poolName = ""    //           Nist5
  }
  algo "pentablake" {
    poolName = ""    //           5 x blake512
  }
  algo "phi1612" {
    poolName = "phi" //           phi
  }
  algo "phi2" {
    poolName = ""
  }
  algo "polytimos" {
    poolName = ""
  }
  algo "power2b" {
    poolName = "" //              MicroBitcoin (MBC)
  }
  algo "quark" {
    poolName = "" //              Quark
  }
  algo "qubit" {
    poolName = "" //              Qubit
  }
  algo "scrypt" {
    poolName = "" //              scrypt(1024, 1, 1) (default)
  }
  algo "scrypt" { //              scryptn2 uses scrypt with special parameters
    poolName = "scryptn2"
    extraParams = "--param-n 1048576"
  }
  algo "scrypt:N" {
    poolName = "" //              scrypt(N, 1, 1)
  }
  algo "sha256d" {
    poolName = "" //              Double SHA-256
  }
  algo "sha256q" {
    poolName = "" //              Quad SHA-256, Pyrite (PYE)
  }
  algo "sha256t" {
    poolName = "" //              Triple SHA-256, Onecoin (OC)
  }
  algo "sha3d" {
    poolName = "" //              Double Keccak256 (BSHA3)
  }
  algo "shavite3" {
    poolName = "" //              Shavite3
  }
  algo "skein" {
    poolName = "" //              Skein+Sha (Skeincoin)
  }
  algo "skein2" {
    poolName = "" //              Double Skein (Woodcoin)
  }
  algo "skunk" {
    poolName = "" //              Signatum (SIGT)
  }
  algo "sonoa" {
    poolName = "" //              Sono
  }
  algo "timetravel" {
    poolName = "" //              timeravel8, Machinecoin (MAC)
  }
  algo "timetravel10" {
    poolName = "" //              Bitcore (BTX)
  }
  algo "tribus" {
    poolName = "" //              Denarius (DNR)
  }
  algo "vanilla" {
    poolName = "" //              blake256r8vnl (VCash)
  }
  algo "veltor" {
    poolName = ""
  }
  algo "verthash" {
    poolName = ""
  }
  algo "whirlpool" {
    poolName = ""
  }
  algo "whirlpoolx" {
    poolName = ""
  }
  algo "x11" {
    poolName = ""    //           Dash
  }
  algo "x11evo" {
    poolName = ""    //           Revolvercoin (XRE)
  }
  algo "x11gost" {
    poolName = "sib" //           sib (SibCoin)
  }
  algo "x12" {
    poolName = ""    //           Galaxie Cash (GCH)
  }
  algo "x13" {
    poolName = ""    //           X13
  }
  algo "x13bcd" {
    poolName = ""    //           bcd
  }
  algo "x13sm3" {
    poolName = ""    //           hsr (Hshare)
  }
  algo "x14" {
    poolName = ""    //           X14
  }
  algo "x15" {
    poolName = ""    //           X15
  }
  algo "x16r" {
    poolName = ""
  }
  algo "x16rv2" {
    poolName = ""
  }
  algo "x16rt" {
    poolName = ""	//        Gincoin (GIN)
  }
  algo "x16rt-veil" {
    poolName = ""	//        Veil (VEIL)
  }
  algo "x16s" {
    poolName = ""
  }
  algo "x17" {
    poolName = ""
  }
  algo "x21s" {
    poolName = ""
  }
  algo "x22i" {
    poolName = ""
  }
  algo "x25x" {
    poolName = ""
  }
  algo "xevan" {
    poolName = ""     // Bitsend (BSD)
  }
  algo "yescrypt" {
    poolName = ""     // Globalboost-Y (BSTY)
  }
  algo "yescryptr8" {
    poolName = ""     // BitZeny (ZNY)
  }
  algo "yescryptr8g" {
    poolName = ""     // Koto (KOTO)
  }
  algo "yescryptr16" {
    poolName = "yescryptR16" //   Eli
  }
  algo "yescryptr32" {
    poolName = "yescryptR32" //   WAVI
  }
  algo "yespower" {
    poolName = ""            //   Cryply
  }
  algo "yespowerr16" {
    poolName = "yespowerR16" //   Yenten (YTN)
  }
  algo "yespower-b2b" {
    poolName = ""            //    generic yespower + blake2b
  }
  algo "zr5" {
    poolName = ""            //    Ziftr
  }
}

software "SRBMiner-Multi" {
  releaseWebsite = "https://github.com/doktor83/SRBMiner-Multi/releases"
  // The path to the mining software. If this is blank, it is skipped.
  filePath = ""
  algoParam = "--algorithm"
  connectForAssessment = true
  poolParam = "--pool"  // Requires a pool to generate stats. Optional.
  passwordParam = "--password"
  walletParam = "--wallet" // Requires a wallet to connect. Optional.
  fileParam = "--log-file" // Some software can log to a file. Optional.
  // NOTE: During tests, great performance was identified setting intensity to 4. Priority set to 1 prevents
  // the miner from overwhelming system/other important processes. This can be set to the default of 2, if desired.
  // To enable GPU mining, remove --disable-gpu.
  otherParams = "--disable-gpu --cpu-threads 0 --cpu-threads-priority 1 --cpu-threads-intensity 4"
  // This is used to find the hash rate in the mining program's screen output (which is saved to a file).
  statSearchPhrase = "Total:"
  // The amount of time to wait before checking output for statistics, in seconds.
  // It can be helpful to give the program a few minutes sometimes, as it often calculates an average
  // hash rate instead of a current hash rate.
  statWaitTime=140
  // How many lines to skip on the output. If the software outputs low hashrate initially, use this to
  // skip those values. 1 will skip 1 line of hashrate output.
  skipLines=0

  // Algorithm maps - The mining software may not use the pool's algorithm name (in the algorithm table).
  // If so, this can be used to map the mining name to the pool name.
  algo "argon2d_dynamic" {
    poolName = "argon2d-dyn" //   argon2d-dyn, Dynamic (DYN)
  }
  algo "blake2b" {
    poolName = ""            //   Blake2b 256
  }
  algo "blake2s" {
    poolName = ""            //   Blake-2 S
  }
  algo "cpupower" {
    poolName = ""
  }
  algo "curvehash" {
    poolName = ""
  }
  algo "cryptonight_gpu" {
    poolName = ""
  }
  algo "cryptonight_upx" {
    poolName = ""
  }
  algo "cryptonight_xhv" {
    poolName = "cryptonight_haven"
  }
  algo "etchash" {
    poolName = ""
  }
  algo "ethash" {
    poolName = ""
  }
  algo "firopow" {
    poolName = ""
  }
  algo "ghostrider" {
    poolName = ""
  }
  algo "heavyhash" {
    poolName = ""
  }
  algo "k12" {
    poolName = ""
  }
  algo "kawpow" {
    poolName = ""
  }
  algo "keccak" {
    poolName = ""            //   Maxcoin
  }
  algo "minotaurx" {
    poolName = ""
  }
  algo "panthera" {
    poolName = ""
  }
  algo "randomarq" {
    poolName = ""
  }
  algo "randomx" {
    poolName = ""
  }
  algo "scryptn2" {
    poolName = ""
  }
  algo "ubqhash" {
    poolName = ""
  }
  algo "verthash" {
    poolName = ""
  }
  algo "verushash" {
    poolName = ""
  }
  algo "yescrypt" {
    poolName = ""     // Globalboost-Y (BSTY)
  }
  algo "yescryptr16" {
    poolName = "" //   Eli
  }
  algo "yescryptr32" {
    poolName = "" //   WAVI
  }
  algo "yespower" {
    poolName = ""            //   Cryply
  }
  algo "yespowerr16" {
    poolName = "" //   Yenten (YTN)
  }
  algo "yespowertide" {
    poolName = ""
  }
  algo "yespowerurx" {
    poolName = ""
  }
}

software "cpuminer-rplant" {
  releaseWebsite = "https://github.com/rplant8/cpuminer-opt-rplant/releases"
  // The path to the mining software. If this is blank, it is skipped.
  filePath = ""
  algoParam = "--algo"
  connectForAssessment = false  // Supports benchmarking without connection
  poolParam = "--url"  // Used for the pool connection
  passwordParam = "--pass"
  walletParam = "--user"
  benchmarkParam = "--benchmark" // Runs in benchmark mode without URL target
  otherParams = ""  // Any other important params
  // This is used to find the hash rate in the mining program's screen output (which is saved to a file).
  statSearchPhrase = "Total:"
  // The amount of time to wait before checking output for statistics, in seconds.
  // It can be helpful to give the program a few minutes sometimes, as it often calculates an average
  // hash rate instead of a current hash rate.
  statWaitTime=60
  // How many lines to skip on the output. If the software outputs low hashrate initially, use this to
  // skip those values. 1 will skip 1 line of hashrate output.
  skipLines=4

  // Algorithm maps - The mining software may not use the pool's algo name (stored in the algorithm table).
  // If so, this can be used to map the mining name to the pool name.
  algo "allium" { //        Garlicoin (GRLC)
    // A blank pool name means either the pool does not support it or the pool name matches the
    // miner's name for the algo.
    poolName = "" 
  }
  algo "anime" {
    poolName = "" //              Animecoin (ANI)
  }
  algo "argon2d250" {
    poolName = ""
  }
  algo "argon2d500" {
    poolName = "argon2d-dyn" //   argon2d-dyn, Dynamic (DYN)
  }
  algo "argon2d4096" {
    poolName = ""            //   argon2d-uis, Unitus (UIS)
  }
  algo "axiom" {
    poolName = ""            //   Shabal-256 MemoHash
  }
  algo "blake" {
    poolName = ""            //   blake256r14 (SFR)
  }
  algo "blake2b" {
    poolName = ""            //   Blake2b 256
  }
  algo "blake2s" {
    poolName = ""            //   Blake-2 S
  }
  algo "blakecoin" {
    poolName = ""            //   blake256r8
  }
  algo "bmw" {
    poolName = ""            //   BMW 256
  }
  algo "bmw512" {
    poolName = ""            //   BMW 512
  }
  algo "cpupower" {
    poolName = ""
  }
  algo "curvehash" {
    poolName = ""
  }
  algo "c11" {
    poolName = ""            //   Chaincoin
  }
  algo "decred" {
    poolName = ""            //   Blake256r14dcr
  }
  algo "deep" {
    poolName = ""            //   Deepcoin (DCN)
  }
  algo "dmd-gr" {
    poolName = ""            //   Diamond
  }
  algo "gr" {
    poolName = "ghostrider"
  }
  algo "groestl" {
    poolName = ""            //   Groestl coin
  }
  algo "heavyhash" {
    poolName = ""
  }
  algo "hex" {
    poolName = ""            //   x16r-hex
  }
  algo "hmq1725" {
    poolName = ""            //   Espers
  }
  algo "hodl" {
    poolName = ""            //   Hodlcoin
  }
  algo "jha" {
    poolName = ""            //   jackppot (Jackpotcoin)
  }
  algo "keccak" {
    poolName = ""            //   Maxcoin
  }
  algo "keccakc" {
    poolName = ""            //   Creative Coin
  }
  algo "lbry" {
    poolName = ""            //   LBC, LBRY Credits
  }
  algo "lyra2h" {
    poolName = ""            //   Hppcoin
  }
  algo "lyra2re" {
    poolName = ""            //   lyra2
  }
  algo "lyra2rev2" {
    poolName = "lyra2v2"     //   lyrav2
  }
  algo "lyra2rev3" {
    poolName = ""            //   lyrav2v3
  }
  algo "lyra2z" {
    poolName = ""
  }
  algo "lyra2z330" {
    poolName = ""    //           Lyra2 330 rows
  }
  algo "myr-gr" {
    poolName = ""    //           Myriad-Groestl
  }
  algo "minotaur" {
    poolName = ""    //           Ringcoin (RNG)
  }
  algo "minotaurx" {
    poolName = ""
  }
  algo "neoscrypt" {
    poolName = ""    //           NeoScrypt(128, 2, 1)
  }
  algo "nist5" {
    poolName = ""    //           Nist5
  }
  algo "pentablake" {
    poolName = ""    //           5 x blake512
  }
  algo "phi1612" {
    poolName = "phi" //           phi
  }
  algo "phi2" {
    poolName = ""
  }
  algo "polytimos" {
    poolName = ""
  }
  algo "power2b" {
    poolName = "" //              MicroBitcoin (MBC)
  }
  algo "quark" {
    poolName = "" //              Quark
  }
  algo "qubit" {
    poolName = "" //              Qubit
  }
  algo "scrypt" {
    poolName = "" //              scrypt(1024, 1, 1) (default)
  }
  algo "scrypt" { //              scryptn2 uses scrypt with special parameters
    poolName = "scryptn2"
    extraParams = "--param-n 1048576"
  }
  algo "scrypt:N" {
    poolName = "" //              scrypt(N, 1, 1)
  }
  algo "sha256d" {
    poolName = "" //              Double SHA-256
  }
  algo "sha256q" {
    poolName = "" //              Quad SHA-256, Pyrite (PYE)
  }
  algo "sha256t" {
    poolName = "" //              Triple SHA-256, Onecoin (OC)
  }
  algo "sha3d" {
    poolName = "" //              Double Keccak256 (BSHA3)
  }
  algo "shavite3" {
    poolName = "" //              Shavite3
  }
  algo "skein" {
    poolName = "" //              Skein+Sha (Skeincoin)
  }
  algo "skein2" {
    poolName = "" //              Double Skein (Woodcoin)
  }
  algo "skunk" {
    poolName = "" //              Signatum (SIGT)
  }
  algo "sonoa" {
    poolName = "" //              Sono
  }
  algo "timetravel" {
    poolName = "" //              timeravel8, Machinecoin (MAC)
  }
  algo "timetravel10" {
    poolName = "" //              Bitcore (BTX)
  }
  algo "tribus" {
    poolName = "" //              Denarius (DNR)
  }
  algo "vanilla" {
    poolName = "" //              blake256r8vnl (VCash)
  }
  algo "veltor" {
    poolName = ""
  }
  algo "whirlpool" {
    poolName = ""
  }
  algo "whirlpoolx" {
    poolName = ""
  }
  algo "x11" {
    poolName = ""    //           Dash
  }
  algo "x11evo" {
    poolName = ""    //           Revolvercoin (XRE)
  }
  algo "x11gost" {
    poolName = "sib" //           sib (SibCoin)
  }
  algo "x12" {
    poolName = ""    //           Galaxie Cash (GCH)
  }
  algo "x13" {
    poolName = ""    //           X13
  }
  algo "x13bcd" {
    poolName = ""    //           bcd
  }
  algo "x13sm3" {
    poolName = ""    //           hsr (Hshare)
  }
  algo "x14" {
    poolName = ""    //           X14
  }
  algo "x15" {
    poolName = ""    //           X15
  }
  algo "x16r" {
    poolName = ""
  }
  algo "x16rv2" {
    poolName = ""
  }
  algo "x16rt" {
    poolName = ""	//        Gincoin (GIN)
  }
  algo "x16rt-veil" {
    poolName = ""	//        Veil (VEIL)
  }
  algo "x16s" {
    poolName = ""
  }
  algo "x17" {
    poolName = ""
  }
  algo "x21s" {
    poolName = ""
  }
  algo "x22i" {
    poolName = ""
  }
  algo "x25x" {
    poolName = ""
  }
  algo "xevan" {
    poolName = ""     // Bitsend (BSD)
  }
  algo "yescrypt" {
    poolName = ""     // Globalboost-Y (BSTY)
  }
  algo "yescryptr8" {
    poolName = ""     // BitZeny (ZNY)
  }
  algo "yescryptr8g" {
    poolName = ""     // Koto (KOTO)
  }
  algo "yescryptr16" {
    poolName = "yescryptR16" //   Eli
  }
  algo "yescryptr32" {
    poolName = "yescryptR32" //   WAVI
  }
  algo "yespower" {
    poolName = ""            //   Cryply
  }
  algo "yespowerr16" {
    poolName = "yespowerR16" //   Yenten (YTN)
  }
  algo "yespower-b2b" {
    poolName = ""            //    generic yespower + blake2b
  }
  algo "yespowertide" {
    poolName = ""
  }
  algo "yespowerurx" {
    poolName = ""
  }
  algo "zr5" {
    poolName = ""            //    Ziftr
  }
}

software "XMRRig" {
  releaseWebsite = "https://github.com/xmrig/xmrig"
  // The path to the mining software. If this is blank, it is skipped.
  filePath = "C:\\Mining\\xmrig-6.16.4\\xmrig.exe"
  algoParam = "--algo"
  connectForAssessment = true
  poolParam = "--url"  // Requires a pool to generate stats. Optional.
  passwordParam = "--pass"
  walletParam = "--user" // Requires a wallet to connect. Optional.
  fileParam = "--log-file" // Some software can log to a file. Optional.
  // Forces hash rates to print every 15 seconds.
  otherParams = "--print-time 15"
  // This is used to find the hash rate in the mining program's screen output (which is saved to a file).
  statSearchPhrase = "10s/60s/15m"
  // The amount of time to wait before checking output for statistics, in seconds.
  // It can be helpful to give the program a few minutes sometimes, as it often calculates an average
  // hash rate instead of a current hash rate.
  statWaitTime=140
  // How many lines to skip on the output. If the software outputs low hashrate initially, use this to
  // skip those values. 1 will skip 1 line of hashrate output.
  skipLines=0
  // This will grab the first number (10s average) and the hash rate unit after the 15m number.
  skipTokens=2

  // Algorithm maps - The mining software may not use the pool's algorithm name (in the algorithm table).
  // If so, this can be used to map the mining name to the pool name.
  // NOTE: cn/ccx did not work in testing with Zergpool. Commenting out.
  //algo "cn/ccx" {
  //  poolName = "cryptonight_gpu"
  //}
  algo "cn/upx2" {
    poolName = "cryptonight_upx"
  }
  algo "cn-heavy/xhv" {
    poolName = "cryptonight_haven"
  }
  // NOTE: During testing, Zergpool indicated this was incompatible/disabled. Might have been a temporary outage?
  algo "kawpow" {
    poolName = ""
  }
  algo "rx/arq" {
    poolName = "randomARQ"
  }
  algo "rx/0" {
    poolName = "randomx"
  }
}

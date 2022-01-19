// Database Connectivity
host="localhost"
port="5432"
database="mining"
user="postgres"
password="whateves"
timezone="America/Chicago"

// Miner Stats Configuration
minerName="YourMiner"  // An identifier for the mining hardware.
// Alter this according to your path/program.
minerPath="C:\\Mining\\cpuminer-opt-3.19.3\\cpuminer-avx2-sha.exe"
// This is used to find the hash rate in the mining program's screen output (which is saved to a file).
statSearchPhrase="Total:"
allowGPU=0  // 1=Yes, 0=No - Whether to include GPU(s) when calculating statistics.
// The amount of time to wait before checking output for statistics, in seconds.
// It can be helpful to give the program a few minutes sometimes, as it often calculates an average
// hash rate instead of a current hash rate.
statWaitTime=60

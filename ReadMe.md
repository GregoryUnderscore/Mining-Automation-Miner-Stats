# **Miner Stats for CPU and GPU Miners**

### **Summary**
Pull and save miner statistics for various mining software (cpuminer-opt supported out-of-the-box). Various SQL queries can then be utilized to aid in mining automation or predictions.

### **Important**
Pool provider statistics are required for profitability estimates/actuals. Before using this, please see the instructions at: https://github.com/GregoryUnderscore/Mining-Automation-ZergPool.com

### **Description**
ZergPool provides several useful statistics for every pool they host. This allows a miner to calculate projections
and possible profit opportunities. However, to properly calculate these projections, a miner's hash rate must be calculated
for all supported algorithms. This can be a painstaking process when done manually. This program makes the process easier. 

In short, the minerStats.go program does the following:
1. Connect to a database defined in the configuration file, MinerStats.hcl.
2. Automatically create the required schema.
3. Determine all algorithms stored for all pools.
4. Generate hash rates for the algorithms and store them into the database.

### **How to Use**

1. Follow the instructions first at https://github.com/GregoryUnderscore/Mining-Automation-ZergPool.com
2. Clone this repository
3. Update the MinerStats.hcl file with the appropriate details.
4. go run minerStats.go

### **Included Reports**
In the sql folder are SQL reports. There are reports to see profitability estimates and actuals (daily). Also, there is a report
that shows all statistics.
# **Miner Stats for CPU and GPU Miners**

### **Summary**
Pull and save miner statistics for various mining software. Cpuminer-opt, SRBMiner, and Cpuminer-rplant are supported 
out-of-the-box. You can add your own setup too for other mining software. Tested on Linux and Windows. Various SQL queries can 
also be utilized to aid in mining automation or predictions (/sql folder).

### **Important**
Pool provider statistics are required for profitability estimates/actuals. Before using this, please see the instructions at: https://github.com/GregoryUnderscore/Mining-Automation-ZergPool.com

### **Description**
ZergPool provides several useful statistics for every pool they host. This allows a miner to calculate projections
and possible profit opportunities. However, to properly calculate these projections, a miner's hash rate must be calculated
for all supported algorithms. This can be a painstaking process when done manually. This program makes the process easier. 

In short, the minerStats.go program does the following:
1. Connects to a database defined in the configuration file, MinerStats.hcl.
2. Automatically creates the required schema.
3. Determines all algorithms stored for all pools.
4. Generates hash rates for the algorithms and store them into the database according to miner software settings
in MinerStats.hcl.

### **How to Use**

1. Follow the instructions first at https://github.com/GregoryUnderscore/Mining-Automation-ZergPool.com
2. Download the latest release and extract it to a folder.
3. Update the MinerStats.hcl file with the appropriate details. There is no need to change the miner software algorithm and 
parameter settings unless you know they are incorrect. They should be accurate (or close). You will need to add your file path to the miner software setup.
4. Run the minerStats.exe or minerStats (on Linux) and wait for it to finish. It could take a few hours if all supported mining software are assessed.

### **Included Reports**
In the sql folder are SQL reports. There are reports to see profitability estimates and actuals (daily). Also, there is a report
that shows all statistics.

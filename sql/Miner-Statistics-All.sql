-- Reports all the statistics for all miners in the database.

SELECT m.name AS "Miner", miner_softwares.name AS "Software", msa."name" AS "Algo", a.name AS "Algo", miner_stats.instant, miner_stats.work_per_second AS "Work", 	
	CASE 
		WHEN mh_factor = 1000000000 THEN 'PH/s'
		WHEN mh_factor = 1000000 THEN 'TH/s'
		WHEN mh_factor = 1000 THEN 'GH/s'
		WHEN mh_factor = 1 THEN 'MH/s'
		WHEN mh_factor = 0.001 THEN 'kH/s'
		WHEN mh_factor = 0.000001 THEN 'H/s'
	END AS "Unit"
FROM miners m
INNER JOIN miner_stats ON
	m.id = miner_stats.miner_id
INNER JOIN miner_softwares ON
	miner_softwares.id = miner_stats.miner_software_id
INNER JOIN algorithms a ON
	a.id = miner_stats.algorithm_id 
INNER JOIN miner_software_algos msa ON
	msa.algorithm_id = a.id AND msa.miner_software_id = miner_softwares.id
ORDER BY m.name, a.name;
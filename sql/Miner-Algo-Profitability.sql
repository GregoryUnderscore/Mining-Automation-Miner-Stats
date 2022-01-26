-- Reports the estimated profitability and 24 hour profit for a miner/algo.
-- The average of all work in the database is used for the calculations.
-- The most profitable algorithm will be at the top (actual 24 hour profit).
-- NOTE: Some algorithms may generate multiple units during their assessments. 
-- If so, those are shown separately, and this will cause apparent duplicates to
-- appear in the list. In general, those are outliers and can be ignored (and will
-- be low in the profit list anyway).

SELECT m.name AS "Miner", ms2.name AS "Software", msa.name AS "Miner Algo", a.name AS "Algo", average_work AS "Work", 	
	CASE 
		WHEN average_stat.mh_factor = 1000000000 THEN 'PH/s'
		WHEN average_stat.mh_factor = 1000000 THEN 'TH/s'
		WHEN average_stat.mh_factor = 1000 THEN 'GH/s'
		WHEN average_stat.mh_factor = 1 THEN 'MH/s'
		WHEN average_stat.mh_factor = 0.001 THEN 'kH/s'
		WHEN average_stat.mh_factor = 0.000001 THEN 'H/s'
	END AS "Unit",
	price*profit_estimate*(average_stat.mh_factor / pools.mh_factor)*average_work AS "D Profit Est",  -- The $ estimate per day
	price*0.001*profit_actual24_hours*(average_stat.mh_factor / pools.mh_factor)*average_work AS "24Hr Profit", -- The $ actual last DAY
	pools.port, pools.mh_factor 
FROM miners m
INNER JOIN ( -- Pulls the latest miner statistic TO use FOR hash factor calculations
	SELECT miner_id, miner_software_id, algorithm_id, max(id) AS latest_stat_id
	FROM miner_stats ms
	GROUP BY ms.miner_id, ms.miner_software_id, ms.algorithm_id 
) latest_stat ON
	latest_stat.miner_id = m.id
INNER JOIN miner_stats ON
	miner_stats.id = latest_stat.latest_stat_id
INNER JOIN miner_softwares ms2 ON
	latest_stat.miner_software_id = ms2.id
INNER JOIN algorithms a ON
	a.id = latest_stat.algorithm_id 
INNER JOIN miner_software_algos msa ON
	msa.algorithm_id = a.id AND msa.miner_software_id = ms2.id
INNER JOIN pools ON
	pools.algorithm_id = a.id
INNER JOIN pool_stats ON
	pool_stats.pool_id = pools.id
INNER JOIN (  -- Pulls the latest pool statistic FOR a pool/algorithm.
	SELECT max(id) AS id
	FROM pool_stats ps 
	GROUP BY ps.pool_id 
) latest_pool_stat ON
	latest_pool_stat.id = pool_stats.id
INNER JOIN coin_prices cp ON
	pool_stats.coin_price_id = cp.id
INNER JOIN ( -- Pulls the average hash rate FOR a miner/algo
	SELECT miner_id, miner_software_id, algorithm_id, AVG(work_per_second) AS average_work, mh_factor  -- Should pull the average statistics FOR a miner/software/algo
	FROM miner_stats ms
	GROUP BY ms.miner_id, ms.miner_software_id, ms.algorithm_id, mh_factor 
) average_stat ON
	average_stat.miner_id = m.id AND average_stat.miner_software_id = ms2.id AND average_stat.algorithm_id = a.id
ORDER BY (price*0.001*profit_actual24_hours*(average_stat.mh_factor / pools.mh_factor)*average_work) DESC;
CREATE MATERIALIZED VIEW if not EXISTS tx_count_hourly
WITH (timescaledb.continuous) AS
	select 
		time_bucket('1 hour'::interval, time) AS timestamp,
		sum(tx_count) as tx_count,
		sum(tx_count) / 3600.0 as tps
	from block_stats
	group by timestamp
CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_minute
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 minute'::interval, time) AS ts,
		sum(bytes_in_block) as bytes_in_block,
		(sum(blobs_size)/60.0) as bps,
		max(case when block_time > 0 then blobs_size::float/(block_time/1000.0) else 0 end) as bps_max,
		min(case when block_time > 0 then blobs_size::float/(block_time/1000.0) else 0 end) as bps_min,
		(sum(tx_count)/60.0) as tps,
		max(case when block_time > 0 then tx_count::float/(block_time/1000.0) else 0 end) as tps_max,
		min(case when block_time > 0 then tx_count::float/(block_time/1000.0) else 0 end) as tps_min,
		avg(block_time) as block_time,
		percentile_agg(block_time) as block_time_pct,
		sum(blobs_size) as blobs_size,
		sum(blobs_count) as blobs_count,
		sum(tx_count) as tx_count,
		sum(events_count) as events_count,
		sum(fee) as fee,
		sum(supply_change) as supply_change,
		sum(rewards) as rewards,
		sum(commissions) as commissions,
		sum(gas_limit) as gas_limit,
		sum(gas_used) as gas_used,
		(case when sum(gas_limit) > 0 then sum(fee) / sum(gas_limit) else 0 end) as gas_price,
		(case when sum(gas_limit) > 0 then sum(gas_used) / sum(gas_limit) else 0 end) as gas_efficiency
	from block_stats
	group by 1
	order by 1 desc
	with no data;

CALL add_view_refresh_job('block_stats_by_minute', NULL, INTERVAL '1 minute');

CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_year
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 year', day.ts) AS ts,
		sum(bytes_in_block) as bytes_in_block,
		sum(blobs_size)/(count(*) * 86400.0) as bps,
		max(bps_max) as bps_max,
		min(bps_min) as bps_min,
		sum(tx_count)/(count(*) * 86400.0) as tps,
		max(tps_max) as tps_max,
		min(tps_min) as tps_min,
		mean(rollup(block_time_pct)) as block_time,
		rollup(block_time_pct) as block_time_pct,
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
	from block_stats_by_day as day
	group by 1
	order by 1 desc
	with no data;

CALL add_view_refresh_job('block_stats_by_year', NULL, INTERVAL '1 hour');
CREATE MATERIALIZED VIEW IF NOT EXISTS block_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=true) AS
	select 
		time_bucket('1 hour'::interval, bbm.ts) AS ts,
		sum(blobs_size)/3600.0 as bps,
		max(bps_max) as bps_max,
		min(bps_min) as bps_min,
		sum(tx_count)/3600.0 as tps,
		max(tps_max) as tps_max,
		min(tps_min) as tps_min,
		avg(block_time) as block_time,
		sum(blobs_size) as blobs_size,
		sum(tx_count) as tx_count,
		sum(events_count) as events_count,
		sum(fee) as fee,
		sum(supply_change) as supply_change,
		sum(gas_limit) as gas_limit,
		sum(gas_used) as gas_used,
		(case when sum(gas_limit) > 0 then sum(fee) / sum(gas_limit) else 0 end) as gas_price,
		(case when sum(gas_limit) > 0 then sum(gas_used) / sum(gas_limit) else 0 end) as gas_efficiency
	from block_stats_by_minute as bbm
	group by 1
	order by 1 desc;

SELECT add_continuous_aggregate_policy('block_stats_by_hour',
  start_offset => NULL,
  end_offset => INTERVAL '1 minute',
  schedule_interval => INTERVAL '15 minute',
  if_not_exists => true)
WHERE NOT (SELECT EXISTS (
    SELECT FROM 
        "_timescaledb_catalog".continuous_agg
    WHERE user_view_schema = 'public' AND user_view_name = 'block_stats_by_hour'
    )
);

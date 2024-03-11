CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 day'::interval, logs.time) AS time, 
        logs.namespace_id, 
        logs.signer_id, 
        sum(logs.size) as size, 
        sum(logs.blobs_count) as blobs_count, 
        max(logs.last_time) as last_time,
        min(logs.first_time) as first_time,
        sum(fee) as fee
    from rollup_stats_by_hour as logs
    group by 1, 2, 3
	with no data;
        
CALL add_view_refresh_job('rollup_stats_by_day', NULL, INTERVAL '5 minute');

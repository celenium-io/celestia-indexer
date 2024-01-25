CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 month'::interval, time) AS time, 
        logs.namespace_id, 
        logs.signer_id, 
        sum(logs.size) as size, 
        sum(logs.blobs_count) as blobs_count, 
        max(logs.last_time) as last_time
    from rollup_stats_by_day as logs
    group by 1, 2, 3;
        
CALL add_view_refresh_job('rollup_stats_by_month', INTERVAL '1 minute', INTERVAL '1 hour');

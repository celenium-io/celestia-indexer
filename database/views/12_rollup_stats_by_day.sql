CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 day'::interval, time) AS time, 
        logs.namespace_id, 
        logs.signer_id, 
        sum(logs.size) as size, 
        sum(logs.blobs_count) as blobs_count, 
        max(logs.last_time) as last_time
    from rollup_stats_by_hour as logs
    group by 1, 2, 3;
        
CALL add_view_refresh_job('rollup_stats_by_day', INTERVAL '1 minute', INTERVAL '1 hour');

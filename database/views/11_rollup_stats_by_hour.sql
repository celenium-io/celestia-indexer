CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, time) AS time, 
        logs.namespace_id, 
        logs.signer_id, 
        sum(logs.size) as size, 
        count(*) as blobs_count, 
        max(logs.time) as last_time
    from blob_log as logs
    group by 1, 2, 3;
        
CALL add_view_refresh_job('rollup_stats_by_hour', INTERVAL '1 minute', INTERVAL '10 minute');

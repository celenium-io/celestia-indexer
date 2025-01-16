CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_tvl_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 month'::interval, logs.time) AS time,
        logs.rollup_id as rollup_id,
        last(logs.value, time) as value
    from tvl as logs
    group by 1, 2
    with no data;
CALL add_view_refresh_job('rollup_tvl_by_month', NULL, INTERVAL '12 hours');
CREATE MATERIALIZED VIEW IF NOT EXISTS tvs_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 day'::interval, logs.time) AS time,
        sum(logs.value) as value
    from tvl as logs
    group by 1
    with no data;
CALL add_view_refresh_job('tvs_by_day', NULL, INTERVAL '4 hours');
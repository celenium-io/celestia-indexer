CREATE MATERIALIZED VIEW IF NOT EXISTS tvs_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 month'::interval, logs.time) AS time,
        min(logs.value)::DOUBLE PRECISION AS min_value,
        max(logs.value)::DOUBLE PRECISION AS max_value,
        last(logs.value, time)::DOUBLE PRECISION AS value
    from tvl as logs
    group by 1
    with no data;
CALL add_view_refresh_job('tvs_by_month', NULL, INTERVAL '12 hours');
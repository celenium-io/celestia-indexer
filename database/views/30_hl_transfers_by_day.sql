CREATE MATERIALIZED VIEW IF NOT EXISTS hl_transfers_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 day'::interval, time) AS time,
        counterparty,
        sum(amount) as amount,
        sum(count) as count
    from hl_transfers_by_hour
    group by 1, 2
    with no data;

CALL add_view_refresh_job('hl_transfers_by_day', NULL, INTERVAL '5 minute');

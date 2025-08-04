CREATE MATERIALIZED VIEW IF NOT EXISTS hl_transfers_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 month'::interval, time) AS time, 
        counterparty,
        sum(amount) as amount,
        sum(count) as count
    from hl_transfers_by_day
    group by 1, 2
    with no data;

CALL add_view_refresh_job('hl_transfers_by_month', NULL, INTERVAL '1 hour');

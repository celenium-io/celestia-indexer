CREATE MATERIALIZED VIEW IF NOT EXISTS hl_transfers_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select
        time_bucket('1 hour'::interval, time) AS time, 
        counterparty,
        sum(amount) as amount,
        count(hl_transfer.id) as count
    from hl_transfer
    group by 1, 2
    with no data;

CALL add_view_refresh_job('hl_transfers_by_hour', NULL, INTERVAL '1 minute');

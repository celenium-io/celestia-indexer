CREATE MATERIALIZED VIEW IF NOT EXISTS accounts_tx_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 day'::interval, time) AS time, 
        tx.address_id,
        sum(fee) as fee,
        sum(gas_wanted) as gas_wanted,
        sum(gas_used) as gas_used,
        sum(count) as count
    from accounts_tx_by_hour as tx
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('accounts_tx_by_day', NULL, INTERVAL '5 minute');

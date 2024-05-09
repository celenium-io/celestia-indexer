CREATE MATERIALIZED VIEW IF NOT EXISTS accounts_tx_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, time) AS time, 
        signer.address_id,
        sum(fee) as fee,
        sum(gas_wanted) as gas_wanted,
        sum(gas_used) as gas_used,
        count(tx.id) as count
    from tx
    inner join signer on signer.tx_id = tx.id
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('accounts_tx_by_hour', NULL, INTERVAL '1 minute');

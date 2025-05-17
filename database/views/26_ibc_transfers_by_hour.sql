CREATE MATERIALIZED VIEW IF NOT EXISTS ibc_transfers_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, time) AS time, 
        channel_id,
        sum(amount) as amount,
        count(ibc_transfer.id) as count
    from ibc_transfer
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('ibc_transfers_by_hour', NULL, INTERVAL '1 minute');

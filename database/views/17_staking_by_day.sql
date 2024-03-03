CREATE MATERIALIZED VIEW IF NOT EXISTS staking_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 day'::interval, time) AS time, 
        logs.validator_id,
        sum(flow) as flow,
        sum(rewards) as rewards,
        sum(commissions) as commissions
    from staking_by_hour as logs
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('staking_by_day', NULL, INTERVAL '5 minute');

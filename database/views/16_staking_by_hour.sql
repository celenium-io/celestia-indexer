CREATE MATERIALIZED VIEW IF NOT EXISTS staking_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, time) AS time, 
        logs.validator_id,
        sum(case when type = 'delegation' then change else 0 end) as flow,
        sum(case when type = 'rewards' and change > 0 then change else 0 end) as rewards,
        sum(case when type = 'commissions' and change > 0 then change else 0 end) as commissions
    from staking_log as logs
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('staking_by_hour', NULL, INTERVAL '1 minute');

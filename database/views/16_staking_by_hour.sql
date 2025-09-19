CREATE MATERIALIZED VIEW IF NOT EXISTS staking_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 hour'::interval, time) AS ts, 
        logs.validator_id,
        sum(case 
            when type = 'delegation' then change
            when type = 'unbonding' then change 
            else 0 
            end) as flow,
        sum(case when type = 'delegation' then change else 0 end) as delegations,
        sum(case when type = 'unbonding' then -change else 0 end) as unbondings,
        sum(case when type = 'rewards' and change > 0 then change else 0 end) as rewards,
        sum(case when type = 'commissions' and change > 0 then change else 0 end) as commissions,
        sum(case when type = 'delegation' then 1 else 0 end) as delegations_count,
        sum(case when type = 'unbonding' then 1 else 0 end) as unbondings_count
    from staking_log as logs
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('staking_by_hour', NULL, INTERVAL '1 minute');

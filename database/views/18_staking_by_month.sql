CREATE MATERIALIZED VIEW IF NOT EXISTS staking_by_month
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
    select 
        time_bucket('1 month'::interval, ts) AS ts, 
        logs.validator_id,
        sum(flow) as flow,
        sum(delegations) as delegations,
        sum(unbondings) as unbondings,
        sum(rewards) as rewards,
        sum(commissions) as commissions,
        sum(delegations_count) as delegations_count,
        sum(unbondings_count) as unbondings_count
    from staking_by_day as logs
    group by 1, 2
	with no data;
        
CALL add_view_refresh_job('staking_by_month', NULL, INTERVAL '1 hour');

CREATE MATERIALIZED VIEW IF NOT EXISTS validator_metrics AS
    with current_state as (
        select * from state limit 1
    ),
    votes as (
        select count(distinct proposal_id) as value, validator_id
        from vote
        left join proposal on proposal.id = proposal_id
        where validator_id is not null and proposal.status = 'applied'
        group by validator_id
    ),
    applied_proposals as (
        select count(*) as value
        from proposal
        where status = 'applied'
    ),
    votes_metric as (
        select 
            (case when applied_proposals.value > 0 then votes.value::float/applied_proposals.value::float else 1 end) as value, 
            validator_id,
            applied_proposals.value as applied_proposals_count,
            votes.value as votes_count
        from votes, applied_proposals
    ),
    self_delegations as (
        select delegation.amount, validator.id 
        from validator
        left join address on address.address = validator.delegator
        left join delegation on address.id = delegation.address_id and validator.id = delegation.validator_id
    ),
    block_missed as (
        select validator_id, count(*) as value 
        from (
            select * from block_signature, current_state
            where block_signature.height >= current_state.last_height - 1000
        ) as signs
        group by validator_id
    )
    select 
        validator.id,
        validator.moniker,
        validator.max_rate,
        validator.max_change_rate,
        validator.creation_time,
        coalesce(votes_metric.applied_proposals_count, 0) as applied_proposals_count, 
        coalesce(votes_metric.votes_count, 0) as votes_count,
        coalesce(self_delegations.amount, 0) as self_delegation_amount,
        validator.stake as stake,
        coalesce(block_missed.value, 0) as block_missed_count,
        coalesce(votes_metric.value, 0) as votes_metric, 
        (1 - (0.65*validator.max_rate + 0.35*validator.max_change_rate)) as commission_metric,
        case when EXTRACT(EPOCH FROM now()) - EXTRACT(EPOCH FROM validator.creation_time) <= 31536000 
        then
            ((EXTRACT(EPOCH FROM now()) - EXTRACT(EPOCH FROM validator.creation_time))/31536000)
        else 1
        end as operation_time_metric,
        case when validator.stake > 0 and self_delegations.amount > validator.min_self_delegation
        then
            1 - (self_delegations.amount - validator.min_self_delegation) / (validator.stake - validator.min_self_delegation)
        else 0
        end as self_delegation_metric,
        coalesce(block_missed.value::float / 1000.0, 0) as block_missed_metric
    from validator
    left join votes_metric on votes_metric.validator_id = validator.id
    left join self_delegations on self_delegations.id = validator.id
    left join block_missed on block_missed.validator_id = validator.id
    order by stake desc;

CALL add_job_refresh_materialized_view();

CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
    select 
        size, 
        blobs_count, 
        last_time, 
        first_time, 
        fee, 
        rollup.*
    from (
        select 
            rollup_id,
            sum(size) as size, 
            sum(blobs_count) as blobs_count, 
            max(last_time) as last_time, 
            min(first_time) as first_time, 
            sum(fee) as fee 
        from (
            select
                namespace_id, 
                signer_id,
                sum(size) as size, 
                sum(blobs_count) as blobs_count, 
                max(last_time) as last_time, 
                min(first_time) as first_time, 
                sum(fee) as fee
            from rollup_stats_by_month
            group by 1, 2
        ) as agg
        inner join rollup_provider as rp on rp.address_id = agg.signer_id AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
        group by 1
    ) as leaderboard
    inner join rollup on rollup.id = leaderboard.rollup_id
    with no data;

CALL add_job_refresh_materialized_view();
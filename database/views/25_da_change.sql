CREATE MATERIALIZED VIEW IF NOT EXISTS da_change AS
    with board1 as (
        select 
            rollup_id,
            sum(size) as size
        from (
            select
                namespace_id, 
                signer_id,
                sum(size) as size
            from rollup_stats_by_hour
            where time > now() - '1 week'::interval
            group by 1, 2
        ) as agg
        inner join rollup_provider as rp on rp.address_id = agg.signer_id AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
        inner join rollup on rollup.id = rp.rollup_id
        where rollup.verified = TRUE
        group by 1
    ), board2 as (
        select 
            rollup_id,
            sum(size) as size
        from (
            select
                namespace_id, 
                signer_id,
                sum(size) as size
            from rollup_stats_by_hour
            where time <= now() - '1 week'::interval and time > now() - '2 week'::interval
            group by 1, 2
        ) as agg
        inner join rollup_provider as rp on rp.address_id = agg.signer_id AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
        inner join rollup on rollup.id = rp.rollup_id
        where rollup.verified = TRUE
        group by 1
    ) 
    select 
        case 
            when coalesce(board2.size, 0) > 0
                then coalesce(board1.size, 0) / coalesce(board2.size, 0) - 1 
            when coalesce(board1.size, 0) > 0 and coalesce(board2.size, 0) = 0
                then 1
            else 0 
            end as da_pct,
        rollup.id as rollup_id
    from rollup
    inner join board1 on rollup.id = board1.rollup_id
    inner join board2 on rollup.id = board2.rollup_id;

CALL add_job_refresh_materialized_view();
CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
   with board as (
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
        inner join rollup_provider as rp on (rp.address_id = agg.signer_id OR rp.address_id = 0) AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
        inner join rollup on rollup.id = rp.rollup_id
        where rollup.verified = TRUE
        group by 1
    ) 
    select 
        board.size, 
        board.blobs_count, 
        board.last_time, 
        board.first_time, 
        board.fee,
        board.size / (select sum(size) from board) as size_pct,
        board.fee / (select sum(fee) from board)as fee_pct,
        board.blobs_count / (select sum(blobs_count) from board)as blobs_count_pct,
        (now() - board.last_time < INTERVAL '1 month') as is_active,
        rollup.*
    from board
    inner join rollup on rollup.id = board.rollup_id;

CALL add_job_refresh_materialized_view();
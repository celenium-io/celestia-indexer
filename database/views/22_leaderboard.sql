CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
    WITH rollups AS (
        SELECT * FROM rollup WHERE verified = TRUE
    ), agg AS MATERIALIZED (
        SELECT
            namespace_id,
            signer_id,
            sum(size)        AS size,
            sum(blobs_count) AS blobs_count,
            max(last_time)   AS last_time,
            min(first_time)  AS first_time,
            sum(fee)         AS fee
        FROM rollup_stats_by_month
        GROUP BY 1, 2
    ), matched AS (
        SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
        FROM agg
        INNER JOIN rollup_provider AS rp
               ON rp.address_id = agg.signer_id AND rp.namespace_id = agg.namespace_id
        UNION ALL
        SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
        FROM agg
        INNER JOIN rollup_provider AS rp
               ON rp.address_id = 0 AND rp.namespace_id = agg.namespace_id
        UNION ALL
        SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
        FROM agg
        INNER JOIN rollup_provider AS rp
               ON rp.address_id = agg.signer_id AND rp.namespace_id = 0
    ), board AS (
        SELECT
            rollup_id,
            sum(size)        AS size,
            sum(blobs_count) AS blobs_count,
            max(last_time)   AS last_time,
            min(first_time)  AS first_time,
            sum(fee)         AS fee
        FROM matched
        GROUP BY 1
    )
    SELECT
        board.size,
        board.blobs_count,
        board.last_time,
        board.first_time,
        board.fee,
        board.size        / sum(board.size)        OVER () AS size_pct,
        board.fee         / sum(board.fee)         OVER () AS fee_pct,
        board.blobs_count / sum(board.blobs_count) OVER () AS blobs_count_pct,
        (now() - board.last_time < INTERVAL '1 month') AS is_active,
        rollups.*
    FROM board
    INNER JOIN rollups ON rollups.id = board.rollup_id;

CALL add_job_refresh_materialized_view();

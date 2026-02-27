CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard_day AS
WITH data AS (
    SELECT * FROM blob_log WHERE time > now() - '1 day'::interval
),
rollup_data AS (
    SELECT data.*, rp.rollup_id FROM data
    INNER JOIN rollup_provider rp
        ON rp.namespace_id = data.namespace_id AND rp.address_id = data.signer_id
    UNION ALL
    SELECT data.*, rp.rollup_id FROM data
    INNER JOIN rollup_provider rp
        ON rp.namespace_id = data.namespace_id AND rp.address_id = 0
    UNION ALL
    SELECT data.*, rp.rollup_id FROM data
    INNER JOIN rollup_provider rp
        ON rp.address_id = data.signer_id AND rp.namespace_id = 0
)
SELECT
    avg(size)                    AS avg_size,
    count(*)                     AS blobs_count,
    sum(size)                    AS total_size,
    sum(rollup_data.fee)         AS total_fee,
    ceil(sum(size) / 86400)      AS throughput,
    count(DISTINCT rollup_data.namespace_id) AS namespace_count,
    count(DISTINCT rollup_data.msg_id)       AS pfb_count,
    (CASE WHEN sum(size) > 0
        THEN ceil(sum(rollup_data.fee) * 1024 * 1024 / sum(size))
        ELSE 0 END)              AS mb_price,
    rollup_id
FROM rollup_data
GROUP BY rollup_id;

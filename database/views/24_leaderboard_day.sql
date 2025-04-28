CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard_day AS
   with 
	data as (
		select * from blob_log where time > now() - '1 day'::interval
	), 
	rollup_data as (
		select data.*, rp.rollup_id from data 
		inner join rollup_provider rp on (rp.namespace_id = 0 or rp.namespace_id = data.namespace_id) and (rp.address_id = data.signer_id OR rp.address_id = 0)
    )
    select 
        avg(size) as avg_size, 
        count(*) as blobs_count,
        sum(size) as total_size,
        sum(rollup_data.fee) as total_fee,
        ceil(sum(size) / (60*60*24)) as throughput,
        count(DISTINCT rollup_data.namespace_id) as namespace_count,
        count(DISTINCT rollup_data.msg_id) as pfb_count,
        (case when sum(size) > 0 then ceil(sum(rollup_data.fee) * 1024 * 1024 / sum(size)) else 0 end) as mb_price,
        rollup_id
    from rollup_data
    group by rollup_id;

CALL add_job_refresh_materialized_view();
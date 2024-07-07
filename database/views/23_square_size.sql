CREATE MATERIALIZED VIEW IF NOT EXISTS square_size
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 day'::interval, time) AS ts,
        square_size,
        count(*) as count_blocks
    from block_stats
    where square_size > 0
    group by 1, 2
	order by 1 desc, 2 desc
	with no data;

CALL add_view_refresh_job('square_size', NULL, INTERVAL '1 hour');

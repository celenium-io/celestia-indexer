CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_day
WITH (timescaledb.continuous, timescaledb.materialized_only=true) AS
	select 
		time_bucket('1 day'::interval, nm.ts) AS ts,
		nm.namespace_id,
		sum(pfb_count) as pfb_count,
		sum(size) as size		
	from namespace_stats_by_hour as nm
	group by 1, 2
	order by 1 desc;

CALL add_view_refresh_job('namespace_stats_by_day', INTERVAL '1 minute', INTERVAL '1 hour');

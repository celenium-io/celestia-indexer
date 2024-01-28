CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_week
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 week'::interval, nm.ts) AS ts,
		nm.namespace_id,
		sum(pfb_count) as pfb_count,
		sum(size) as size		
	from namespace_stats_by_day as nm
	group by 1, 2
	order by 1 desc
	with no data;

CALL add_view_refresh_job('namespace_stats_by_week', NULL, INTERVAL '1 hour');
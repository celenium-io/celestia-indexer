CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, nm.time) AS ts,
		nm.namespace_id,
		count(*) filter (where source = 'pfb') as pfb_count,
		count(*) filter (where source = 'fibre') as pff_count,
		sum(size) as size		
	from namespace_message as nm
	group by 1, 2
	order by 1 desc
	with no data;

CALL add_view_refresh_job('namespace_stats_by_hour', NULL, INTERVAL '1 minutes');

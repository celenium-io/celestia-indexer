CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=true) AS
	select 
		time_bucket('1 hour'::interval, nm.time) AS ts,
		nm.namespace_id,
		count(*) as pfb_count,
		sum(size) as size		
	from namespace_message as nm
	group by 1, 2
	order by 1 desc;


SELECT add_continuous_aggregate_policy('namespace_stats_by_hour',
  start_offset => NULL,
  end_offset => INTERVAL '1 minute',
  schedule_interval => INTERVAL '15 minute',
  if_not_exists => true)
WHERE NOT (SELECT EXISTS (
    SELECT FROM 
        "_timescaledb_catalog".continuous_agg
    WHERE user_view_schema = 'public' AND user_view_name = 'namespace_stats_by_hour'
    )
);

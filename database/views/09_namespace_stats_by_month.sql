CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_month
  WITH (timescaledb.continuous, timescaledb.materialized_only=true) AS
    select 
      time_bucket('1 month', nm.ts) AS ts,
      nm.namespace_id,
      count(*) as pfb_count,
      sum(size) as size		
    from namespace_stats_by_day as nm
    group by 1, 2
    order by 1 desc;

CALL add_view_refresh_job('namespace_stats_by_month', INTERVAL '10 minute', INTERVAL '1 hour');

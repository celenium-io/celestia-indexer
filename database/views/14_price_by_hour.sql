CREATE MATERIALIZED VIEW IF NOT EXISTS price_by_hour
WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
	select 
		time_bucket('1 hour'::interval, price.time) AS time,
		first(price.open, price.time) as open,
        max(high) as high,
        min(low) as low,
        last(price.close, price.time) as close
	from price
	group by 1
	order by 1 desc
	with no data;

CALL add_view_refresh_job('price_by_hour', NULL, INTERVAL '1 minute');
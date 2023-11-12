CREATE MATERIALIZED VIEW IF NOT EXISTS gas_price_candlesticks_hourly
WITH (timescaledb.continuous) AS
	select 
	 time_bucket('1 hour'::interval, time) AS timestamp,
	 candlestick_agg("time", fee/gas_wanted, gas_wanted) as value,
	 sum(fee) as fee,
	 sum(gas_used) as gas_used
	from tx
	where gas_wanted  > 0 and "status" = 'success'
	group by timestamp
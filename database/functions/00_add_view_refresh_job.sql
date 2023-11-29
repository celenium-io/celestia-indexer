CREATE OR REPLACE PROCEDURE add_view_refresh_job(view_name text, end_offset interval, schedule interval) AS
$$
declare
	mat_id text;
	id bigint;
begin	
	select mat_hypertable_id::text into mat_id from "_timescaledb_catalog".continuous_agg where user_view_name = view_name;

	if not exists (select from timescaledb_information.jobs where hypertable_name = '_materialized_hypertable_' || mat_id)
	then
		SELECT add_continuous_aggregate_policy(view_name,
			  start_offset => NULL,
			  end_offset => end_offset,
			  schedule_interval => schedule,
			  if_not_exists => true) INTO id;
	end if;
end;
$$ LANGUAGE plpgsql;
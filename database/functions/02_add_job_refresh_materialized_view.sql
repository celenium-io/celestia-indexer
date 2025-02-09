CREATE OR REPLACE PROCEDURE add_job_refresh_materialized_view()
    LANGUAGE PLPGSQL AS
    $$
    BEGIN
        if not exists (select from timescaledb_information.jobs where proc_name = 'refresh_materialized_view')
        then
            PERFORM add_job('refresh_materialized_view', '1h', config => NULL);
        end if;

        if not exists (select from timescaledb_information.jobs where proc_name = 'refresh_short_materialized_view')
        then
            PERFORM add_job('refresh_short_materialized_view', '5m', config => NULL);
        end if;

        if not exists (select from timescaledb_information.jobs where proc_name = 'refresh_da_change_materialized_view')
        then
            PERFORM add_job('refresh_da_change_materialized_view', '5m', config => NULL);
        end if;
    END
    $$;
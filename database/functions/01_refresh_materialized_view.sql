
CREATE OR REPLACE PROCEDURE refresh_materialized_view(job_id INT, config JSONB)
    LANGUAGE PLPGSQL AS
    $$
    BEGIN
        REFRESH MATERIALIZED VIEW leaderboard;
    END
    $$;


CREATE OR REPLACE PROCEDURE refresh_short_materialized_view(job_id INT, config JSONB)
    LANGUAGE PLPGSQL AS
    $$
    BEGIN
        REFRESH MATERIALIZED VIEW leaderboard_day;
    END
    $$;

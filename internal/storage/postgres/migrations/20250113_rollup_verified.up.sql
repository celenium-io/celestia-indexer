ALTER TABLE public."rollup" ADD COLUMN IF NOT EXISTS verified bool NULL;

--bun:split

COMMENT ON COLUMN public."rollup".verified IS 'Flag is set when rollup was approved';

--bun:split

UPDATE rollup set verified = TRUE where id > 0;

--bun:split

REFRESH MATERIALIZED VIEW leaderboard;

--bun:split

REFRESH MATERIALIZED VIEW leaderboard_day;
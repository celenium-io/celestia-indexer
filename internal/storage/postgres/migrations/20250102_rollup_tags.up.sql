ALTER TABLE public."rollup" ADD COLUMN IF NOT EXISTS tags _varchar NULL;

--bun:split

COMMENT ON COLUMN public."rollup".tags IS 'Rollup tags';

--bun:split

UPDATE rollup SET tags = ARRAY[category] WHERE category is not NULL;

--bun:split

REFRESH MATERIALIZED VIEW leaderboard;

--bun:split

REFRESH MATERIALIZED VIEW leaderboard_day;
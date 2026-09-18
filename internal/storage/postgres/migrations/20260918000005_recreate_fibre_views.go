// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upRecreateFibreViews, downRecreateFibreViews)
}

// execAll runs each statement on its own: continuous aggregates and DDL like
// CREATE MATERIALIZED VIEW cannot always share a multi-statement exec, and
// running them one by one lets IF NOT EXISTS/IF EXISTS make every step safe
// to retry after a partial failure.
func execAll(ctx context.Context, db *bun.DB, statements ...string) error {
	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return errors.Wrapf(err, "%.60s", stmt)
		}
	}
	return nil
}

// namespace_stats_by_hour is the base of a materialized-view chain (day <-
// year/month/week), and rollup_stats_by_hour is the base of another
// (day <- month <- leaderboard) that also feeds da_change. Dropping a base
// view CASCADEs onto every view built on top of it, fibre-aware or not, so
// da_change (its own definition is unchanged) has to be recreated here too.
func upRecreateFibreViews(ctx context.Context, db *bun.DB) error {
	if err := execAll(ctx, db,
		`DROP MATERIALIZED VIEW IF EXISTS namespace_stats_by_hour CASCADE`,
		`DROP MATERIALIZED VIEW IF EXISTS rollup_stats_by_hour CASCADE`,
		`DROP MATERIALIZED VIEW IF EXISTS leaderboard_day`,
	); err != nil {
		return errors.Wrap(err, "dropping fibre-affected views")
	}

	return execAll(ctx, db,
		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_hour
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
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_hour', NULL, INTERVAL '1 minutes')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_day
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 day'::interval, nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(pff_count) as pff_count,
				sum(size) as size
			from namespace_stats_by_hour as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_day', NULL, INTERVAL '5 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_year
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 year', nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(pff_count) as pff_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_year', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_month
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 month', nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(pff_count) as pff_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_month', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_week
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 week'::interval, nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(pff_count) as pff_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_week', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_hour
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 hour'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				count(*) filter (where source = 'pfb') as blobs_count,
				count(*) filter (where source = 'fibre') as fibre_blobs_count,
				max(logs.time) as last_time,
				min(logs.time) as first_time,
				sum(logs.fee) as fee
			from blob_log as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_hour', NULL, INTERVAL '1 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_day
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 day'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				sum(logs.blobs_count) as blobs_count,
				sum(logs.fibre_blobs_count) as fibre_blobs_count,
				max(logs.last_time) as last_time,
				min(logs.first_time) as first_time,
				sum(fee) as fee
			from rollup_stats_by_hour as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_day', NULL, INTERVAL '5 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_month
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 month'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				sum(logs.blobs_count) as blobs_count,
				sum(logs.fibre_blobs_count) as fibre_blobs_count,
				max(logs.last_time) as last_time,
				min(logs.first_time) as first_time,
				sum(fee) as fee
			from rollup_stats_by_day as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_month', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
			WITH rollups AS (
				SELECT * FROM rollup WHERE verified = TRUE
			), agg AS MATERIALIZED (
				SELECT
					namespace_id,
					signer_id,
					sum(size)              AS size,
					sum(blobs_count)       AS blobs_count,
					sum(fibre_blobs_count) AS fibre_blobs_count,
					max(last_time)   AS last_time,
					min(first_time)  AS first_time,
					sum(fee)         AS fee
				FROM rollup_stats_by_month
				GROUP BY 1, 2
			), matched AS (
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.fibre_blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = agg.signer_id AND rp.namespace_id = agg.namespace_id
				UNION ALL
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.fibre_blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = 0 AND rp.namespace_id = agg.namespace_id
				UNION ALL
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.fibre_blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = agg.signer_id AND rp.namespace_id = 0
			), board AS (
				SELECT
					rollup_id,
					sum(size)              AS size,
					sum(blobs_count)       AS blobs_count,
					sum(fibre_blobs_count) AS fibre_blobs_count,
					max(last_time)   AS last_time,
					min(first_time)  AS first_time,
					sum(fee)         AS fee
				FROM matched
				GROUP BY 1
			)
			SELECT
				board.size,
				board.blobs_count,
				board.fibre_blobs_count,
				board.last_time,
				board.first_time,
				board.fee,
				board.size        / sum(board.size)        OVER () AS size_pct,
				board.fee         / sum(board.fee)         OVER () AS fee_pct,
				(board.blobs_count + board.fibre_blobs_count) / sum(board.blobs_count + board.fibre_blobs_count) OVER () AS blobs_count_pct,
				(now() - board.last_time < INTERVAL '1 month') AS is_active,
				rollups.*
			FROM board
			INNER JOIN rollups ON rollups.id = board.rollup_id`,
		`CALL add_job_refresh_materialized_view()`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS da_change AS
			with board1 as (
				select
					rollup_id,
					sum(size) as size
				from (
					select
						namespace_id,
						signer_id,
						sum(size) as size
					from rollup_stats_by_hour
					where time > now() - '1 week'::interval
					group by 1, 2
				) as agg
				inner join rollup_provider as rp on (rp.address_id = agg.signer_id OR rp.address_id = 0) AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
				inner join rollup on rollup.id = rp.rollup_id
				where rollup.verified = TRUE
				group by 1
			), board2 as (
				select
					rollup_id,
					sum(size) as size
				from (
					select
						namespace_id,
						signer_id,
						sum(size) as size
					from rollup_stats_by_hour
					where time <= now() - '1 week'::interval and time > now() - '2 week'::interval
					group by 1, 2
				) as agg
				inner join rollup_provider as rp on (rp.address_id = agg.signer_id OR rp.address_id = 0) AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
				inner join rollup on rollup.id = rp.rollup_id
				where rollup.verified = TRUE
				group by 1
			)
			select
				case
					when coalesce(board2.size, 0) > 0
						then coalesce(board1.size, 0) / coalesce(board2.size, 0) - 1
					when coalesce(board1.size, 0) > 0 and coalesce(board2.size, 0) = 0
						then 1
					else 0
					end as da_pct,
				rollup.id as rollup_id
			from rollup
			left join board1 on rollup.id = board1.rollup_id
			left join board2 on rollup.id = board2.rollup_id`,
		`CALL add_job_refresh_materialized_view()`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard_day AS
			WITH data AS (
				SELECT * FROM blob_log WHERE time > now() - '1 day'::interval
			),
			rollup_data AS (
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.namespace_id = data.namespace_id AND rp.address_id = data.signer_id
				UNION ALL
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.namespace_id = data.namespace_id AND rp.address_id = 0
				UNION ALL
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.address_id = data.signer_id AND rp.namespace_id = 0
			)
			SELECT
				avg(size)                    AS avg_size,
				count(*) filter (where source = 'pfb')   AS blobs_count,
				count(*) filter (where source = 'fibre') AS fibre_blobs_count,
				sum(size)                    AS total_size,
				sum(rollup_data.fee)         AS total_fee,
				ceil(sum(size) / 86400)      AS throughput,
				count(DISTINCT rollup_data.namespace_id) AS namespace_count,
				count(DISTINCT rollup_data.msg_id) filter (where source = 'pfb')  AS pfb_count,
				count(DISTINCT rollup_data.msg_id) filter (where source = 'fibre')  AS pff_count,
				(CASE WHEN sum(size) > 0
					THEN ceil(sum(rollup_data.fee) * 1024 * 1024 / sum(size))
					ELSE 0 END)              AS mb_price,
				rollup_id
			FROM rollup_data
			GROUP BY rollup_id`,
	)
}

func downRecreateFibreViews(ctx context.Context, db *bun.DB) error {
	if err := execAll(ctx, db,
		`DROP MATERIALIZED VIEW IF EXISTS namespace_stats_by_hour CASCADE`,
		`DROP MATERIALIZED VIEW IF EXISTS rollup_stats_by_hour CASCADE`,
		`DROP MATERIALIZED VIEW IF EXISTS leaderboard_day`,
	); err != nil {
		return errors.Wrap(err, "dropping fibre-affected views")
	}

	return execAll(ctx, db,
		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_hour
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 hour'::interval, nm.time) AS ts,
				nm.namespace_id,
				count(*) as pfb_count,
				sum(size) as size
			from namespace_message as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_hour', NULL, INTERVAL '1 minutes')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_day
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 day'::interval, nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(size) as size
			from namespace_stats_by_hour as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_day', NULL, INTERVAL '5 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_year
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 year', nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_year', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_month
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 month', nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_month', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS namespace_stats_by_week
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 week'::interval, nm.ts) AS ts,
				nm.namespace_id,
				sum(pfb_count) as pfb_count,
				sum(size) as size
			from namespace_stats_by_day as nm
			group by 1, 2
			order by 1 desc
			with no data`,
		`CALL add_view_refresh_job('namespace_stats_by_week', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_hour
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 hour'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				count(*) as blobs_count,
				max(logs.time) as last_time,
				min(logs.time) as first_time,
				sum(logs.fee) as fee
			from blob_log as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_hour', NULL, INTERVAL '1 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_day
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 day'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				sum(logs.blobs_count) as blobs_count,
				max(logs.last_time) as last_time,
				min(logs.first_time) as first_time,
				sum(fee) as fee
			from rollup_stats_by_hour as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_day', NULL, INTERVAL '5 minute')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS rollup_stats_by_month
		WITH (timescaledb.continuous, timescaledb.materialized_only=false) AS
			select
				time_bucket('1 month'::interval, logs.time) AS time,
				logs.namespace_id,
				logs.signer_id,
				sum(logs.size) as size,
				sum(logs.blobs_count) as blobs_count,
				max(logs.last_time) as last_time,
				min(logs.first_time) as first_time,
				sum(fee) as fee
			from rollup_stats_by_day as logs
			group by 1, 2, 3
			with no data`,
		`CALL add_view_refresh_job('rollup_stats_by_month', NULL, INTERVAL '1 hour')`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard AS
			WITH rollups AS (
				SELECT * FROM rollup WHERE verified = TRUE
			), agg AS MATERIALIZED (
				SELECT
					namespace_id,
					signer_id,
					sum(size)        AS size,
					sum(blobs_count) AS blobs_count,
					max(last_time)   AS last_time,
					min(first_time)  AS first_time,
					sum(fee)         AS fee
				FROM rollup_stats_by_month
				GROUP BY 1, 2
			), matched AS (
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = agg.signer_id AND rp.namespace_id = agg.namespace_id
				UNION ALL
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = 0 AND rp.namespace_id = agg.namespace_id
				UNION ALL
				SELECT rp.rollup_id, agg.size, agg.blobs_count, agg.last_time, agg.first_time, agg.fee
				FROM agg
				INNER JOIN rollup_provider AS rp
				       ON rp.address_id = agg.signer_id AND rp.namespace_id = 0
			), board AS (
				SELECT
					rollup_id,
					sum(size)        AS size,
					sum(blobs_count) AS blobs_count,
					max(last_time)   AS last_time,
					min(first_time)  AS first_time,
					sum(fee)         AS fee
				FROM matched
				GROUP BY 1
			)
			SELECT
				board.size,
				board.blobs_count,
				board.last_time,
				board.first_time,
				board.fee,
				board.size        / sum(board.size)        OVER () AS size_pct,
				board.fee         / sum(board.fee)         OVER () AS fee_pct,
				board.blobs_count / sum(board.blobs_count) OVER () AS blobs_count_pct,
				(now() - board.last_time < INTERVAL '1 month') AS is_active,
				rollups.*
			FROM board
			INNER JOIN rollups ON rollups.id = board.rollup_id`,
		`CALL add_job_refresh_materialized_view()`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS da_change AS
			with board1 as (
				select
					rollup_id,
					sum(size) as size
				from (
					select
						namespace_id,
						signer_id,
						sum(size) as size
					from rollup_stats_by_hour
					where time > now() - '1 week'::interval
					group by 1, 2
				) as agg
				inner join rollup_provider as rp on (rp.address_id = agg.signer_id OR rp.address_id = 0) AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
				inner join rollup on rollup.id = rp.rollup_id
				where rollup.verified = TRUE
				group by 1
			), board2 as (
				select
					rollup_id,
					sum(size) as size
				from (
					select
						namespace_id,
						signer_id,
						sum(size) as size
					from rollup_stats_by_hour
					where time <= now() - '1 week'::interval and time > now() - '2 week'::interval
					group by 1, 2
				) as agg
				inner join rollup_provider as rp on (rp.address_id = agg.signer_id OR rp.address_id = 0) AND (rp.namespace_id = agg.namespace_id OR rp.namespace_id = 0)
				inner join rollup on rollup.id = rp.rollup_id
				where rollup.verified = TRUE
				group by 1
			)
			select
				case
					when coalesce(board2.size, 0) > 0
						then coalesce(board1.size, 0) / coalesce(board2.size, 0) - 1
					when coalesce(board1.size, 0) > 0 and coalesce(board2.size, 0) = 0
						then 1
					else 0
					end as da_pct,
				rollup.id as rollup_id
			from rollup
			left join board1 on rollup.id = board1.rollup_id
			left join board2 on rollup.id = board2.rollup_id`,
		`CALL add_job_refresh_materialized_view()`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard_day AS
			WITH data AS (
				SELECT * FROM blob_log WHERE time > now() - '1 day'::interval
			),
			rollup_data AS (
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.namespace_id = data.namespace_id AND rp.address_id = data.signer_id
				UNION ALL
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.namespace_id = data.namespace_id AND rp.address_id = 0
				UNION ALL
				SELECT data.*, rp.rollup_id FROM data
				INNER JOIN rollup_provider rp
					ON rp.address_id = data.signer_id AND rp.namespace_id = 0
			)
			SELECT
				avg(size)                    AS avg_size,
				count(*)                     AS blobs_count,
				sum(size)                    AS total_size,
				sum(rollup_data.fee)         AS total_fee,
				ceil(sum(size) / 86400)      AS throughput,
				count(DISTINCT rollup_data.namespace_id) AS namespace_count,
				count(DISTINCT rollup_data.msg_id)       AS pfb_count,
				(CASE WHEN sum(size) > 0
					THEN ceil(sum(rollup_data.fee) * 1024 * 1024 / sum(size))
					ELSE 0 END)              AS mb_price,
				rollup_id
			FROM rollup_data
			GROUP BY rollup_id`,
	)
}

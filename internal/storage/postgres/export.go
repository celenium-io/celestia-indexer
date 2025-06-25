// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/dipdup-net/go-lib/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Export struct {
	*bun.DB
}

func NewExport(cfg config.Database) *Export {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&application_name=export",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return &Export{db}
}

func (e *Export) ToCsv(ctx context.Context, writer io.Writer, query string) error {
	conn, err := e.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	rawQuery := fmt.Sprintf("COPY (%s) TO STDOUT WITH CSV HEADER", bun.Safe(query))
	_, err = pgdriver.CopyTo(ctx, conn, writer, rawQuery)
	return err
}

func (e *Export) Close() error {
	if err := e.DB.Close(); err != nil {
		return err
	}
	return nil
}

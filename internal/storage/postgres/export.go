// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"fmt"
	"io"

	"github.com/dipdup-io/go-lib/database"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Export struct {
	*database.Bun
}

func NewExport(db *database.Bun) *Export {
	return &Export{db}
}

func (e *Export) ToCsv(ctx context.Context, writer io.Writer, query string) error {
	pool := e.Pool()
	if pool == nil {
		return errors.New("pool is nil")
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "acquire connection pool")
	}
	defer conn.Release()

	rawQuery := fmt.Sprintf("COPY (%s) TO STDOUT WITH CSV HEADER", bun.Safe(query))
	_, err = conn.Conn().PgConn().CopyTo(ctx, writer, rawQuery)
	return err
}

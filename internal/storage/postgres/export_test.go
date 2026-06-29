// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (s *StorageTestSuite) TestExportToCsv() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	query := s.storage.
		Connection().DB().
		NewSelect().
		Model((*storage.Address)(nil)).
		Order("id")

	var buf bytes.Buffer
	err := s.storage.export.ToCsv(ctx, &buf, query)
	s.Require().NoError(err)

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	rows, err := reader.ReadAll()
	s.Require().NoError(err)
	s.Require().Len(rows, 6) // 1 header + 5 fixtures

	header := rows[0]
	s.Require().ElementsMatch(
		[]string{"id", "height", "last_height", "hash", "address", "name", "is_forwarding"},
		header,
	)

	col := make(map[string]int, len(header))
	for i, name := range header {
		col[name] = i
	}

	first := rows[1]
	s.Require().Equal("1", first[col["id"]])
	s.Require().Equal("100", first[col["height"]])
	s.Require().Equal("celestia1mm8yykm46ec3t0dgwls70g0jvtm055wk9ayal8", first[col["address"]])
	s.Require().Equal("f", first[col["is_forwarding"]])

	last := rows[5]
	s.Require().Equal("5", last[col["id"]])
	s.Require().Equal("t", last[col["is_forwarding"]])
}

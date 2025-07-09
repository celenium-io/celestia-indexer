// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"
)

func (s *StorageTestSuite) TestExportToCsv() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	var buf bytes.Buffer
	err := s.storage.export.ToCsv(ctx, &buf, "select * from address")
	s.Require().NoError(err)

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	rows, err := reader.ReadAll()
	s.Require().NoError(err)
	s.Require().Len(rows, 5)
}

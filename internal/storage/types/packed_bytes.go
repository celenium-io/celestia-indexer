// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"

	"github.com/andybalholm/brotli"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

var json = sonic.ConfigFastest

type PackedBytes map[string]any

var _ sql.Scanner = (*PackedBytes)(nil)

func (pb *PackedBytes) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.Errorf("invalid packed bytes type: %T", src)
	}

	result := bytes.NewBuffer(b)
	return json.NewDecoder(brotli.NewReader(result)).Decode(pb)
}

var _ driver.Valuer = (*PackedBytes)(nil)

func (pb PackedBytes) Value() (driver.Value, error) {
	return pb.ToBytes()
}

func (pb PackedBytes) ToBytes() ([]byte, error) {
	b, err := json.Marshal(pb)
	if err != nil {
		return nil, err
	}
	result := bytes.NewBuffer(nil)
	writer := brotli.NewWriterLevel(result, brotli.BestSpeed)

	if _, err := writer.Write(b); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

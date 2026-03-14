// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"

	"github.com/andybalholm/brotli"
	"github.com/bytedance/sonic"
	"github.com/celenium-io/celestia-indexer/internal/pool"
	"github.com/pkg/errors"
)

var json = sonic.ConfigFastest

var (
	brotliPool = pool.New(
		func() *brotli.Writer { return brotli.NewWriterLevel(nil, brotli.BestSpeed) },
	)
	bufPool = pool.New(
		func() *bytes.Buffer { return bytes.NewBuffer(make([]byte, 0, 512)) },
	)
)

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
	buf := bufPool.Get()
	buf.Reset()
	defer bufPool.Put(buf)

	writer := brotliPool.Get()
	writer.Reset(buf)
	defer brotliPool.Put(writer)

	if _, err := writer.Write(b); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

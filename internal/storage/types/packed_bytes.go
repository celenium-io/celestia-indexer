// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	errKeyNotFound = errors.New("key not found")
)

type PackedBytes map[string]any

var _ sql.Scanner = (*PackedBytes)(nil)

func (pb *PackedBytes) Scan(src any) error {
	if src == nil {
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return errors.Errorf("invalid packed bytes type: %T", src)
	}
	return msgpack.Unmarshal(b, pb)
}

var _ driver.Valuer = (*PackedBytes)(nil)

func (pb PackedBytes) Value() (driver.Value, error) {
	return pb.ToBytes()
}

func (pb PackedBytes) ToBytes() ([]byte, error) {
	return msgpack.Marshal(pb)
}

func (pb PackedBytes) GetString(key string) (string, error) {
	val, ok := pb[key]
	if !ok {
		return "", errors.Wrap(errKeyNotFound, key)
	}
	str, ok := val.(string)
	if !ok {
		return "", errors.Errorf("key is not a string type: %s", key)
	}
	return str, nil
}

func (pb PackedBytes) GetStringOrDefault(key string) string {
	val, ok := pb[key]
	if !ok {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

type Hex []byte

var nullBytes = []byte("null")

func HexFromString(s string) (Hex, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return Hex(data), nil
}

func (h *Hex) UnmarshalJSON(data []byte) error {
	if h == nil {
		return nil
	}

	if bytes.Equal(nullBytes, data) {
		*h = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*h = make(Hex, len(b))
	_ = copy(*h, b)
	return nil
}

func (h Hex) MarshalJSON() ([]byte, error) {
	if h == nil {
		return nullBytes, nil
	}
	return []byte(strconv.Quote(h.String())), nil
}

func (h *Hex) Scan(src interface{}) (err error) {
	switch val := src.(type) {
	case []byte:
		*h = make(Hex, len(val))
		_ = copy(*h, val)
	case nil:
		*h = make(Hex, 0)
	default:
		return errors.Errorf("unknown hex database type: %T", src)
	}
	return nil
}

var _ driver.Valuer = (*Hex)(nil)

func (h Hex) Value() (driver.Value, error) {
	return []byte(h), nil
}

func (h Hex) Bytes() []byte {
	return []byte(h)
}

func (h Hex) String() string {
	return strings.ToUpper(hex.EncodeToString([]byte(h)))
}

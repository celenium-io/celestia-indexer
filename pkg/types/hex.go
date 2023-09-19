package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Hex []byte

var nullBytes = "null"

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

	if nullBytes == string(data) {
		*h = nil
		return nil
	}
	length := len(data)
	if length%2 == 1 {
		return errors.Errorf("odd hex lenght: %d %v", length, data)
	}
	if data[0] != '"' || data[length-1] != '"' {
		return errors.Errorf("hex should be quotted string: got=%s", data)
	}

	data = bytes.Trim(data, `"`)
	*h = make(Hex, hex.DecodedLen(length-1))
	if length-1 == 0 {
		return nil
	}
	_, err := hex.Decode(*h, data)
	return err
}

func (h Hex) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte(nullBytes), nil
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

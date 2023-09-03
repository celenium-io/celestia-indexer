package types

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      int64  `json:"id"`
	JsonRpc string `json:"jsonrpc"`
}

type Response[T any] struct {
	Id      int64  `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Error   *Error `json:"error,omitempty"`
	Result  T      `json:"result"`
}

// Error -
type Error struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Error -
func (e Error) Error() string {
	return fmt.Sprintf("code=%d message=%s data=%s", e.Code, e.Message, string(e.Data))
}

package websocket

import "errors"

var (
	ErrUnknownMethod     = errors.New("unknown method")
	ErrUnknownChannel    = errors.New("unknown channel")
	ErrUnavailableFilter = errors.New("unknown filter value")
)

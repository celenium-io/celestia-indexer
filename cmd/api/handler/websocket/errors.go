package websocket

import "errors"

var (
	ErrUnknownMethod     = errors.New("unknown method")
	ErrUnknownChannel    = errors.New("unknown channel")
	ErrTimeout           = errors.New("connection timeout")
	ErrUnavailableFilter = errors.New("unknown filter value")
)

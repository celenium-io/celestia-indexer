package websocket

import "errors"

var (
	ErrUnknownMethod      = errors.New("unknown method")
	ErrUnknownChannel     = errors.New("unknown channel")
	ErrTimeout            = errors.New("connection timeout")
	ErrUnavailiableFilter = errors.New("unknown filter value")
)

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"encoding/json"
	"errors"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
)

// methods
const (
	MethodSubscribe   = "subscribe"
	MethodUnsubscribe = "unsubscribe"
)

// channels
const (
	ChannelHead     = "head"
	ChannelBlocks   = "blocks"
	ChannelGasPrice = "gas_price"
	ChannelError    = "error"
)

type Message struct {
	Method string          `json:"method" validate:"required,oneof=subscribe unsubscribe"`
	Body   json.RawMessage `json:"body"   validate:"required"`
}

type Subscribe struct {
	Channel string          `json:"channel" validate:"required,oneof=head blocks gas_price"`
	Filters json.RawMessage `json:"filters" validate:"required"`
}

type Unsubscribe struct {
	Channel string `json:"channel" validate:"required,oneof=head blocks gas_price"`
}

type TransactionFilters struct {
	Status   []string `json:"status,omitempty"`
	Messages []string `json:"msg_type,omitempty"`
}

type INotification interface {
	*responses.Block | *responses.State | *responses.GasPrice
}

type Notification[T INotification] struct {
	Channel string `json:"channel"`
	Body    T      `json:"body"`
}

func (n Notification[T]) GetChannel() string {
	return n.Channel
}

func NewBlockNotification(block responses.Block) Notification[*responses.Block] {
	return Notification[*responses.Block]{
		Channel: ChannelBlocks,
		Body:    &block,
	}
}

func NewStateNotification(state responses.State) Notification[*responses.State] {
	return Notification[*responses.State]{
		Channel: ChannelHead,
		Body:    &state,
	}
}

func NewGasPriceNotification(value responses.GasPrice) Notification[*responses.GasPrice] {
	return Notification[*responses.GasPrice]{
		Channel: ChannelGasPrice,
		Body:    &value,
	}
}

// error codes reported to the client. Codes are stable and safe to expose;
// internal error details are never sent to the client to avoid leaking
// sensitive information.
const (
	ErrCodeInvalidMessage = 1
	ErrCodeUnknownMethod  = 2
	ErrCodeUnknownChannel = 3
)

// ErrorMessage is sent to the client when an incoming message cannot be handled.
type ErrorMessage struct {
	Channel string    `json:"channel"`
	Body    ErrorBody `json:"body"`
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (m ErrorMessage) GetChannel() string {
	return m.Channel
}

func newErrorMessage(code int, message string) ErrorMessage {
	return ErrorMessage{
		Channel: ChannelError,
		Body: ErrorBody{
			Code:    code,
			Message: message,
		},
	}
}

// predefined client-facing errors
var (
	errInvalidMessage = newErrorMessage(ErrCodeInvalidMessage, "invalid message")
	errUnknownMethod  = newErrorMessage(ErrCodeUnknownMethod, "unknown method")
	errUnknownChannel = newErrorMessage(ErrCodeUnknownChannel, "unknown channel")
)

// errorMessage maps an internal error to a safe, predefined client-facing message.
// Unmapped errors fall back to a generic "invalid message" so that internal
// details are never exposed.
func errorMessage(err error) ErrorMessage {
	switch {
	case errors.Is(err, ErrUnknownMethod):
		return errUnknownMethod
	case errors.Is(err, ErrUnknownChannel):
		return errUnknownChannel
	default:
		return errInvalidMessage
	}
}

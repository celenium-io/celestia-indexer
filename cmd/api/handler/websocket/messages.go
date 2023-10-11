// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import "github.com/goccy/go-json"

// methods
const (
	MethodSubscribe   = "subscribe"
	MethodUnsubscribe = "unsubscribe"
)

// channels
const (
	ChannelHead = "head"
	ChannelTx   = "tx"
)

type Message struct {
	Method string          `json:"method" validate:"required,oneof=subscribe,unsubscribe"`
	Body   json.RawMessage `json:"body"   validate:"required"`
}

type Subscribe struct {
	Channel string          `json:"channel" validate:"required,oneof=head tx"`
	Filters json.RawMessage `json:"filters" validate:"required"`
}

type Unsubscribe struct {
	Channel string `json:"channel" validate:"required,oneof=head tx"`
}

type TransactionFilters struct {
	Status   []string `json:"status,omitempty"`
	Messages []string `json:"msg_type,omitempty"`
}

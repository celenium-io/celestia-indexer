// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/goccy/go-json"
)

// methods
const (
	MethodSubscribe   = "subscribe"
	MethodUnsubscribe = "unsubscribe"
)

// channels
const (
	ChannelHead   = "head"
	ChannelBlocks = "blocks"
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

type INotification interface {
	*responses.Block | *responses.State
}

type Notification[T INotification] struct {
	Channel string `json:"channel"`
	Body    T      `json:"body"`
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

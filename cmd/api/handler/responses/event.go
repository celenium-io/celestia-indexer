// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
)

type Event struct {
	Id       uint64         `example:"321"                       format:"int64"     json:"id"              swaggertype:"integer"`
	Height   pkgTypes.Level `example:"100"                       format:"int64"     json:"height"          swaggertype:"integer"`
	Time     time.Time      `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggertype:"string"`
	Position int64          `example:"1"                         format:"int64"     json:"position"        swaggertype:"integer"`
	TxId     uint64         `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggertype:"integer"`

	Type types.EventType `example:"commission" json:"type"`

	Data map[string]any `json:"data"`
}

func NewEvent(event storage.Event) Event {
	result := Event{
		Id:       event.Id,
		Height:   event.Height,
		Time:     event.Time,
		Position: event.Position,
		Type:     event.Type,
		Data:     event.Data,
	}

	if event.TxId != nil {
		result.TxId = *event.TxId
	}

	return result
}

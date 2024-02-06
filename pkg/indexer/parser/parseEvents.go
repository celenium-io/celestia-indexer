// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/rs/zerolog/log"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func parseEvents(b types.BlockData, events []types.Event) []storage.Event {
	result := make([]storage.Event, len(events))

	for i := range events {
		parseEvent(b, events[i], i, &result[i])
	}

	return result
}

func parseEvent(b types.BlockData, eN types.Event, index int, resultEvent *storage.Event) {
	eventType, err := storageTypes.ParseEventType(eN.Type)
	if err != nil {
		log.Err(err).Msgf("got type %v", eN.Type)
		eventType = storageTypes.EventTypeUnknown
	}

	resultEvent.Height = b.Height
	resultEvent.Time = b.Block.Time
	resultEvent.Position = int64(index)
	resultEvent.Type = eventType
	resultEvent.Data = make(map[string]any, len(eN.Attributes))

	for i := range eN.Attributes {
		resultEvent.Data[string(eN.Attributes[i].Key)] = string(eN.Attributes[i].Value)
	}
}

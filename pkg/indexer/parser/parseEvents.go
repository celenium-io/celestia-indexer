package parser

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func parseEvents(b types.BlockData, events []types.Event) []storage.Event {
	result := make([]storage.Event, len(events))

	for i, eN := range events {
		eS := parseEvent(b, eN, i)
		result[i] = eS
	}

	return result
}

func parseEvent(b types.BlockData, eN types.Event, index int) storage.Event {
	eventType, err := storageTypes.ParseEventType(eN.Type)
	if err != nil {
		log.Err(err).Msgf("got type %v", eN.Type)
		eventType = storageTypes.EventTypeUnknown
	}

	event := storage.Event{
		Height:   b.Height,
		Time:     b.Block.Time,
		Position: uint64(index),
		Type:     eventType,
		Data:     make(map[string]any),
	}

	for _, attr := range eN.Attributes {
		key := decodeEventAttribute(attr.Key)
		value := decodeEventAttribute(attr.Value)
		event.Data[key] = value
	}

	return event
}

var b64 = base64.StdEncoding

func decodeEventAttribute(data string) string {
	dst, err := b64.DecodeString(data)
	if err != nil {
		return data
	}
	return string(dst)
}

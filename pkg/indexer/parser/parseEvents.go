package parser

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
)

func parseEvents(b types.BlockData, events []nodeTypes.Event) []storage.Event {
	result := make([]storage.Event, len(events))

	for i, eN := range events {
		eS := parseEvent(b, eN, i)
		result[i] = eS
	}

	return result
}

func parseEvent(b types.BlockData, eN nodeTypes.Event, index int) storage.Event {
	event := storage.Event{
		Height:   b.Height,
		Time:     b.Block.Time,
		Position: uint64(index),
		Type:     storageTypes.EventType(eN.Type), // TODO errors
		Data:     make(map[string]any),
	}

	for _, attr := range eN.Attributes {
		event.Data[string(attr.Key)] = attr.Value // TODO: create actual unmarshalling bytes
	}

	return event
}

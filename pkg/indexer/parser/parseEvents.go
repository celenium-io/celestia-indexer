// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package parser

import (
	"github.com/rs/zerolog/log"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func parseEvents(ctx *context.Context, b types.BlockData, events []types.Event) ([]storage.Event, error) {
	result := make([]storage.Event, len(events))

	for i := range events {
		if err := parseEvent(ctx, b, events[i], i, &result[i], false); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func parseBlockEvents(ctx *context.Context, b types.BlockData, events []types.Event, firstTxEvent *types.Event) ([]storage.Event, error) {
	result := make([]storage.Event, len(events))

	var (
		count        = 0
		isDuplicated bool
	)

	for i := range events {
		isDuplicated = isDuplicated || (firstTxEvent != nil && firstTxEvent.Compare(events[i]))
		if isDuplicated {
			count++
		}
		isDuplicated = isDuplicated && count <= ctx.TxEventsCount

		if err := parseEvent(ctx, b, events[i], i, &result[i], isDuplicated); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func parseEvent(ctx *context.Context, b types.BlockData, eN types.Event, index int, resultEvent *storage.Event, duplicated bool) error {
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
		resultEvent.Data[eN.Attributes[i].Key] = eN.Attributes[i].Value
	}

	if duplicated {
		return nil
	}
	return processEvent(ctx, resultEvent)
}

func processEvent(ctx *context.Context, event *storage.Event) error {
	switch event.Type {
	case storageTypes.EventTypeBurn:
		ctx.SubSupply(event.Data)
	case storageTypes.EventTypeMint:
		ctx.SetInflation(event.Data)
		ctx.AddSupply(event.Data)
	case storageTypes.EventTypeCoinReceived:
		return parseCoinReceived(ctx, event.Data, event.Height)
	case storageTypes.EventTypeCoinSpent:
		return parseCoinSpent(ctx, event.Data, event.Height)
	case storageTypes.EventTypeCompleteUnbonding:
		return parseCompleteUnbonding(ctx, event.Data, event.Height)
	case storageTypes.EventTypeCommission:
		return parseCommission(ctx, event.Data)
	case storageTypes.EventTypeRewards:
		return parseRewards(ctx, event.Data)
	case storageTypes.EventTypeProposerReward:
		return parseRewards(ctx, event.Data)
	case storageTypes.EventTypeSlash:
		return parseSlash(ctx, event.Data)
	case storageTypes.EventTypeActiveProposal:
		return parseProposal(ctx, event.Data)
	case storageTypes.EventTypeInactiveProposal:
		return parseProposal(ctx, event.Data)
	case storageTypes.EventTypeHyperlanecorepostDispatchv1EventCreateIgp:
		return parseCreateIgp(ctx, event.Data)
	case storageTypes.EventTypeHyperlanecorepostDispatchv1EventSetIgp:
		return parseSetIgp(ctx, event.Data)
	case storageTypes.EventTypeHyperlanecorepostDispatchv1EventSetDestinationGasConfig:
		return parseSetDestinationGasConfig(ctx, event.Data)
	}

	return nil
}

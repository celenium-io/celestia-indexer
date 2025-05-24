// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package parser

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func parseEvents(ctx *context.Context, b types.BlockData, events []types.Event) ([]storage.Event, error) {
	result := make([]storage.Event, len(events))

	for i := range events {
		if err := parseEvent(ctx, b, events[i], i, &result[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func parseEvent(ctx *context.Context, b types.BlockData, eN types.Event, index int, resultEvent *storage.Event) error {
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
		if b.AppVersion <= 3 {
			key, err := base64.StdEncoding.DecodeString(eN.Attributes[i].Key)
			if err != nil {
				return errors.Wrapf(err, "decode event attribute key: %s appversion=%d", eN.Attributes[i].Key, b.AppVersion)
			}
			value, err := base64.StdEncoding.DecodeString(eN.Attributes[i].Value)
			if err != nil {
				return errors.Wrapf(err, "decode event attribute key: %s appversion=%d", eN.Attributes[i].Key, b.AppVersion)
			}
			resultEvent.Data[string(key)] = string(value)
		} else {
			resultEvent.Data[eN.Attributes[i].Key] = eN.Attributes[i].Value
		}
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
	}

	return nil
}

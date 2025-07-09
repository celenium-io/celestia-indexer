// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleUpdateClient(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/ibc.core.client.v1.MsgUpdateClient" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processUpdateClient(ctx, events, msg, idx)
}

func processUpdateClient(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	uc, err := decode.NewUpdateClient(events[*idx].Data)
	if err != nil {
		return errors.Wrap(err, "parse update client event")
	}

	header, err := decoder.HeaderFromMap(msg.Data, "Header")
	if err != nil {
		return errors.Wrap(err, "receiving Header from message")
	}

	msg.IbcClient = &storage.IbcClient{
		Id:                   uc.Id,
		UpdatedAt:            msg.Time,
		ChainId:              header.Header.ChainID,
		LatestRevisionHeight: uc.ConsensusHeight,
		LatestRevisionNumber: uc.Revision,
	}

	*idx += 2
	return nil
}

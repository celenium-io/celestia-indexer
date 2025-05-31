// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleUnjail(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.slashing.v1beta1.MsgUnjail" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processUnjail(ctx, events, msg, idx)
}

func processUnjail(ctx *context.Context, events []storage.Event, _ *storage.Message, idx *int) error {
	if events[*idx].Type != types.EventTypeMessage {
		return errors.Errorf("slashing unexpected event type: %s", events[*idx].Type)
	}

	module := decoder.StringFromMap(events[*idx].Data, "module")
	if module != types.ModuleNameSlashing.String() {
		return errors.Errorf("slashing unexpected module name: %s", module)
	}

	sender := decoder.StringFromMap(events[*idx].Data, "sender")
	if sender == "" {
		return errors.Errorf("slashing unexpected sender value: %s", sender)
	}

	jailed := false
	v := storage.EmptyValidator()
	v.Address = sender
	v.Jailed = &jailed
	ctx.AddValidator(v)

	*idx += 1
	return nil
}

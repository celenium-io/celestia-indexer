// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
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
	return processUnjail(ctx, events, msg, idx)
}

func processUnjail(ctx *context.Context, events []storage.Event, _ *storage.Message, idx *int) error {
	if events[*idx].Type != types.EventTypeMessage {
		return errors.Errorf("slashing unexpected event type: %s", events[*idx].Type)
	}

	module := decoder.StringFromMap(events[*idx].Data, "module")
	if module == "" {
		*idx += 1
		module = decoder.StringFromMap(events[*idx].Data, "module")
	}
	if module != types.ModuleNameSlashing.String() {
		return errors.Errorf("slashing unexpected module name: %s", module)
	}

	sender := decoder.StringFromMap(events[*idx].Data, "sender")
	if sender == "" {
		return errors.Errorf("slashing unexpected sender value: %s", sender)
	}

	prefix, hash, err := pkgTypes.Address(sender).Decode()
	if err != nil {
		return errors.Wrap(err, "parsing validator address")
	}

	jailed := false
	v := storage.EmptyValidator()

	if prefix == pkgTypes.AddressPrefixValoper {
		v.Address = sender
	} else {
		addr, err := pkgTypes.NewValoperAddressFromBytes(hash)
		if err != nil {
			return errors.Wrap(err, "encoding validator address")
		}
		v.Address = addr.String()
	}

	v.Jailed = &jailed
	ctx.AddValidator(v)

	*idx += 1
	return nil
}

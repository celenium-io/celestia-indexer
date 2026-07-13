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

func handleUnjail(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/cosmos.slashing.v1beta1.MsgUnjail" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	return processUnjail(ctx, c, msg)
}

func processUnjail(ctx *context.Context, c *Cursor, _ *storage.Message) error {
	event, ok := c.Next()
	if !ok {
		return errors.New("not enough events for unjail")
	}
	if event.Type != types.EventTypeMessage {
		return errors.Errorf("slashing unexpected event type: %s", event.Type)
	}

	module := decoder.StringFromMap(event.Data, "module")
	if module == "" {
		event, ok = c.Next()
		if !ok {
			return errors.New("not enough events for unjail after parsing module name")
		}
		module = decoder.StringFromMap(event.Data, "module")
	}
	if module != types.ModuleNameSlashing.String() {
		return errors.Errorf("slashing unexpected module name: %s", module)
	}

	sender := decoder.StringFromMap(event.Data, "sender")
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

	return nil
}

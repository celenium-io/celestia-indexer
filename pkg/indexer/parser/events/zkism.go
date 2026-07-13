// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleCreateZkISM(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/celestia.zkism.v1.MsgCreateInterchainSecurityModule" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processCreateZkISM(ctx, c, msg)
}

func processCreateZkISM(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule {
			continue
		}

		e, err := decode.NewZkISMCreateEvent(event.Data)
		if err != nil {
			return errors.Wrap(err, "parsing create zkism event")
		}

		ism := &storage.ZkISM{
			Height:              ctx.Block.Height,
			Time:                ctx.Block.Time,
			ExternalId:          e.Id,
			State:               e.State,
			MerkleTreeAddress:   e.MerkleTreeAddress,
			Groth16VKey:         e.Groth16VKey,
			StateTransitionVKey: e.StateTransitionVKey,
			StateMembershipVKey: e.StateMembershipVKey,
			TxId:                msg.TxId,
		}
		if e.Creator != "" {
			ism.Creator = &storage.Address{Address: e.Creator}
		}

		ctx.AddZkISM(ism)
		break
	}
	return nil
}

func handleUpdateZkISM(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/celestia.zkism.v1.MsgUpdateInterchainSecurityModule" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processUpdateZkISM(ctx, c, msg)
}

func processUpdateZkISM(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule {
			continue
		}

		e, err := decode.NewZkISMUpdateEvent(event.Data)
		if err != nil {
			return errors.Wrap(err, "parsing update zkism event")
		}

		var addr *storage.Address
		if signer := msg.Data.GetStringOrDefault("Signer"); signer != "" {
			addr = &storage.Address{
				Address:    signer,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			}
		}

		ctx.AddZkISM(&storage.ZkISM{
			ExternalId: e.Id,
			State:      e.NewState,
		})
		ctx.AddZkIsmUpdate(&storage.ZkISMUpdate{
			Height:          ctx.Block.Height,
			Time:            ctx.Block.Time,
			NewState:        e.NewState,
			TxId:            msg.TxId,
			ZkISMExternalId: e.Id,
			Signer:          addr,
		})
		break
	}
	return nil
}

func handleSubmitZkISMMessages(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/celestia.zkism.v1.MsgSubmitMessages" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processSubmitZkISMMessages(ctx, c, msg)
}

func processSubmitZkISMMessages(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeCelestiazkismv1EventSubmitMessages {
			continue
		}

		e, err := decode.NewZkISMSubmitMessagesEvent(event.Data)
		if err != nil {
			return errors.Wrap(err, "parsing submit zkism messages event")
		}

		var addr *storage.Address
		if signer := msg.Data.GetStringOrDefault("Signer"); signer != "" {
			addr = &storage.Address{
				Address:    signer,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			}
		}

		for i := range e.MessageIds {
			ctx.AddZkIsmMessage(&storage.ZkISMMessage{
				Height:          ctx.Block.Height,
				Time:            ctx.Block.Time,
				StateRoot:       e.StateRoot,
				MessageId:       e.MessageIds[i],
				TxId:            msg.TxId,
				Signer:          addr,
				ZkISMExternalId: e.Id,
			})
		}

		break
	}
	return nil
}

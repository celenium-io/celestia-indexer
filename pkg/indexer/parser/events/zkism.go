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

func handleCreateZkISM(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/celestia.zkism.v1.MsgCreateInterchainSecurityModule" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processCreateZkISM(ctx, events, msg, idx)
}

func processCreateZkISM(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeCelestiazkismv1EventCreateInterchainSecurityModule {
			e, err := decode.NewZkISMCreateEvent(events[*idx].Data)
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
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}
	return nil
}

func handleUpdateZkISM(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/celestia.zkism.v1.MsgUpdateInterchainSecurityModule" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	if err := processUpdateZkISM(ctx, events, msg, idx); err != nil {
		return err
	}
	return nil
}

func processUpdateZkISM(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule {
			e, err := decode.NewZkISMUpdateEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing update zkism event")
			}

			var addr *storage.Address
			if signer := msg.Data.GetStringOrDefault("Signer"); signer != "" {
				addr = &storage.Address{
					Address:    signer,
					Height:     msg.Height,
					LastHeight: msg.Height,
					Balance:    storage.EmptyBalance(),
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
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}
	return nil
}

func handleSubmitZkISMMessages(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/celestia.zkism.v1.MsgSubmitMessages" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	if err := processSubmitZkISMMessages(ctx, events, msg, idx); err != nil {
		return err
	}
	return nil
}

func processSubmitZkISMMessages(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeCelestiazkismv1EventSubmitMessages {
			e, err := decode.NewZkISMSubmitMessagesEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing submit zkism messages event")
			}

			var addr *storage.Address
			if signer := msg.Data.GetStringOrDefault("Signer"); signer != "" {
				addr = &storage.Address{
					Address:    signer,
					Height:     msg.Height,
					LastHeight: msg.Height,
					Balance:    storage.EmptyBalance(),
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

			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}
	return nil
}

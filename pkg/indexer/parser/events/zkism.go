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
				StateRoot:           e.StateRoot,
				MerkleTreeAddress:   e.MerkleTreeAddress,
				Groth16VKey:         e.Groth16VKey,
				StateTransitionVKey: e.StateTransitionVKey,
				StateMembershipVKey: e.StateMembershipVKey,
			}
			if e.Creator != "" {
				ism.Creator = &storage.Address{Address: e.Creator}
			}

			ctx.AddZkISM(ism)
			msg.ZkISM = ism
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
	return processUpdateZkISM(ctx, events, msg, idx)
}

func processUpdateZkISM(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeCelestiazkismv1EventUpdateInterchainSecurityModule {
			e, err := decode.NewZkISMUpdateEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing update zkism event")
			}

			update := &storage.ZkISMUpdate{
				Height:       ctx.Block.Height,
				Time:         ctx.Block.Time,
				NewState:     e.NewState,
				NewStateRoot: e.NewStateRoot,
			}
			if e.Signer != "" {
				update.Signer = &storage.Address{Address: e.Signer}
			}

			ctx.AddZkISM(&storage.ZkISM{
				ExternalId: e.Id,
				State:      e.NewState,
				StateRoot:  e.NewStateRoot,
			})
			msg.ZkISMUpdate = update
			msg.ZkISMUpdate.ZkISMExternalId = e.Id
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
	return processSubmitZkISMMessages(ctx, events, msg, idx)
}

func processSubmitZkISMMessages(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeCelestiazkismv1EventSubmitMessages {
			e, err := decode.NewZkISMSubmitMessagesEvent(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing submit zkism messages event")
			}

			msgs := make([]*storage.ZkISMMessage, len(e.MessageIds))
			for i, msgId := range e.MessageIds {
				m := &storage.ZkISMMessage{
					Height:    ctx.Block.Height,
					Time:      ctx.Block.Time,
					StateRoot: e.StateRoot,
					MessageId: msgId,
				}
				if e.Signer != "" {
					m.Signer = &storage.Address{Address: e.Signer}
				}
				m.ZkISMExternalId = e.Id
				msgs[i] = m
			}

			msg.ZkISMMessages = msgs
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}
	return nil
}

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleCreateMailbox(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.core.v1.MsgCreateMailbox" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processCreateMailbox(ctx, events, msg, idx)
}

func processCreateMailbox(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeHyperlanecorev1EventCreateMailbox {
			createMailbox, err := decode.NewCreateMailbox(events[*idx].Data)
			if err != nil {
				return errors.Wrap(err, "parsing create mailbox event")
			}
			mailboxId, err := util.DecodeHexAddress(createMailbox.MailboxId)
			if err != nil {
				return errors.Wrap(err, "decode mailbox id")
			}
			defaultIsm, err := util.DecodeHexAddress(createMailbox.DefaultIsm)
			if err != nil {
				return errors.Wrap(err, "decode default ISM")
			}

			mailbox := &storage.HLMailbox{
				Height:  ctx.Block.Height,
				Time:    ctx.Block.Time,
				Mailbox: mailboxId.Bytes(),
				Owner: &storage.Address{
					Address: createMailbox.Owner,
				},
				DefaultIsm: defaultIsm.Bytes(),
				Domain:     createMailbox.LocalDomain,
			}

			if len(createMailbox.DefaultHook) > 0 {
				defaultHook, err := util.DecodeHexAddress(createMailbox.DefaultHook)
				if err != nil {
					return errors.Wrap(err, "decode default hook")
				}
				mailbox.DefaultHook = defaultHook.Bytes()
			}

			if len(createMailbox.RequiredHook) > 0 {
				requiredHook, err := util.DecodeHexAddress(createMailbox.RequiredHook)
				if err != nil {
					return errors.Wrap(err, "decode required hook")
				}
				mailbox.DefaultHook = requiredHook.Bytes()
			}

			msg.HLMailbox = mailbox
			end = true
		}

		action := decoder.StringFromMap(events[*idx].Data, "action")
		end = len(events)-1 == *idx || action != "" || end
		*idx += 1
	}

	toTheNextAction(events, idx)
	return nil
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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

const null = "null"

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
				return errors.Wrapf(err, "decode mailbox id: %s", createMailbox.MailboxId)
			}
			defaultIsm, err := util.DecodeHexAddress(createMailbox.DefaultIsm)
			if err != nil {
				return errors.Wrapf(err, "decode default ISM: %s", createMailbox.DefaultIsm)
			}

			mailbox := &storage.HLMailbox{
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Mailbox:    mailboxId.Bytes(),
				InternalId: mailboxId.GetInternalId(),
				Owner: &storage.Address{
					Address: createMailbox.Owner,
				},
				DefaultIsm: defaultIsm.Bytes(),
				Domain:     createMailbox.LocalDomain,
			}

			if len(createMailbox.DefaultHook) > 0 && createMailbox.DefaultHook != null {
				defaultHook, err := util.DecodeHexAddress(createMailbox.DefaultHook)
				if err != nil {
					return errors.Wrapf(err, "decode default hook: %s", createMailbox.DefaultHook)
				}
				mailbox.DefaultHook = defaultHook.Bytes()
			}

			if len(createMailbox.RequiredHook) > 0 && createMailbox.RequiredHook != null {
				requiredHook, err := util.DecodeHexAddress(createMailbox.RequiredHook)
				if err != nil {
					return errors.Wrapf(err, "decode required hook: %s", createMailbox.RequiredHook)
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
	return nil
}

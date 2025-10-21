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

func handleSetMailbox(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/hyperlane.core.v1.MsgSetMailbox" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processSetMailbox(ctx, events, msg, idx)
}

func processSetMailbox(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	end := false
	for !end {
		if events[*idx].Type == types.EventTypeHyperlanecorev1EventSetMailbox {
			setMailbox, err := decode.NewSetMailbox(events[*idx].Data)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(setMailbox.MailboxId)
			if err != nil {
				return errors.Wrap(err, "decode mailbox id")
			}

			mailbox := &storage.HLMailbox{
				Height:     ctx.Block.Height,
				Time:       ctx.Block.Time,
				Mailbox:    mailboxId.Bytes(),
				InternalId: mailboxId.GetInternalId(),
				Owner: &storage.Address{
					Address: setMailbox.Owner,
				},
			}

			if len(setMailbox.DefaultIsm) > 0 && setMailbox.DefaultIsm != null {
				defaultIsm, err := util.DecodeHexAddress(setMailbox.DefaultIsm)
				if err != nil {
					return errors.Wrapf(err, "decode default ISM: %s", setMailbox.DefaultIsm)
				}
				mailbox.DefaultIsm = defaultIsm.Bytes()
			}

			if len(setMailbox.DefaultHook) > 0 && setMailbox.DefaultHook != null {
				defaultHook, err := util.DecodeHexAddress(setMailbox.DefaultHook)
				if err != nil {
					return errors.Wrapf(err, "decode default hook: %s", setMailbox.DefaultHook)
				}
				mailbox.DefaultHook = defaultHook.Bytes()
			}

			if len(setMailbox.NewOwner) > 0 {
				mailbox.DefaultHook = []byte(setMailbox.NewOwner)
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

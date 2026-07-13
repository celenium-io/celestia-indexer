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

func handleSetMailbox(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/hyperlane.core.v1.MsgSetMailbox" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()
	return processSetMailbox(ctx, c, msg)
}

func processSetMailbox(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for event := range c.MsgEvents("action") {
		if event.Type != types.EventTypeHyperlanecorev1EventSetMailbox {
			continue
		}

		setMailbox, err := decode.NewSetMailbox(event.Data)
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
				Address:    setMailbox.Owner,
				Height:     msg.Height,
				LastHeight: msg.Height,
				Balances:   []storage.Balance{storage.EmptyBalance()},
			},
		}
		if err := ctx.AddAddress(mailbox.Owner); err != nil {
			return errors.Wrap(err, "add address")
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

		if len(setMailbox.NewOwner) > 0 && setMailbox.NewOwner != setMailbox.Owner {
			newOwner := &storage.Address{
				Address:    setMailbox.NewOwner,
				Balances:   []storage.Balance{storage.EmptyBalance()},
				Height:     msg.Height,
				LastHeight: msg.Height,
			}
			if err := ctx.AddAddress(newOwner); err != nil {
				return errors.Wrap(err, "add address")
			}
			mailbox.Owner = newOwner
		}

		ctx.AddHlMailbox(mailbox)
		break
	}
	return nil
}

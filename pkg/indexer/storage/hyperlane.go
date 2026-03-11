// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveHlMailboxes(
	ctx context.Context,
	tx storage.Transaction,
	mailboxes []*storage.HLMailbox,
	addrToId map[string]uint64,
) error {
	if len(mailboxes) == 0 {
		return nil
	}

	for i := range mailboxes {
		if mailboxes[i].Owner != nil {
			if addrId, ok := addrToId[mailboxes[i].Owner.Address]; ok {
				mailboxes[i].OwnerId = addrId
			}
		}
	}

	return tx.SaveHyperlaneMailbox(ctx, mailboxes...)
}

func saveHlTokens(
	ctx context.Context,
	tx storage.Transaction,
	tokens []*storage.HLToken,
	addrToId map[string]uint64,
) error {
	if len(tokens) == 0 {
		return nil
	}

	for i := range tokens {
		if tokens[i].Owner != nil {
			if addrId, ok := addrToId[tokens[i].Owner.Address]; ok {
				tokens[i].OwnerId = addrId
			}
		}

		if tokens[i].Mailbox != nil {
			mailbox, err := tx.HyperlaneMailbox(ctx, tokens[i].Mailbox.InternalId)
			if err != nil {
				return errors.Wrapf(err, "can't find mailbox for token: %x", tokens[i].Mailbox)
			}
			tokens[i].MailboxId = mailbox.Id
		}
	}

	return tx.SaveHyperlaneTokens(ctx, tokens...)
}

func saveHlTransfers(
	ctx context.Context,
	tx storage.Transaction,
	transfers []*storage.HLTransfer,
	addrToId map[string]uint64,
) error {
	if len(transfers) == 0 {
		return nil
	}

	gasPayments := make([]*storage.HLGasPayment, 0, len(transfers))
	for i := range transfers {
		if transfers[i].Relayer != nil {
			if addrId, ok := addrToId[transfers[i].Relayer.Address]; ok {
				transfers[i].RelayerId = addrId
			}
		}
		if transfers[i].Address != nil {
			if addrId, ok := addrToId[transfers[i].Address.Address]; ok {
				transfers[i].AddressId = addrId
			}
		}

		if transfers[i].Mailbox != nil {
			mailbox, err := tx.HyperlaneMailbox(ctx, transfers[i].Mailbox.InternalId)
			if err != nil {
				return errors.Wrapf(err, "can't find mailbox for token: %x", transfers[i].Mailbox)
			}
			transfers[i].MailboxId = mailbox.Id
		}

		if transfers[i].Token != nil {
			token, err := tx.HyperlaneToken(ctx, transfers[i].Token.TokenId)
			if err != nil {
				return errors.Wrapf(err, "can't find token for transfer: %x", transfers[i].Token.TokenId)
			}
			transfers[i].TokenId = token.Id
		}

		if transfers[i].GasPayment != nil {
			igp, err := tx.HyperlaneIgp(ctx, transfers[i].GasPayment.Igp.IgpId)
			if err != nil {
				return errors.Wrapf(err, "can't find igp for transfer: %x", transfers[i].GasPayment.Igp.IgpId)
			}

			transfers[i].GasPayment.IgpId = igp.Id
		}
	}

	if err := tx.SaveHyperlaneTransfers(ctx, transfers...); err != nil {
		return err
	}

	for i := range transfers {
		if transfers[i].GasPayment != nil {
			transfers[i].GasPayment.TransferId = transfers[i].Id
			gasPayments = append(gasPayments, transfers[i].GasPayment)
		}
	}

	if len(gasPayments) > 0 {
		if err := tx.SaveHyperlaneGasPayments(ctx, gasPayments...); err != nil {
			return errors.Wrap(err, "hyperlane gas payments saving")
		}
	}

	return nil
}

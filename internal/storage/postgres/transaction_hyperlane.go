// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"

	models "github.com/celenium-io/celestia-indexer/internal/storage"
)

func (tx Transaction) SaveHyperlaneIgps(ctx context.Context, igps ...*models.HLIGP) error {
	if len(igps) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&igps).
		Column("id", "height", "time", "igp_id", "owner_id", "denom").
		On("CONFLICT (igp_id) DO UPDATE").
		Set("height = EXCLUDED.height").
		Set("time = EXCLUDED.time").
		Set("owner_id = EXCLUDED.owner_id").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveHyperlaneIgpConfigs(ctx context.Context, configs ...models.HLIGPConfig) error {
	if len(configs) == 0 {
		return nil
	}

	_, err := tx.Tx().NewInsert().Model(&configs).
		Column("id", "height", "time", "gas_overhead", "gas_price", "remote_domain", "token_exchange_rate").
		On("CONFLICT (id, remote_domain) DO UPDATE").
		Set("height = EXCLUDED.height").
		Set("time = EXCLUDED.time").
		Set("gas_overhead = EXCLUDED.gas_overhead").
		Set("gas_price = EXCLUDED.gas_price").
		Set("token_exchange_rate = EXCLUDED.token_exchange_rate").
		Exec(ctx)
	return err
}

func (tx Transaction) SaveHyperlaneMailbox(ctx context.Context, mailbox ...*models.HLMailbox) error {
	if len(mailbox) == 0 {
		return nil
	}

	for i := range mailbox {
		query := tx.Tx().NewInsert().
			Model(mailbox[i]).
			Column("height", "time", "tx_id", "mailbox", "internal_id", "owner_id", "default_ism", "default_hook", "required_hook", "domain", "sent_messages", "received_messages").
			On("CONFLICT (internal_id) DO UPDATE")

		if mailbox[i].Owner != nil || mailbox[i].OwnerRenounced {
			query.Set("owner_id = EXCLUDED.owner_id")
		}
		if mailbox[i].SentMessages > 0 {
			query.Set("sent_messages = hl_mailbox.sent_messages + EXCLUDED.sent_messages")
		}
		if mailbox[i].ReceivedMessages > 0 {
			query.Set("received_messages = hl_mailbox.received_messages + EXCLUDED.received_messages")
		}
		if mailbox[i].DefaultIsm != nil {
			query.Set("default_ism = EXCLUDED.default_ism")
		}
		if mailbox[i].DefaultHook != nil {
			query.Set("default_hook = EXCLUDED.default_hook")
		}
		if mailbox[i].RequiredHook != nil {
			query.Set("required_hook = EXCLUDED.required_hook")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveHyperlaneTokens(ctx context.Context, tokens ...*models.HLToken) error {
	if len(tokens) == 0 {
		return nil
	}

	for i := range tokens {
		query := tx.Tx().NewInsert().
			Model(tokens[i]).
			Column("height", "time", "tx_id", "mailbox_id", "owner_id", "type", "denom", "token_id", "sent_transfers", "received_transfers", "sent", "received").
			On("CONFLICT (token_id) DO UPDATE")

		if tokens[i].Owner != nil || tokens[i].OwnerRenounced {
			query.Set("owner_id = EXCLUDED.owner_id")
		}
		if tokens[i].SentTransfers > 0 {
			query.Set("sent_transfers = hl_token.sent_transfers + EXCLUDED.sent_transfers")
		}
		if tokens[i].ReceiveTransfers > 0 {
			query.Set("received_transfers = hl_token.received_transfers + EXCLUDED.received_transfers")
		}
		if !tokens[i].Sent.IsZero() {
			query.Set("sent = hl_token.sent + EXCLUDED.sent")
		}
		if !tokens[i].Received.IsZero() {
			query.Set("received = hl_token.received + EXCLUDED.received")
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tx Transaction) SaveHyperlaneTransfers(ctx context.Context, transfers ...*models.HLTransfer) error {
	if len(transfers) == 0 {
		return nil
	}
	_, err := tx.Tx().NewInsert().Model(&transfers).Returning("id").Exec(ctx)
	return err
}

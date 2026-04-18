// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveForwarding(
	ctx context.Context,
	tx storage.Transaction,
	forwardings []*storage.Forwarding,
	addrToId map[string]uint64,
) error {
	if len(forwardings) == 0 {
		return nil
	}

	for i := range forwardings {
		if forwardings[i].Address != nil {
			if addrId, ok := addrToId[forwardings[i].Address.Address]; ok {
				forwardings[i].AddressId = addrId
			}
		}
		if forwardings[i].Token != nil {
			token, err := tx.HyperlaneToken(ctx, forwardings[i].Token.TokenId)
			if err != nil {
				return errors.Wrapf(err, "can't find token for forwarding: %x", forwardings[i].Token.TokenId)
			}
			forwardings[i].TokenId = token.Id
		} else {
			return errors.Errorf("token is nil in forwarding")
		}
	}

	return tx.SaveForwardings(ctx, forwardings...)
}

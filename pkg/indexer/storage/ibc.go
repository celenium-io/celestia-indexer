// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

func saveIbcClients(
	ctx context.Context,
	tx storage.Transaction,
	clients []*storage.IbcClient,
	addrToId map[string]uint64,
) (int64, error) {
	if len(clients) == 0 {
		return 0, nil
	}

	for i := range clients {
		if clients[i].Creator != nil {
			if addrId, ok := addrToId[clients[i].Creator.Address]; ok {
				clients[i].CreatorId = addrId
			}
		}
	}

	return tx.SaveIbcClients(ctx, clients...)
}

func saveIbcChannels(
	ctx context.Context,
	tx storage.Transaction,
	channels []*storage.IbcChannel,
	addrToId map[string]uint64,
) error {
	if len(channels) == 0 {
		return nil
	}

	for i := range channels {
		if channels[i].ConnectionId != "" {
			conn, err := tx.IbcConnection(ctx, channels[i].ConnectionId)
			if err != nil {
				return errors.Wrap(err, "receiving connection for channel")
			}
			channels[i].ClientId = conn.ClientId
		}

		if channels[i].Creator != nil {
			if addrId, ok := addrToId[channels[i].Creator.Address]; ok {
				channels[i].CreatorId = addrId
			}
		}
	}

	return tx.SaveIbcChannels(ctx, channels...)
}

func saveIbcTransfers(
	ctx context.Context,
	tx storage.Transaction,
	transfers []*storage.IbcTransfer,
	addrToId map[string]uint64,
) error {
	if len(transfers) == 0 {
		return nil
	}

	for i := range transfers {
		if transfers[i].Sender != nil {
			if addrId, ok := addrToId[transfers[i].Sender.Address]; ok {
				transfers[i].SenderId = &addrId
			}
		}
		if transfers[i].Receiver != nil {
			if addrId, ok := addrToId[transfers[i].Receiver.Address]; ok {
				transfers[i].ReceiverId = &addrId
			}
		}
	}

	return tx.SaveIbcTransfers(ctx, transfers...)
}

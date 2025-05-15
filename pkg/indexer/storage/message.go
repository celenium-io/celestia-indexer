// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"maps"
	"slices"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

var errCantFindAddress = errors.New("can't find address")

func (module *Module) saveMessages(
	ctx context.Context,
	tx storage.Transaction,
	messages []*storage.Message,
	addrToId map[string]uint64,
) (int64, error) {
	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return 0, err
	}

	var (
		namespaceMsgs   []storage.NamespaceMessage
		msgAddress      []storage.MsgAddress
		blobLogs        = make([]storage.BlobLog, 0)
		vestingAccounts = make([]*storage.VestingAccount, 0)
		ibcClients      = make(map[string]*storage.IbcClient)
		ibcConnections  = make(map[string]*storage.IbcConnection)
		ibcChannels     = make([]*storage.IbcChannel, 0)
		grants          = make(map[string]storage.Grant)
		namespaces      = make(map[string]uint64)
		addedMsgId      = make(map[uint64]struct{})
		msgAddrMap      = make(map[string]struct{})
	)
	for i := range messages {
		for j := range messages[i].Namespace {
			nsId := messages[i].Namespace[j].Id
			key := messages[i].Namespace[j].String()
			if nsId == 0 {
				if _, ok := addedMsgId[messages[i].Id]; ok { // in case of duplication of writing to one namespace inside one messages
					continue
				}

				id, ok := namespaces[key]
				if !ok {
					continue
				}
				nsId = id
			} else {
				namespaces[key] = nsId
			}

			addedMsgId[messages[i].Id] = struct{}{}
			namespaceMsgs = append(namespaceMsgs, storage.NamespaceMessage{
				MsgId:       messages[i].Id,
				NamespaceId: nsId,
				Time:        messages[i].Time,
				Height:      messages[i].Height,
				TxId:        messages[i].TxId,
				Size:        uint64(messages[i].Namespace[j].Size),
			})
		}

		for j := range messages[i].Addresses {
			id, ok := addrToId[messages[i].Addresses[j].String()]
			if !ok {
				continue
			}
			msgAddressEntity := storage.MsgAddress{
				MsgId:     messages[i].Id,
				AddressId: id,
				Type:      messages[i].Addresses[j].Type,
			}
			key := msgAddressEntity.String()
			if _, ok := msgAddrMap[key]; !ok {
				msgAddress = append(msgAddress, msgAddressEntity)
				msgAddrMap[key] = struct{}{}
			}
		}

		for j := range messages[i].BlobLogs {
			if err := processPayForBlob(addrToId, namespaces, messages[i], messages[i].BlobLogs[j]); err != nil {
				return 0, err
			}

			blobLogs = append(blobLogs, *messages[i].BlobLogs[j])
		}

		if messages[i].VestingAccount != nil {
			addrId, ok := addrToId[messages[i].VestingAccount.Address.Address]
			if !ok {
				continue
			}
			messages[i].VestingAccount.AddressId = addrId
			messages[i].VestingAccount.TxId = &messages[i].TxId
			vestingAccounts = append(vestingAccounts, messages[i].VestingAccount)
		}

		for j := range messages[i].Grants {
			if err := processGrants(addrToId, &messages[i].Grants[j]); err != nil {
				return 0, err
			}
			grants[messages[i].Grants[j].String()] = messages[i].Grants[j]
		}

		if messages[i].IbcClient != nil {
			messages[i].IbcClient.TxId = messages[i].TxId
			if messages[i].IbcClient.Creator != nil {
				if addrId, ok := addrToId[messages[i].IbcClient.Creator.Address]; ok {
					messages[i].IbcClient.CreatorId = addrId
				}
			}
			if client, ok := ibcClients[messages[i].IbcClient.Id]; ok {
				client.ConnectionCount += messages[i].IbcClient.ConnectionCount
			} else {
				ibcClients[messages[i].IbcClient.Id] = messages[i].IbcClient
			}
		}

		if messages[i].IbcConnection != nil {
			if messages[i].IbcConnection.Height > 0 {
				messages[i].IbcConnection.CreateTxId = messages[i].TxId
			}
			if messages[i].IbcConnection.ConnectionHeight > 0 {
				messages[i].IbcConnection.ConnectionTxId = messages[i].TxId
			}

			if conn, ok := ibcConnections[messages[i].IbcConnection.ConnectionId]; ok {
				conn.ChannelsCount += messages[i].IbcConnection.ChannelsCount
			} else {
				ibcConnections[messages[i].IbcConnection.ConnectionId] = messages[i].IbcConnection
			}
		}

		if messages[i].IbcChannel != nil {
			if messages[i].IbcChannel.Height > 0 {
				messages[i].IbcChannel.CreateTxId = messages[i].TxId
			}
			if messages[i].IbcChannel.ConfirmationHeight > 0 {
				messages[i].IbcChannel.ConfirmationTxId = messages[i].TxId
			}

			if messages[i].IbcChannel.ConnectionId != "" {
				conn, err := tx.IbcConnection(ctx, messages[i].IbcChannel.ConnectionId)
				if err != nil {
					return 0, errors.Wrap(err, "receiving connection for channel")
				}
				messages[i].IbcChannel.ClientId = conn.ClientId
			}

			if messages[i].IbcChannel.Creator != nil {
				if addrId, ok := addrToId[messages[i].IbcChannel.Creator.Address]; ok {
					messages[i].IbcChannel.CreatorId = addrId
				}
			}

			ibcChannels = append(ibcChannels, messages[i].IbcChannel)
		}
	}

	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return 0, err
	}
	if err := tx.SaveMsgAddresses(ctx, msgAddress...); err != nil {
		return 0, err
	}
	if err := tx.SaveBlobLogs(ctx, blobLogs...); err != nil {
		return 0, err
	}

	grantsArr := make([]storage.Grant, 0)
	for _, g := range grants {
		grantsArr = append(grantsArr, g)
	}

	if err := tx.SaveGrants(ctx, grantsArr...); err != nil {
		return 0, err
	}

	if len(vestingAccounts) > 0 {
		if err := tx.SaveVestingAccounts(ctx, vestingAccounts...); err != nil {
			return 0, err
		}

		vestingPeriods := make([]storage.VestingPeriod, 0)
		for i := range vestingAccounts {
			for j := range vestingAccounts[i].VestingPeriods {
				vestingAccounts[i].VestingPeriods[j].VestingAccountId = vestingAccounts[i].Id
			}
			vestingPeriods = append(vestingPeriods, vestingAccounts[i].VestingPeriods...)
		}

		if err := tx.SaveVestingPeriods(ctx, vestingPeriods...); err != nil {
			return 0, err
		}
	}

	ibcClientsCount, err := tx.SaveIbcClients(ctx, slices.Collect(maps.Values(ibcClients))...)
	if err != nil {
		return 0, errors.Wrap(err, "ibc clients saving")
	}

	if err := tx.SaveIbcConnections(ctx, slices.Collect(maps.Values(ibcConnections))...); err != nil {
		return 0, errors.Wrap(err, "ibc connections saving")
	}
	if err := tx.SaveIbcChannels(ctx, ibcChannels...); err != nil {
		return 0, errors.Wrap(err, "ibc channels saving")
	}

	return ibcClientsCount, nil
}

func processPayForBlob(addrToId map[string]uint64, namespaces map[string]uint64, msg *storage.Message, blob *storage.BlobLog) error {
	if blob.Namespace == nil {
		return errors.New("nil namespace in pay for blob message")
	}
	nsId, ok := namespaces[blob.Namespace.String()]
	if !ok {
		return errors.Errorf("can't find namespace for pay for blob message: %s", blob.Namespace.String())
	}
	if blob.Signer == nil {
		return errors.New("nil signer address in pay for blob message")
	}
	signerId, ok := addrToId[blob.Signer.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "signer for pay for blob message: %s", blob.Signer.Address)
	}
	blob.MsgId = msg.Id
	blob.TxId = msg.TxId
	blob.SignerId = signerId
	blob.NamespaceId = nsId
	return nil
}

func processGrants(addrToId map[string]uint64, grant *storage.Grant) error {
	granteeId, ok := addrToId[grant.Grantee.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "grantee: %s", grant.Grantee.Address)
	}
	grant.GranteeId = granteeId
	granterId, ok := addrToId[grant.Granter.Address]
	if !ok {
		return errors.Wrapf(errCantFindAddress, "granter: %s", grant.Granter.Address)
	}
	grant.GranterId = granterId
	return nil
}

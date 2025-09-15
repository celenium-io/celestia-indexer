// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"maps"
	"slices"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var errCantFindAddress = errors.New("can't find address")
var signalsThreshold = decimal.NewFromFloat(0.833333333) // 5/6

func (module *Module) saveMessages(
	ctx context.Context,
	tx storage.Transaction,
	messages []*storage.Message,
	addrToId map[string]uint64,
	state storage.State,
) (int64, uint64, error) {
	if err := tx.SaveMessages(ctx, messages...); err != nil {
		return 0, state.Version, err
	}

	var (
		namespaceMsgs      []storage.NamespaceMessage
		msgAddress         []storage.MsgAddress
		valMsgs            []storage.MsgValidator
		blobLogs           = make([]storage.BlobLog, 0)
		vestingAccounts    = make([]*storage.VestingAccount, 0)
		ibcClients         = make(map[string]*storage.IbcClient)
		ibcConnections     = make(map[string]*storage.IbcConnection)
		ibcChannels        = make(map[string]*storage.IbcChannel)
		ibcTransfers       = make([]*storage.IbcTransfer, 0)
		hyperlaneMailbox   = make(map[uint64]*storage.HLMailbox, 0)
		hyperlaneTokens    = make(map[string]*storage.HLToken, 0)
		hyperlaneTransfers = make([]*storage.HLTransfer, 0)
		grants             = make(map[string]storage.Grant)
		signals            = make([]*storage.SignalVersion, 0)
		upgrades           = make([]*storage.Upgrade, 0)
		namespaces         = make(map[string]uint64)
		addedMsgId         = make(map[uint64]struct{})
		msgAddrMap         = make(map[string]struct{})
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
				return 0, state.Version, err
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
				return 0, state.Version, err
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
					return 0, state.Version, errors.Wrap(err, "receiving connection for channel")
				}
				messages[i].IbcChannel.ClientId = conn.ClientId
			}

			if messages[i].IbcChannel.Creator != nil {
				if addrId, ok := addrToId[messages[i].IbcChannel.Creator.Address]; ok {
					messages[i].IbcChannel.CreatorId = addrId
				}
			}

			if conn, ok := ibcChannels[messages[i].IbcChannel.Id]; ok {
				conn.Received = conn.Received.Add(messages[i].IbcChannel.Received)
				conn.Sent = conn.Sent.Add(messages[i].IbcChannel.Sent)
				conn.TransfersCount += messages[i].IbcChannel.TransfersCount
			} else {
				ibcChannels[messages[i].IbcChannel.Id] = messages[i].IbcChannel
			}
		}

		if messages[i].IbcTransfer != nil {
			messages[i].IbcTransfer.TxId = messages[i].TxId

			if messages[i].IbcTransfer.Sender != nil {
				if addrId, ok := addrToId[messages[i].IbcTransfer.Sender.Address]; ok {
					messages[i].IbcTransfer.SenderId = &addrId
				}
			}
			if messages[i].IbcTransfer.Receiver != nil {
				if addrId, ok := addrToId[messages[i].IbcTransfer.Receiver.Address]; ok {
					messages[i].IbcTransfer.ReceiverId = &addrId
				}
			}

			ibcTransfers = append(ibcTransfers, messages[i].IbcTransfer)
		}

		if messages[i].HLMailbox != nil {
			messages[i].HLMailbox.TxId = messages[i].TxId

			if messages[i].HLMailbox.Owner != nil {
				if addrId, ok := addrToId[messages[i].HLMailbox.Owner.Address]; ok {
					messages[i].HLMailbox.OwnerId = addrId
				}
			}

			hyperlaneMailbox[messages[i].HLMailbox.InternalId] = messages[i].HLMailbox
		}

		if messages[i].HLToken != nil {
			messages[i].HLToken.TxId = messages[i].TxId

			if messages[i].HLToken.Owner != nil {
				if addrId, ok := addrToId[messages[i].HLToken.Owner.Address]; ok {
					messages[i].HLToken.OwnerId = addrId
				}
			}

			if messages[i].HLToken.Mailbox != nil {
				mailbox, err := tx.HyperlaneMailbox(ctx, messages[i].HLToken.Mailbox.InternalId)
				if err != nil {
					return 0, state.Version, errors.Wrapf(err, "can't find mailbox for token: %x", messages[i].HLToken.Mailbox)
				}
				messages[i].HLToken.MailboxId = mailbox.Id
			}

			hyperlaneTokens[messages[i].HLToken.String()] = messages[i].HLToken
		}

		if messages[i].HLTransfer != nil {
			messages[i].HLTransfer.TxId = messages[i].TxId

			if messages[i].HLTransfer.Relayer != nil {
				if addrId, ok := addrToId[messages[i].HLTransfer.Relayer.Address]; ok {
					messages[i].HLTransfer.RelayerId = addrId
				}
			}
			if messages[i].HLTransfer.Address != nil {
				if addrId, ok := addrToId[messages[i].HLTransfer.Address.Address]; ok {
					messages[i].HLTransfer.AddressId = addrId
				}
			}
			if messages[i].HLTransfer.Mailbox != nil {
				mailbox, err := tx.HyperlaneMailbox(ctx, messages[i].HLTransfer.Mailbox.InternalId)
				if err != nil {
					return 0, state.Version, errors.Wrapf(err, "can't find mailbox for transfer: %x", messages[i].HLTransfer.Mailbox)
				}
				messages[i].HLTransfer.MailboxId = mailbox.Id
				messages[i].HLTransfer.Mailbox.Id = mailbox.Id

				if hlm, ok := hyperlaneMailbox[messages[i].HLTransfer.Mailbox.InternalId]; ok {
					if messages[i].HLTransfer.Token != nil {
						hlm.ReceivedMessages += messages[i].HLTransfer.Token.ReceiveTransfers
						hlm.SentMessages += messages[i].HLTransfer.Token.SentTransfers
					}
				} else {
					hyperlaneMailbox[messages[i].HLTransfer.Mailbox.InternalId] = messages[i].HLTransfer.Mailbox
				}
			}
			if messages[i].HLTransfer.Token != nil {
				token, err := tx.HyperlaneToken(ctx, messages[i].HLTransfer.Token.TokenId)
				if err != nil {
					return 0, state.Version, errors.Wrapf(err, "can't find token for transfer: %x", messages[i].HLTransfer.Token.TokenId)
				}
				messages[i].HLTransfer.TokenId = token.Id
				messages[i].HLTransfer.Token.Id = token.Id

				if hlt, ok := hyperlaneTokens[messages[i].HLTransfer.Token.String()]; ok {
					hlt.ReceiveTransfers += messages[i].HLTransfer.Token.ReceiveTransfers
					hlt.SentTransfers += messages[i].HLTransfer.Token.SentTransfers
					hlt.Received = hlt.Received.Add(messages[i].HLTransfer.Token.Received)
					hlt.Sent = hlt.Sent.Add(messages[i].HLTransfer.Token.Sent)
				} else {
					hyperlaneTokens[messages[i].HLTransfer.Token.String()] = messages[i].HLTransfer.Token
				}
			}

			hyperlaneTransfers = append(hyperlaneTransfers, messages[i].HLTransfer)
		}

		if messages[i].SignalVersion != nil {
			messages[i].SignalVersion.TxId = messages[i].TxId
			messages[i].SignalVersion.MsgId = messages[i].Id
			validatorId, ok := module.validatorsByAddress[messages[i].SignalVersion.Validator.Address]
			if !ok {
				return 0, state.Version, errors.New("address:" + messages[i].SignalVersion.Validator.Address + " not found in validator map")
			}

			validator, err := tx.Validator(ctx, validatorId)
			if err != nil {
				return 0, state.Version, errors.Wrapf(err, "can't find validator for address: %s", messages[i].SignalVersion.Validator.Address)
			}

			messages[i].SignalVersion.VotingPower = validator.Stake
			messages[i].SignalVersion.ValidatorId = validatorId
			messages[i].SignalVersion.Validator = &validator

			signals = append(signals, messages[i].SignalVersion)
		}

		if messages[i].Upgrade != nil {
			messages[i].Upgrade.TxId = messages[i].TxId
			messages[i].Upgrade.MsgId = messages[i].Id
			signerId, ok := addrToId[messages[i].Upgrade.Signer.Address]
			if !ok {
				return 0, state.Version, errors.Errorf("address %s not found in addrToId map", messages[i].Upgrade.Signer.Address)
			}

			messages[i].Upgrade.SignerId = signerId

			vp, validators, err := module.totalVotingPower(ctx, tx)
			if err != nil {
				return 0, state.Version, errors.Wrapf(err, "receiving total voting power")
			}

			voted := decimal.Zero
			for _, v := range validators {
				if v.Version <= state.Version {
					continue
				}
				voted = voted.Add(v.Stake)
				if voted.GreaterThan(vp.Mul(signalsThreshold)) {
					messages[i].Upgrade.Version = v.Version
					state.Version = v.Version

					if err := tx.UpdateSignalsAfterUpgrade(ctx, v.Version); err != nil {
						return 0, state.Version, errors.Wrap(err, "updating signals after upgrade")
					}

					break
				}
			}

			upgrades = append(upgrades, messages[i].Upgrade)
		}

		if len(messages[i].Validators) > 0 {
			for _, val := range messages[i].Validators {
				id, ok := module.validatorsByAddress[val]
				if !ok {
					return 0, errors.Errorf("validator %s not found", val)
				}

				valMsgs = append(valMsgs, storage.MsgValidator{
					Height:      messages[i].Height,
					Time:        messages[i].Time,
					MsgId:       messages[i].Id,
					ValidatorId: id,
				})
			}
		}
	}

	if err := tx.SaveNamespaceMessage(ctx, namespaceMsgs...); err != nil {
		return 0, state.Version, err
	}
	if err := tx.SaveMsgAddresses(ctx, msgAddress...); err != nil {
		return 0, state.Version, err
	}
	if err := tx.SaveMsgValidator(ctx, valMsgs...); err != nil {
		return 0, err
	}
	if err := tx.SaveBlobLogs(ctx, blobLogs...); err != nil {
		return 0, state.Version, err
	}

	grantsArr := make([]storage.Grant, 0)
	for _, g := range grants {
		grantsArr = append(grantsArr, g)
	}

	if err := tx.SaveGrants(ctx, grantsArr...); err != nil {
		return 0, state.Version, err
	}

	if len(vestingAccounts) > 0 {
		if err := tx.SaveVestingAccounts(ctx, vestingAccounts...); err != nil {
			return 0, state.Version, err
		}

		vestingPeriods := make([]storage.VestingPeriod, 0)
		for i := range vestingAccounts {
			for j := range vestingAccounts[i].VestingPeriods {
				vestingAccounts[i].VestingPeriods[j].VestingAccountId = vestingAccounts[i].Id
			}
			vestingPeriods = append(vestingPeriods, vestingAccounts[i].VestingPeriods...)
		}

		if err := tx.SaveVestingPeriods(ctx, vestingPeriods...); err != nil {
			return 0, state.Version, err
		}
	}

	ibcClientsCount, err := tx.SaveIbcClients(ctx, slices.Collect(maps.Values(ibcClients))...)
	if err != nil {
		return 0, state.Version, errors.Wrap(err, "ibc clients saving")
	}

	if err := tx.SaveIbcConnections(ctx, slices.Collect(maps.Values(ibcConnections))...); err != nil {
		return 0, state.Version, errors.Wrap(err, "ibc connections saving")
	}
	if err := tx.SaveIbcChannels(ctx, slices.Collect(maps.Values(ibcChannels))...); err != nil {
		return 0, state.Version, errors.Wrap(err, "ibc channels saving")
	}
	if err := tx.SaveIbcTransfers(ctx, ibcTransfers...); err != nil {
		return 0, state.Version, errors.Wrap(err, "ibc transfers saving")
	}
	if err := tx.SaveHyperlaneMailbox(ctx, slices.Collect(maps.Values(hyperlaneMailbox))...); err != nil {
		return 0, state.Version, errors.Wrap(err, "hyperlane mailbox saving")
	}
	if err := tx.SaveHyperlaneTokens(ctx, slices.Collect(maps.Values(hyperlaneTokens))...); err != nil {
		return 0, state.Version, errors.Wrap(err, "hyperlane tokens saving")
	}
	if err := tx.SaveHyperlaneTransfers(ctx, hyperlaneTransfers...); err != nil {
		return 0, state.Version, errors.Wrap(err, "hyperlane transfers saving")
	}
	if err := tx.SaveSignals(ctx, signals...); err != nil {
		return 0, state.Version, errors.Wrap(err, "signals saving")
	}
	if err := tx.SaveUpgrades(ctx, upgrades...); err != nil {
		return 0, state.Version, errors.Wrap(err, "upgrades saving")
	}

	return ibcClientsCount, state.Version, nil
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

func (module *Module) totalVotingPower(ctx context.Context, tx storage.Transaction) (decimal.Decimal, []storage.Validator, error) {
	maxValsConsts, err := module.constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return decimal.Zero, nil, errors.Wrap(err, "get max validators value")
	}
	maxVals, err := strconv.Atoi(maxValsConsts.Value)
	if err != nil {
		return decimal.Zero, nil, errors.Wrap(err, "parse max validators value")
	}

	validators, err := tx.BondedValidators(ctx, maxVals)
	if err != nil {
		return decimal.Zero, nil, errors.Wrap(err, "get validators")
	}

	power := decimal.Zero
	for i := range validators {
		power.Add(validators[i].Stake)
	}
	return power, validators, nil
}

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/dipdup-io/celestia-indexer/pkg/indexer/decode/handle"
	"time"

	"github.com/rs/zerolog/log"

	appBlobTypes "github.com/celestiaorg/celestia-app/x/blob/types"
	qgbTypes "github.com/celestiaorg/celestia-app/x/qgb/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosFeegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	storageTypes "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/dipdup-io/celestia-indexer/pkg/types"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type DecodedMsg struct {
	Msg       storage.Message
	BlobsSize int64
	Addresses []storage.AddressWithType
}

func Message(
	msg cosmosTypes.Msg,
	height types.Level,
	time time.Time,
	position int,
	status storageTypes.Status,
) (d DecodedMsg, err error) {
	d.Msg.Height = height
	d.Msg.Time = time
	d.Msg.Position = int64(position)
	d.Msg.Data = structs.Map(msg)

	switch typedMsg := msg.(type) {

	// distribution module
	case *cosmosDistributionTypes.MsgSetWithdrawAddress:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetWithdrawAddress(height, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawDelegatorReward(height, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawValidatorCommission(height, typedMsg)
	case *cosmosDistributionTypes.MsgFundCommunityPool:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgFundCommunityPool(height, typedMsg)

	// staking module
	case *cosmosStakingTypes.MsgCreateValidator:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Validator, err = handle.MsgCreateValidator(height, status, typedMsg)
	case *cosmosStakingTypes.MsgEditValidator:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Validator, err = handle.MsgEditValidator(height, status, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgBeginRedelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUndelegate(height, typedMsg)
	case *cosmosStakingTypes.MsgCancelUnbondingDelegation:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCancelUnbondingDelegation(height, typedMsg)

	// slashing module
	case *cosmosSlashingTypes.MsgUnjail:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUnjail(height, typedMsg)

	// bank module
	case *cosmosBankTypes.MsgSend:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSend(height, typedMsg)
	case *cosmosBankTypes.MsgMultiSend:
		log.Warn().Msg("MsgMultiSend detected")
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgMultiSend(height, typedMsg)

	// vesting module
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateVestingAccount(height, typedMsg)
	case *cosmosVestingTypes.MsgCreatePermanentLockedAccount:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreatePermanentLockedAccount(height, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreatePeriodicVestingAccount(height, typedMsg)

	// blob module
	case *appBlobTypes.MsgPayForBlobs:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Namespace, d.BlobsSize, err = handle.MsgPayForBlobs(height, typedMsg)

	// feegrant module
	case *cosmosFeegrant.MsgGrantAllowance:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgGrantAllowance(height, typedMsg)
	case *cosmosFeegrant.MsgRevokeAllowance:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRevokeAllowance(height, typedMsg)

	// qgb module
	case *qgbTypes.MsgRegisterEVMAddress:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterEVMAddress(height, typedMsg)

	// authz module
	case *authz.MsgGrant:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgGrant(height, typedMsg)
	case *authz.MsgExec:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgExec(height, typedMsg)
	case *authz.MsgRevoke:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRevoke(height, typedMsg)

	// gov module
	case *cosmosGovTypesV1.MsgSubmitProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitProposal(height, typedMsg.Proposer)
	case *cosmosGovTypesV1Beta1.MsgSubmitProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitProposal(height, typedMsg.Proposer)
	case *cosmosGovTypesV1.MsgExecLegacyContent:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgExecLegacyContent(height, typedMsg)
	case *cosmosGovTypesV1.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVote(height, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVote(height, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgVoteWeighted:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteWeighted(height, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVoteWeighted:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteWeighted(height, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgDeposit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDeposit(height, typedMsg.Depositor)
	case *cosmosGovTypesV1Beta1.MsgDeposit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDeposit(height, typedMsg.Depositor)

	default:
		log.Err(errors.New("unknown message type")).Msgf("got type %T", msg)
		d.Msg.Type = storageTypes.MsgUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	d.Addresses = append(d.Addresses, d.Msg.Addresses...)
	return
}

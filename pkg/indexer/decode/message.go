// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/handle"
	"github.com/cosmos/cosmos-sdk/x/authz"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/nft"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	interchainAccounts "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	fee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	coreClient "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	coreConnection "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	coreChannel "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"

	"github.com/rs/zerolog/log"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
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
		d.Msg.Type, d.Msg.Addresses, d.Msg.Namespace, d.Msg.BlobLogs, d.BlobsSize, err = handle.MsgPayForBlobs(height, time, status, typedMsg)

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

	// ibc module
	case *ibcTypes.MsgTransfer:
		d.Msg.Type, d.Msg.Addresses, err = handle.IBCTransfer(height, typedMsg)

	// crisis module
	case *crisisTypes.MsgVerifyInvariant:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVerifyInvariant(height, typedMsg)

	// evidence module
	case *evidenceTypes.MsgSubmitEvidence:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitEvidence(height, typedMsg)

	// nft module
	case *nft.MsgSend:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSendNFT(height, typedMsg)

	// group module
	case *group.MsgCreateGroup:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroup(height, typedMsg)
	case *group.MsgUpdateGroupMembers:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupMembers(height, typedMsg)
	case *group.MsgUpdateGroupAdmin:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupAdmin(height, typedMsg)
	case *group.MsgUpdateGroupMetadata:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupMetadata(height, typedMsg)
	case *group.MsgCreateGroupPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroupPolicy(height, typedMsg)
	case *group.MsgUpdateGroupPolicyAdmin:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyAdmin(height, typedMsg)
	case *group.MsgCreateGroupWithPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroupWithPolicy(height, typedMsg)
	case *group.MsgUpdateGroupPolicyDecisionPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyDecisionPolicy(height, typedMsg)
	case *group.MsgUpdateGroupPolicyMetadata:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyMetadata(height, typedMsg)
	case *group.MsgSubmitProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitProposalGroup(height, typedMsg)
	case *group.MsgWithdrawProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawProposal(height, typedMsg)
	case *group.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteGroup(height, typedMsg)
	case *group.MsgExec:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgExecGroup(height, typedMsg)
	case *group.MsgLeaveGroup:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgLeaveGroup(height, typedMsg)

	// upgrade module
	case *upgrade.MsgSoftwareUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSoftwareUpgrade(height, typedMsg)
	case *upgrade.MsgCancelUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCancelUpgrade(height, typedMsg)

	// interchainAccounts module
	case *interchainAccounts.MsgRegisterInterchainAccount:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterInterchainAccount(height, typedMsg)
	case *interchainAccounts.MsgSendTx:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSendTx(height, typedMsg)

	// fee module
	case *fee.MsgRegisterPayee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterPayee(height, typedMsg)
	case *fee.MsgRegisterCounterpartyPayee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterCounterpartyPayee(height, typedMsg)
	case *fee.MsgPayPacketFee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgPayPacketFee(height, typedMsg)
	case *fee.MsgPayPacketFeeAsync:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgPayPacketFeeAsync()

	// coreClient module
	case *coreClient.MsgCreateClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateClient(height, typedMsg)
	case *coreClient.MsgUpdateClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateClient(height, typedMsg)
	case *coreClient.MsgUpgradeClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpgradeClient(height, typedMsg)
	case *coreClient.MsgSubmitMisbehaviour:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitMisbehaviour(height, typedMsg)

	// coreConnection module
	case *coreConnection.MsgConnectionOpenInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenInit(height, typedMsg)
	case *coreConnection.MsgConnectionOpenTry:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenTry(height, typedMsg)
	case *coreConnection.MsgConnectionOpenAck:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenAck(height, typedMsg)
	case *coreConnection.MsgConnectionOpenConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenConfirm(height, typedMsg)

	// coreChannel module
	case *coreChannel.MsgChannelOpenInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenInit(height, typedMsg)
	case *coreChannel.MsgChannelOpenTry:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenTry(height, typedMsg)
	case *coreChannel.MsgChannelOpenAck:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenAck(height, typedMsg)
	case *coreChannel.MsgChannelOpenConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenConfirm(height, typedMsg)
	case *coreChannel.MsgChannelCloseInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelCloseInit(height, typedMsg)
	case *coreChannel.MsgChannelCloseConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelCloseConfirm(height, typedMsg)
	case *coreChannel.MsgRecvPacket:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRecvPacket(height, typedMsg)
	case *coreChannel.MsgTimeout:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgTimeout(height, typedMsg)
	case *coreChannel.MsgTimeoutOnClose:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgTimeoutOnClose(height, typedMsg)
	case *coreChannel.MsgAcknowledgement:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgAcknowledgement(height, typedMsg)

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

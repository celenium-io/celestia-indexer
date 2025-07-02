// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	evidenceTypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/nft"
	upgrade "cosmossdk.io/x/upgrade/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/handle"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	interchainAccounts "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	fee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreClient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	coreConnection "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/rs/zerolog/log"

	cosmosFeegrant "cosmossdk.io/x/feegrant"
	hyperlanePostDispatch "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	hyperlaneCore "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
	appBlobTypes "github.com/celestiaorg/celestia-app/v4/x/blob/types"
	minfeeTypes "github.com/celestiaorg/celestia-app/v4/x/minfee/types"
	appSignalTypes "github.com/celestiaorg/celestia-app/v4/x/signal/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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

type measurable interface {
	Size() int
}

func Message(
	ctx *context.Context,
	msg cosmosTypes.Msg,
	position int,
	status storageTypes.Status,
) (d DecodedMsg, err error) {
	d.Msg.Position = int64(position)
	d.Msg.Data = structs.Map(msg)
	d.Msg.Height = ctx.Block.Height
	d.Msg.Time = ctx.Block.Time

	switch typedMsg := msg.(type) {

	// distribution module
	case *cosmosDistributionTypes.MsgSetWithdrawAddress:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetWithdrawAddress(ctx, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawDelegatorReward(ctx, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawValidatorCommission(ctx, typedMsg)
	case *cosmosDistributionTypes.MsgFundCommunityPool:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgFundCommunityPool(ctx, typedMsg)

	// staking module
	case *cosmosStakingTypes.MsgCreateValidator:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateValidator(ctx, status, typedMsg)
		if err != nil {
			return d, err
		}
		if pk, ok := typedMsg.Pubkey.GetCachedValue().(cryptotypes.PubKey); ok {
			d.Msg.Data["Pubkey"] = map[string]any{
				"key":  pk.Bytes(),
				"type": pk.Type(),
			}
		}
	case *cosmosStakingTypes.MsgEditValidator:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgEditValidator(ctx, status, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDelegate(ctx, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgBeginRedelegate(ctx, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUndelegate(ctx, typedMsg)
	case *cosmosStakingTypes.MsgCancelUnbondingDelegation:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCancelUnbondingDelegation(ctx, typedMsg)

	// slashing module
	case *cosmosSlashingTypes.MsgUnjail:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUnjail(ctx, typedMsg)

	// bank module
	case *cosmosBankTypes.MsgSend:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSend(ctx, typedMsg)
	case *cosmosBankTypes.MsgMultiSend:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgMultiSend(ctx, typedMsg)

	// vesting module
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.Msg.Type, d.Msg.Addresses, d.Msg.VestingAccount, err = handle.MsgCreateVestingAccount(ctx, status, typedMsg)
	case *cosmosVestingTypes.MsgCreatePermanentLockedAccount:
		d.Msg.Type, d.Msg.Addresses, d.Msg.VestingAccount, err = handle.MsgCreatePermanentLockedAccount(ctx, status, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.Msg.Type, d.Msg.Addresses, d.Msg.VestingAccount, err = handle.MsgCreatePeriodicVestingAccount(ctx, status, typedMsg)

	// blob module
	case *appBlobTypes.MsgPayForBlobs:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Namespace, d.Msg.BlobLogs, d.BlobsSize, err = handle.MsgPayForBlobs(ctx, status, typedMsg)

	// feegrant module
	case *cosmosFeegrant.MsgGrantAllowance:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Grants, err = handle.MsgGrantAllowance(ctx, status, typedMsg)
	case *cosmosFeegrant.MsgRevokeAllowance:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Grants, err = handle.MsgRevokeAllowance(ctx, status, typedMsg)

	// qgb module
	case *legacy.MsgRegisterEVMAddress:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterEVMAddress(ctx, typedMsg)

	// authz module
	case *authz.MsgGrant:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Grants, err = handle.MsgGrant(ctx, status, typedMsg)
	case *authz.MsgExec:
		d.Msg.Type, d.Msg.Addresses, d.Msg.InternalMsgs, err = handle.MsgExec(ctx, status, typedMsg)
		if err != nil {
			return d, err
		}

		msgs := make([]any, 0)
		for i := range typedMsg.Msgs {
			msg, err := cosmosTypes.GetMsgFromTypeURL(cfg.Codec, typedMsg.Msgs[i].TypeUrl)
			if err != nil {
				return d, err
			}
			if err := cfg.Codec.UnpackAny(typedMsg.Msgs[i], &msg); err != nil {
				return d, err
			}
			msgs = append(msgs, structs.Map(msg))
		}
		d.Msg.Data["Msgs"] = msgs

	case *authz.MsgRevoke:
		d.Msg.Type, d.Msg.Addresses, d.Msg.Grants, err = handle.MsgRevoke(ctx, status, typedMsg)

	// gov module
	case *cosmosGovTypesV1.MsgSubmitProposal:
		var msgs []any
		d.Msg.Type, d.Msg.Addresses, msgs, d.Msg.Proposal, err = handle.MsgSubmitProposalV1(ctx, cfg.Codec, status, typedMsg)
		if err != nil {
			return d, err
		}
		if len(msgs) > 0 {
			d.Msg.Data["Messages"] = msgs
		}
	case *cosmosGovTypesV1Beta1.MsgSubmitProposal:
		var content any
		d.Msg.Type, d.Msg.Addresses, content, d.Msg.Proposal, err = handle.MsgSubmitProposalV1Beta(ctx, cfg.Codec, status, typedMsg)
		if err != nil {
			return d, err
		}
		if content != nil {
			d.Msg.Data["Content"] = content
		}
	case *cosmosGovTypesV1.MsgExecLegacyContent:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgExecLegacyContent(ctx, typedMsg)
	case *cosmosGovTypesV1.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVote(ctx, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVote(ctx, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgVoteWeighted:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteWeighted(ctx, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVoteWeighted:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteWeighted(ctx, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgDeposit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDeposit(ctx, typedMsg.Depositor)
	case *cosmosGovTypesV1Beta1.MsgDeposit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgDeposit(ctx, typedMsg.Depositor)

	// ibc module
	case *ibcTypes.MsgTransfer:
		d.Msg.Type, d.Msg.Addresses, err = handle.IBCTransfer(ctx, typedMsg)

	// crisis module
	case *crisisTypes.MsgVerifyInvariant:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVerifyInvariant(ctx, typedMsg)

	// evidence module
	case *evidenceTypes.MsgSubmitEvidence:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitEvidence(ctx, typedMsg)

	// nft module
	case *nft.MsgSend:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSendNFT(ctx, typedMsg)

	// group module
	case *group.MsgCreateGroup:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroup(ctx, typedMsg)
	case *group.MsgUpdateGroupMembers:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupMembers(ctx, typedMsg)
	case *group.MsgUpdateGroupAdmin:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupAdmin(ctx, typedMsg)
	case *group.MsgUpdateGroupMetadata:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupMetadata(ctx, typedMsg)
	case *group.MsgCreateGroupPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroupPolicy(ctx, typedMsg)
	case *group.MsgUpdateGroupPolicyAdmin:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyAdmin(ctx, typedMsg)
	case *group.MsgCreateGroupWithPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateGroupWithPolicy(ctx, typedMsg)
	case *group.MsgUpdateGroupPolicyDecisionPolicy:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyDecisionPolicy(ctx, typedMsg)
	case *group.MsgUpdateGroupPolicyMetadata:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateGroupPolicyMetadata(ctx, typedMsg)
	case *group.MsgSubmitProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitProposalGroup(ctx, typedMsg)
	case *group.MsgWithdrawProposal:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgWithdrawProposal(ctx, typedMsg)
	case *group.MsgVote:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgVoteGroup(ctx, typedMsg)
	case *group.MsgExec:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgExecGroup(ctx, typedMsg)
	case *group.MsgLeaveGroup:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgLeaveGroup(ctx, typedMsg)

	// upgrade module
	case *upgrade.MsgSoftwareUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSoftwareUpgrade(ctx, typedMsg)
	case *upgrade.MsgCancelUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCancelUpgrade(ctx, typedMsg)

	// interchainAccounts module
	case *interchainAccounts.MsgRegisterInterchainAccount:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterInterchainAccount(ctx, typedMsg)
	case *interchainAccounts.MsgSendTx:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSendTx(ctx, typedMsg)

	// fee module
	case *fee.MsgRegisterPayee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterPayee(ctx, typedMsg)
	case *fee.MsgRegisterCounterpartyPayee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRegisterCounterpartyPayee(ctx, typedMsg)
	case *fee.MsgPayPacketFee:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgPayPacketFee(ctx, typedMsg)
	case *fee.MsgPayPacketFeeAsync:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgPayPacketFeeAsync()

	// coreClient module
	case *coreClient.MsgCreateClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateClient(ctx, status, d.Msg.Data, typedMsg)
	case *coreClient.MsgUpdateClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateClient(ctx, status, d.Msg.Data, typedMsg)
	case *coreClient.MsgUpgradeClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpgradeClient(ctx, typedMsg)
	case *coreClient.MsgRecoverClient:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRecoverClient(ctx, typedMsg)
	case *coreClient.MsgIBCSoftwareUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgIBCSoftwareUpgrade(ctx, typedMsg)
	case *coreClient.MsgUpdateParams:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateParams(ctx, typedMsg)
	case *coreClient.MsgSubmitMisbehaviour: //nolint
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSubmitMisbehaviour(ctx, typedMsg)

	// coreConnection module
	case *coreConnection.MsgConnectionOpenInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenInit(ctx, typedMsg)
	case *coreConnection.MsgConnectionOpenTry:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenTry(ctx, typedMsg)
	case *coreConnection.MsgConnectionOpenAck:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenAck(ctx, typedMsg)
	case *coreConnection.MsgConnectionOpenConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgConnectionOpenConfirm(ctx, typedMsg)

	// coreChannel module
	case *coreChannel.MsgChannelOpenInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenInit(ctx, typedMsg)
	case *coreChannel.MsgChannelOpenTry:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenTry(ctx, typedMsg)
	case *coreChannel.MsgChannelOpenAck:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenAck(ctx, typedMsg)
	case *coreChannel.MsgChannelOpenConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelOpenConfirm(ctx, typedMsg)
	case *coreChannel.MsgChannelCloseInit:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelCloseInit(ctx, typedMsg)
	case *coreChannel.MsgChannelCloseConfirm:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgChannelCloseConfirm(ctx, typedMsg)
	case *coreChannel.MsgRecvPacket:
		d.Msg.Type, d.Msg.Addresses, d.Msg.IbcTransfer, d.Msg.IbcChannel, err = handle.MsgRecvPacket(ctx, cfg.Codec, d.Msg.Data, typedMsg)
	case *coreChannel.MsgTimeout:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgTimeout(ctx, typedMsg)
	case *coreChannel.MsgTimeoutOnClose:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgTimeoutOnClose(ctx, typedMsg)
	case *coreChannel.MsgAcknowledgement:
		d.Msg.Type, d.Msg.Addresses, d.Msg.IbcTransfer, d.Msg.IbcChannel, err = handle.MsgAcknowledgement(ctx, cfg.Codec, d.Msg.Data, typedMsg)

	// signal module
	case *appSignalTypes.MsgSignalVersion:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSignalVersion(ctx, typedMsg.GetValidatorAddress())
	case *appSignalTypes.MsgTryUpgrade:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgTryUpgrade(ctx, typedMsg.GetSigner())

	// hyperlane
	case *hyperlaneCore.MsgCreateMailbox:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateMailbox(ctx, typedMsg)
	case *hyperlaneCore.MsgProcessMessage:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgProcessMessage(ctx, typedMsg)
	case *hyperlaneCore.MsgSetMailbox:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetMailbox(ctx, typedMsg)
	case *hyperlaneWarp.MsgCreateCollateralToken:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateCollateralToken(ctx, typedMsg)
	case *hyperlaneWarp.MsgCreateSyntheticToken:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateSyntheticToken(ctx, typedMsg)
	case *hyperlaneWarp.MsgSetToken:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetToken(ctx, typedMsg)
	case *hyperlaneWarp.MsgEnrollRemoteRouter:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgEnrollRemoteRouter(ctx, typedMsg)
	case *hyperlaneWarp.MsgUnrollRemoteRouter:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUnrollRemoteRouter(ctx, typedMsg)
	case *hyperlaneWarp.MsgRemoteTransfer:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgRemoteTransfer(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgClaim:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgClaim(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgCreateIgp:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateIgp(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgSetIgpOwner:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetIgpOwner(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgPayForGas:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgPayForGas(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgSetDestinationGasConfig:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgSetDestinationGasConfig(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgCreateMerkleTreeHook:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateMerkleTreeHook(ctx, typedMsg)
	case *hyperlanePostDispatch.MsgCreateNoopHook:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgCreateNoopHook(ctx, typedMsg)

	case *minfeeTypes.MsgUpdateMinfeeParams:
		d.Msg.Type, d.Msg.Addresses, err = handle.MsgUpdateMinfeeParams(ctx, typedMsg)

	default:
		log.Err(errors.New("unknown message type")).Msgf("got type %T", msg)
		d.Msg.Type = storageTypes.MsgUnknown
	}

	if err != nil {
		err = errors.Wrapf(err, "while decoding msg(%T) on position=%d", msg, position)
	}

	if d.Msg.Type != storageTypes.MsgUnknown {
		if m, ok := msg.(measurable); ok {
			d.Msg.Size = m.Size()
		} else {
			return d, errors.Errorf("message %T does not implement Size method: %##v", msg, msg)
		}
	}

	d.Addresses = append(d.Addresses, d.Msg.Addresses...)
	return
}

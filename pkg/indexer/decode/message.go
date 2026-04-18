// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	circuitTypes "cosmossdk.io/x/circuit/types"
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
	interchainAccountsHost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	fee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibcTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	coreClient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	coreConnection "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	coreChannel "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/rs/zerolog/log"

	cosmosFeegrant "cosmossdk.io/x/feegrant"
	hyperlaneICS "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	hyperlanePostDispatch "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	hyperlaneCore "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	hyperlaneWarp "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/legacy"
	appBlobTypes "github.com/celestiaorg/celestia-app/v8/x/blob/types"
	fwdTypes "github.com/celestiaorg/celestia-app/v8/x/forwarding/types"
	minfeeTypes "github.com/celestiaorg/celestia-app/v8/x/minfee/types"
	appSignalTypes "github.com/celestiaorg/celestia-app/v8/x/signal/types"
	zkismTypes "github.com/celestiaorg/celestia-app/v8/x/zkism/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosVestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	cosmosGovTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	cosmosSlashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
)

type DecodedMsg struct {
	Msg       storage.Message
	BlobsSize int64
	BlobLogs  []*storage.BlobLog
}

type measurable interface {
	Size() int
}

func Message(
	ctx *context.Context,
	msg cosmosTypes.Msg,
	position int,
	status storageTypes.Status,
	txId uint64,
) (d DecodedMsg, err error) {
	d.Msg.Position = int64(position)
	d.Msg.Data, err = msgToMap(msg)
	if err != nil {
		return d, errors.Wrap(err, "msg to map")
	}
	d.Msg.Height = ctx.Block.Height
	d.Msg.Time = ctx.Block.Time
	d.Msg.TxId = txId

	if err := d.Msg.SetId(ctx.GetMsgPosition()); err != nil {
		return d, err
	}

	switch typedMsg := msg.(type) {

	// distribution module
	case *cosmosDistributionTypes.MsgSetWithdrawAddress:
		d.Msg.Type, err = handle.MsgSetWithdrawAddress(ctx, d.Msg.Id, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawDelegatorReward:
		d.Msg.Type, err = handle.MsgWithdrawDelegatorReward(ctx, d.Msg.Id, typedMsg)
	case *cosmosDistributionTypes.MsgWithdrawValidatorCommission:
		d.Msg.Type, err = handle.MsgWithdrawValidatorCommission(ctx, d.Msg.Id, typedMsg)
	case *cosmosDistributionTypes.MsgFundCommunityPool:
		d.Msg.Type, err = handle.MsgFundCommunityPool(ctx, d.Msg.Id, typedMsg)
	case *cosmosDistributionTypes.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsDistr(ctx, d.Msg.Id, typedMsg)

	// staking module
	case *cosmosStakingTypes.MsgCreateValidator:
		d.Msg.Type, d.Msg.Validators, err = handle.MsgCreateValidator(ctx, status, d.Msg.Id, typedMsg)
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
		d.Msg.Type, d.Msg.Validators, err = handle.MsgEditValidator(ctx, status, d.Msg.Id, typedMsg)
	case *cosmosStakingTypes.MsgDelegate:
		d.Msg.Type, err = handle.MsgDelegate(ctx, d.Msg.Id, typedMsg)
	case *cosmosStakingTypes.MsgBeginRedelegate:
		d.Msg.Type, err = handle.MsgBeginRedelegate(ctx, d.Msg.Id, typedMsg)
	case *cosmosStakingTypes.MsgUndelegate:
		d.Msg.Type, err = handle.MsgUndelegate(ctx, d.Msg.Id, typedMsg)
	case *cosmosStakingTypes.MsgCancelUnbondingDelegation:
		d.Msg.Type, err = handle.MsgCancelUnbondingDelegation(ctx, d.Msg.Id, typedMsg)
	case *cosmosStakingTypes.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsStaking(ctx, d.Msg.Id, typedMsg)

	// slashing module
	case *cosmosSlashingTypes.MsgUnjail:
		d.Msg.Type, err = handle.MsgUnjail(ctx, d.Msg.Id, typedMsg)
	case *cosmosSlashingTypes.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsSlashing(ctx, d.Msg.Id, typedMsg)

	// bank module
	case *cosmosBankTypes.MsgSend:
		d.Msg.Type, err = handle.MsgSend(ctx, d.Msg.Id, typedMsg)
	case *cosmosBankTypes.MsgMultiSend:
		d.Msg.Type, err = handle.MsgMultiSend(ctx, d.Msg.Id, typedMsg)
	case *cosmosBankTypes.MsgSetSendEnabled:
		d.Msg.Type, err = handle.MsgSetSendEnabled(ctx, d.Msg.Id, typedMsg)
	case *cosmosBankTypes.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsBank(ctx, d.Msg.Id, typedMsg)

	// vesting module
	case *cosmosVestingTypes.MsgCreateVestingAccount:
		d.Msg.Type, err = handle.MsgCreateVestingAccount(ctx, status, txId, d.Msg.Id, typedMsg)
	case *cosmosVestingTypes.MsgCreatePermanentLockedAccount:
		d.Msg.Type, err = handle.MsgCreatePermanentLockedAccount(ctx, status, txId, d.Msg.Id, typedMsg)
	case *cosmosVestingTypes.MsgCreatePeriodicVestingAccount:
		d.Msg.Type, err = handle.MsgCreatePeriodicVestingAccount(ctx, status, txId, d.Msg.Id, typedMsg)

	// blob module
	case *appBlobTypes.MsgPayForBlobs:
		d.Msg.Type, d.BlobLogs, d.BlobsSize, err = handle.MsgPayForBlobs(ctx, status, d.Msg.Id, txId, typedMsg)
	case *appBlobTypes.MsgUpdateBlobParams:
		d.Msg.Type, err = handle.MsgUpdateBlobParams(ctx, status, d.Msg.Id, typedMsg)

	// feegrant module
	case *cosmosFeegrant.MsgGrantAllowance:
		d.Msg.Type, err = handle.MsgGrantAllowance(ctx, status, d.Msg.Id, typedMsg)
	case *cosmosFeegrant.MsgRevokeAllowance:
		d.Msg.Type, err = handle.MsgRevokeAllowance(ctx, status, d.Msg.Id, typedMsg)

	// qgb module
	case *legacy.MsgRegisterEVMAddress:
		d.Msg.Type, err = handle.MsgRegisterEVMAddress(ctx, d.Msg.Id, typedMsg)

	// authz module
	case *authz.MsgPruneExpiredGrants:
		d.Msg.Type, err = handle.MsgPruneExpiredGrants(ctx, d.Msg.Id, typedMsg)
	case *authz.MsgGrant:
		d.Msg.Type, err = handle.MsgGrant(ctx, status, d.Msg.Id, typedMsg)
	case *authz.MsgExec:
		d.Msg.Type, d.Msg.InternalMsgs, err = handle.MsgExec(ctx, status, d.Msg.Id, typedMsg)
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
			m, mapErr := msgToMap(msg)
			if mapErr != nil {
				return d, errors.Wrap(mapErr, "msg to map")
			}
			msgs = append(msgs, map[string]any(m))
		}
		d.Msg.Data["Msgs"] = msgs

	case *authz.MsgRevoke:
		d.Msg.Type, err = handle.MsgRevoke(ctx, status, d.Msg.Id, typedMsg)

	// gov module
	case *cosmosGovTypesV1.MsgSubmitProposal:
		var msgs []any
		d.Msg.Type, msgs, d.Msg.Proposal, err = handle.MsgSubmitProposalV1(ctx, cfg.Codec, status, d.Msg.Id, typedMsg)
		if err != nil {
			return d, err
		}
		if len(msgs) > 0 {
			d.Msg.Data["Messages"] = msgs
		}
	case *cosmosGovTypesV1Beta1.MsgSubmitProposal:
		var content any
		d.Msg.Type, content, d.Msg.Proposal, err = handle.MsgSubmitProposalV1Beta(ctx, cfg.Codec, status, d.Msg.Id, typedMsg)
		if err != nil {
			return d, err
		}
		if content != nil {
			d.Msg.Data["Content"] = content
		}
	case *cosmosGovTypesV1.MsgExecLegacyContent:
		d.Msg.Type, err = handle.MsgExecLegacyContent(ctx, d.Msg.Id, typedMsg)
	case *cosmosGovTypesV1.MsgVote:
		d.Msg.Type, err = handle.MsgVote(ctx, d.Msg.Id, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVote:
		d.Msg.Type, err = handle.MsgVote(ctx, d.Msg.Id, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgVoteWeighted:
		d.Msg.Type, err = handle.MsgVoteWeighted(ctx, d.Msg.Id, typedMsg.Voter)
	case *cosmosGovTypesV1Beta1.MsgVoteWeighted:
		d.Msg.Type, err = handle.MsgVoteWeighted(ctx, d.Msg.Id, typedMsg.Voter)
	case *cosmosGovTypesV1.MsgDeposit:
		d.Msg.Type, err = handle.MsgDeposit(ctx, d.Msg.Id, typedMsg.Depositor)
	case *cosmosGovTypesV1Beta1.MsgDeposit:
		d.Msg.Type, err = handle.MsgDeposit(ctx, d.Msg.Id, typedMsg.Depositor)
	case *cosmosGovTypesV1.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsGov(ctx, d.Msg.Id, typedMsg)

	// ibc module
	case *ibcTypes.MsgTransfer:
		d.Msg.Type, err = handle.IBCTransfer(ctx, d.Msg.Id, typedMsg)

	// crisis module
	case *crisisTypes.MsgVerifyInvariant:
		d.Msg.Type, err = handle.MsgVerifyInvariant(ctx, d.Msg.Id, typedMsg)

	// evidence module
	case *evidenceTypes.MsgSubmitEvidence:
		d.Msg.Type, err = handle.MsgSubmitEvidence(ctx, d.Msg.Id, typedMsg)

	// nft module
	case *nft.MsgSend:
		d.Msg.Type, err = handle.MsgSendNFT(ctx, d.Msg.Id, typedMsg)

	// group module
	case *group.MsgCreateGroup:
		d.Msg.Type, err = handle.MsgCreateGroup(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupMembers:
		d.Msg.Type, err = handle.MsgUpdateGroupMembers(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupAdmin:
		d.Msg.Type, err = handle.MsgUpdateGroupAdmin(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupMetadata:
		d.Msg.Type, err = handle.MsgUpdateGroupMetadata(ctx, d.Msg.Id, typedMsg)
	case *group.MsgCreateGroupPolicy:
		d.Msg.Type, err = handle.MsgCreateGroupPolicy(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupPolicyAdmin:
		d.Msg.Type, err = handle.MsgUpdateGroupPolicyAdmin(ctx, d.Msg.Id, typedMsg)
	case *group.MsgCreateGroupWithPolicy:
		d.Msg.Type, err = handle.MsgCreateGroupWithPolicy(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupPolicyDecisionPolicy:
		d.Msg.Type, err = handle.MsgUpdateGroupPolicyDecisionPolicy(ctx, d.Msg.Id, typedMsg)
	case *group.MsgUpdateGroupPolicyMetadata:
		d.Msg.Type, err = handle.MsgUpdateGroupPolicyMetadata(ctx, d.Msg.Id, typedMsg)
	case *group.MsgSubmitProposal:
		d.Msg.Type, err = handle.MsgSubmitProposalGroup(ctx, d.Msg.Id, typedMsg)
	case *group.MsgWithdrawProposal:
		d.Msg.Type, err = handle.MsgWithdrawProposal(ctx, d.Msg.Id, typedMsg)
	case *group.MsgVote:
		d.Msg.Type, err = handle.MsgVoteGroup(ctx, d.Msg.Id, typedMsg)
	case *group.MsgExec:
		d.Msg.Type, err = handle.MsgExecGroup(ctx, d.Msg.Id, typedMsg)
	case *group.MsgLeaveGroup:
		d.Msg.Type, err = handle.MsgLeaveGroup(ctx, d.Msg.Id, typedMsg)

	// upgrade module
	case *upgrade.MsgSoftwareUpgrade:
		d.Msg.Type, err = handle.MsgSoftwareUpgrade(ctx, d.Msg.Id, typedMsg)
	case *upgrade.MsgCancelUpgrade:
		d.Msg.Type, err = handle.MsgCancelUpgrade(ctx, d.Msg.Id, typedMsg)

	// interchainAccounts module
	case *interchainAccounts.MsgRegisterInterchainAccount:
		d.Msg.Type, err = handle.MsgRegisterInterchainAccount(ctx, d.Msg.Id, typedMsg)
	case *interchainAccounts.MsgSendTx:
		d.Msg.Type, err = handle.MsgSendTx(ctx, d.Msg.Id, typedMsg)
	case *interchainAccounts.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsIcaController(ctx, d.Msg.Id, typedMsg)
	case *interchainAccountsHost.MsgModuleQuerySafe:
		d.Msg.Type, err = handle.MsgModuleQuerySafe(ctx, d.Msg.Id, typedMsg)
	case *interchainAccountsHost.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsIcaHost(ctx, d.Msg.Id, typedMsg)

	// fee module
	case *fee.MsgRegisterPayee:
		d.Msg.Type, err = handle.MsgRegisterPayee(ctx, d.Msg.Id, typedMsg)
	case *fee.MsgRegisterCounterpartyPayee:
		d.Msg.Type, err = handle.MsgRegisterCounterpartyPayee(ctx, d.Msg.Id, typedMsg)
	case *fee.MsgPayPacketFee:
		d.Msg.Type, err = handle.MsgPayPacketFee(ctx, d.Msg.Id, typedMsg)
	case *fee.MsgPayPacketFeeAsync:
		d.Msg.Type, err = handle.MsgPayPacketFeeAsync()

	// coreClient module
	case *coreClient.MsgCreateClient:
		d.Msg.Type, err = handle.MsgCreateClient(ctx, status, d.Msg.Data, d.Msg.Id, typedMsg)
	case *coreClient.MsgUpdateClient:
		d.Msg.Type, err = handle.MsgUpdateClient(ctx, status, d.Msg.Data, d.Msg.Id, typedMsg)
	case *coreClient.MsgUpgradeClient:
		d.Msg.Type, err = handle.MsgUpgradeClient(ctx, d.Msg.Id, typedMsg)
	case *coreClient.MsgRecoverClient:
		d.Msg.Type, err = handle.MsgRecoverClient(ctx, d.Msg.Id, typedMsg)
	case *coreClient.MsgIBCSoftwareUpgrade:
		d.Msg.Type, err = handle.MsgIBCSoftwareUpgrade(ctx, d.Msg.Id, typedMsg)
	case *coreClient.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParams(ctx, d.Msg.Id, typedMsg)
	case *coreClient.MsgSubmitMisbehaviour: //nolint
		d.Msg.Type, err = handle.MsgSubmitMisbehaviour(ctx, d.Msg.Id, typedMsg)

	// coreConnection module
	case *coreConnection.MsgConnectionOpenInit:
		d.Msg.Type, err = handle.MsgConnectionOpenInit(ctx, d.Msg.Id, typedMsg)
	case *coreConnection.MsgConnectionOpenTry:
		d.Msg.Type, err = handle.MsgConnectionOpenTry(ctx, d.Msg.Id, typedMsg)
	case *coreConnection.MsgConnectionOpenAck:
		d.Msg.Type, err = handle.MsgConnectionOpenAck(ctx, d.Msg.Id, typedMsg)
	case *coreConnection.MsgConnectionOpenConfirm:
		d.Msg.Type, err = handle.MsgConnectionOpenConfirm(ctx, d.Msg.Id, typedMsg)
	case *coreConnection.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsConnection(ctx, d.Msg.Id, typedMsg)

	// coreChannel module
	case *coreChannel.MsgChannelOpenInit:
		d.Msg.Type, err = handle.MsgChannelOpenInit(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgChannelOpenTry:
		d.Msg.Type, err = handle.MsgChannelOpenTry(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgChannelOpenAck:
		d.Msg.Type, err = handle.MsgChannelOpenAck(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgChannelOpenConfirm:
		d.Msg.Type, err = handle.MsgChannelOpenConfirm(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgChannelCloseInit:
		d.Msg.Type, err = handle.MsgChannelCloseInit(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgChannelCloseConfirm:
		d.Msg.Type, err = handle.MsgChannelCloseConfirm(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgRecvPacket:
		d.Msg.Type, err = handle.MsgRecvPacket(ctx, status, cfg.Codec, d.Msg.Data, txId, d.Msg.Id, typedMsg)
	case *coreChannel.MsgTimeout:
		d.Msg.Type, err = handle.MsgTimeout(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgTimeoutOnClose:
		d.Msg.Type, err = handle.MsgTimeoutOnClose(ctx, d.Msg.Id, typedMsg)
	case *coreChannel.MsgAcknowledgement:
		d.Msg.Type, err = handle.MsgAcknowledgement(ctx, status, cfg.Codec, d.Msg.Data, txId, d.Msg.Id, typedMsg)
	case *coreChannel.MsgUpdateParams:
		d.Msg.Type, err = handle.MsgUpdateParamsChannel(ctx, d.Msg.Id, typedMsg)

	// signal module
	case *appSignalTypes.MsgSignalVersion:
		d.Msg.Type, err = handle.MsgSignalVersion(ctx, status, txId, d.Msg.Id, typedMsg)
	case *appSignalTypes.MsgTryUpgrade:
		d.Msg.Type, err = handle.MsgTryUpgrade(ctx, status, txId, d.Msg.Id, typedMsg)

	// hyperlane
	case *hyperlaneCore.MsgCreateMailbox:
		d.Msg.Type, err = handle.MsgCreateMailbox(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneCore.MsgProcessMessage:
		d.Msg.Type, err = handle.MsgProcessMessage(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneCore.MsgSetMailbox:
		d.Msg.Type, err = handle.MsgSetMailbox(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgCreateCollateralToken:
		d.Msg.Type, err = handle.MsgCreateCollateralToken(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgCreateSyntheticToken:
		d.Msg.Type, err = handle.MsgCreateSyntheticToken(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgSetToken:
		d.Msg.Type, err = handle.MsgSetToken(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgEnrollRemoteRouter:
		d.Msg.Type, err = handle.MsgEnrollRemoteRouter(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgUnrollRemoteRouter:
		d.Msg.Type, err = handle.MsgUnrollRemoteRouter(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneWarp.MsgRemoteTransfer:
		d.Msg.Type, err = handle.MsgRemoteTransfer(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgClaim:
		d.Msg.Type, err = handle.MsgClaim(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgCreateIgp:
		d.Msg.Type, err = handle.MsgCreateIgp(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgSetIgpOwner:
		d.Msg.Type, err = handle.MsgSetIgpOwner(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgPayForGas:
		d.Msg.Type, err = handle.MsgPayForGas(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgSetDestinationGasConfig:
		d.Msg.Type, err = handle.MsgSetDestinationGasConfig(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgCreateMerkleTreeHook:
		d.Msg.Type, err = handle.MsgCreateMerkleTreeHook(ctx, d.Msg.Id, typedMsg)
	case *hyperlanePostDispatch.MsgCreateNoopHook:
		d.Msg.Type, err = handle.MsgCreateNoopHook(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneICS.MsgAnnounceValidator:
		d.Msg.Type, err = handle.MsgAnnounceValidator(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneICS.MsgCreateMerkleRootMultisigIsm:
		d.Msg.Type, err = handle.MsgCreateMerkleRootMultisigIsm(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneICS.MsgCreateMessageIdMultisigIsm:
		d.Msg.Type, err = handle.MsgCreateMessageIdMultisigIsm(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneICS.MsgCreateNoopIsm:
		d.Msg.Type, err = handle.MsgCreateNoopIsm(ctx, d.Msg.Id, typedMsg)
	case *hyperlaneICS.MsgCreateRoutingIsm:
		d.Msg.Type, err = handle.MsgCreateRoutingIsm(ctx, d.Msg.Id, typedMsg)

	case *minfeeTypes.MsgUpdateMinfeeParams:
		d.Msg.Type, err = handle.MsgUpdateMinfeeParams(ctx, d.Msg.Id, typedMsg)

	// circuit
	case *circuitTypes.MsgAuthorizeCircuitBreaker:
		d.Msg.Type, err = handle.MsgAuthorizeCircuitBreaker(ctx, d.Msg.Id, typedMsg)
	case *circuitTypes.MsgResetCircuitBreaker:
		d.Msg.Type, err = handle.MsgResetCircuitBreaker(ctx, d.Msg.Id, typedMsg)
	case *circuitTypes.MsgTripCircuitBreaker:
		d.Msg.Type, err = handle.MsgTripCircuitBreaker(ctx, d.Msg.Id, typedMsg)

	// forwarding
	case *fwdTypes.MsgForward:
		d.Msg.Type, err = handle.MsgForward(ctx, d.Msg.Id, typedMsg)

	// zkism
	case *zkismTypes.MsgCreateInterchainSecurityModule:
		d.Msg.Type, err = handle.MsgCreateInterchainSecurityModule(ctx, d.Msg.Id, typedMsg)
	case *zkismTypes.MsgUpdateInterchainSecurityModule:
		d.Msg.Type, err = handle.MsgUpdateInterchainSecurityModule(ctx, d.Msg.Id, typedMsg)
	case *zkismTypes.MsgSubmitMessages:
		d.Msg.Type, err = handle.MsgSubmitMessages(ctx, d.Msg.Id, typedMsg)

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

	return
}

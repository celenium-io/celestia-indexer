// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.7
// Revision: bf63e108589bbd2327b13ec2c5da532aad234029
// Build Date: 2023-07-25T23:27:55Z
// Built By: goreleaser

package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// MsgUnknown is a MsgType of type MsgUnknown.
	MsgUnknown MsgType = "MsgUnknown"
	// MsgSetWithdrawAddress is a MsgType of type MsgSetWithdrawAddress.
	MsgSetWithdrawAddress MsgType = "MsgSetWithdrawAddress"
	// MsgWithdrawDelegatorReward is a MsgType of type MsgWithdrawDelegatorReward.
	MsgWithdrawDelegatorReward MsgType = "MsgWithdrawDelegatorReward"
	// MsgWithdrawValidatorCommission is a MsgType of type MsgWithdrawValidatorCommission.
	MsgWithdrawValidatorCommission MsgType = "MsgWithdrawValidatorCommission"
	// MsgFundCommunityPool is a MsgType of type MsgFundCommunityPool.
	MsgFundCommunityPool MsgType = "MsgFundCommunityPool"
	// MsgCreateValidator is a MsgType of type MsgCreateValidator.
	MsgCreateValidator MsgType = "MsgCreateValidator"
	// MsgEditValidator is a MsgType of type MsgEditValidator.
	MsgEditValidator MsgType = "MsgEditValidator"
	// MsgDelegate is a MsgType of type MsgDelegate.
	MsgDelegate MsgType = "MsgDelegate"
	// MsgBeginRedelegate is a MsgType of type MsgBeginRedelegate.
	MsgBeginRedelegate MsgType = "MsgBeginRedelegate"
	// MsgUndelegate is a MsgType of type MsgUndelegate.
	MsgUndelegate MsgType = "MsgUndelegate"
	// MsgCancelUnbondingDelegation is a MsgType of type MsgCancelUnbondingDelegation.
	MsgCancelUnbondingDelegation MsgType = "MsgCancelUnbondingDelegation"
	// MsgUnjail is a MsgType of type MsgUnjail.
	MsgUnjail MsgType = "MsgUnjail"
	// MsgSend is a MsgType of type MsgSend.
	MsgSend MsgType = "MsgSend"
	// MsgMultiSend is a MsgType of type MsgMultiSend.
	MsgMultiSend MsgType = "MsgMultiSend"
	// MsgCreateVestingAccount is a MsgType of type MsgCreateVestingAccount.
	MsgCreateVestingAccount MsgType = "MsgCreateVestingAccount"
	// MsgCreatePermanentLockedAccount is a MsgType of type MsgCreatePermanentLockedAccount.
	MsgCreatePermanentLockedAccount MsgType = "MsgCreatePermanentLockedAccount"
	// MsgCreatePeriodicVestingAccount is a MsgType of type MsgCreatePeriodicVestingAccount.
	MsgCreatePeriodicVestingAccount MsgType = "MsgCreatePeriodicVestingAccount"
	// MsgPayForBlobs is a MsgType of type MsgPayForBlobs.
	MsgPayForBlobs MsgType = "MsgPayForBlobs"
	// MsgGrant is a MsgType of type MsgGrant.
	MsgGrant MsgType = "MsgGrant"
	// MsgExec is a MsgType of type MsgExec.
	MsgExec MsgType = "MsgExec"
	// MsgRevoke is a MsgType of type MsgRevoke.
	MsgRevoke MsgType = "MsgRevoke"
	// MsgGrantAllowance is a MsgType of type MsgGrantAllowance.
	MsgGrantAllowance MsgType = "MsgGrantAllowance"
	// MsgRevokeAllowance is a MsgType of type MsgRevokeAllowance.
	MsgRevokeAllowance MsgType = "MsgRevokeAllowance"
	// MsgRegisterEVMAddress is a MsgType of type MsgRegisterEVMAddress.
	MsgRegisterEVMAddress MsgType = "MsgRegisterEVMAddress"
	// MsgSubmitProposal is a MsgType of type MsgSubmitProposal.
	MsgSubmitProposal MsgType = "MsgSubmitProposal"
	// MsgExecLegacyContent is a MsgType of type MsgExecLegacyContent.
	MsgExecLegacyContent MsgType = "MsgExecLegacyContent"
	// MsgVote is a MsgType of type MsgVote.
	MsgVote MsgType = "MsgVote"
	// MsgVoteWeighted is a MsgType of type MsgVoteWeighted.
	MsgVoteWeighted MsgType = "MsgVoteWeighted"
	// MsgDeposit is a MsgType of type MsgDeposit.
	MsgDeposit MsgType = "MsgDeposit"
	// IBCTransfer is a MsgType of type IBCTransfer.
	IBCTransfer MsgType = "IBCTransfer"
	// MsgVerifyInvariant is a MsgType of type MsgVerifyInvariant.
	MsgVerifyInvariant MsgType = "MsgVerifyInvariant"
	// MsgSubmitEvidence is a MsgType of type MsgSubmitEvidence.
	MsgSubmitEvidence MsgType = "MsgSubmitEvidence"
	// MsgSendNFT is a MsgType of type MsgSendNFT.
	MsgSendNFT MsgType = "MsgSendNFT"
	// MsgCreateGroup is a MsgType of type MsgCreateGroup.
	MsgCreateGroup MsgType = "MsgCreateGroup"
	// MsgUpdateGroupMembers is a MsgType of type MsgUpdateGroupMembers.
	MsgUpdateGroupMembers MsgType = "MsgUpdateGroupMembers"
	// MsgUpdateGroupAdmin is a MsgType of type MsgUpdateGroupAdmin.
	MsgUpdateGroupAdmin MsgType = "MsgUpdateGroupAdmin"
	// MsgUpdateGroupMetadata is a MsgType of type MsgUpdateGroupMetadata.
	MsgUpdateGroupMetadata MsgType = "MsgUpdateGroupMetadata"
	// MsgCreateGroupPolicy is a MsgType of type MsgCreateGroupPolicy.
	MsgCreateGroupPolicy MsgType = "MsgCreateGroupPolicy"
	// MsgUpdateGroupPolicyAdmin is a MsgType of type MsgUpdateGroupPolicyAdmin.
	MsgUpdateGroupPolicyAdmin MsgType = "MsgUpdateGroupPolicyAdmin"
	// MsgCreateGroupWithPolicy is a MsgType of type MsgCreateGroupWithPolicy.
	MsgCreateGroupWithPolicy MsgType = "MsgCreateGroupWithPolicy"
	// MsgUpdateGroupPolicyDecisionPolicy is a MsgType of type MsgUpdateGroupPolicyDecisionPolicy.
	MsgUpdateGroupPolicyDecisionPolicy MsgType = "MsgUpdateGroupPolicyDecisionPolicy"
	// MsgUpdateGroupPolicyMetadata is a MsgType of type MsgUpdateGroupPolicyMetadata.
	MsgUpdateGroupPolicyMetadata MsgType = "MsgUpdateGroupPolicyMetadata"
	// MsgSubmitProposalGroup is a MsgType of type MsgSubmitProposalGroup.
	MsgSubmitProposalGroup MsgType = "MsgSubmitProposalGroup"
	// MsgWithdrawProposal is a MsgType of type MsgWithdrawProposal.
	MsgWithdrawProposal MsgType = "MsgWithdrawProposal"
	// MsgVoteGroup is a MsgType of type MsgVoteGroup.
	MsgVoteGroup MsgType = "MsgVoteGroup"
	// MsgExecGroup is a MsgType of type MsgExecGroup.
	MsgExecGroup MsgType = "MsgExecGroup"
	// MsgLeaveGroup is a MsgType of type MsgLeaveGroup.
	MsgLeaveGroup MsgType = "MsgLeaveGroup"
	// MsgSoftwareUpgrade is a MsgType of type MsgSoftwareUpgrade.
	MsgSoftwareUpgrade MsgType = "MsgSoftwareUpgrade"
	// MsgCancelUpgrade is a MsgType of type MsgCancelUpgrade.
	MsgCancelUpgrade MsgType = "MsgCancelUpgrade"
	// MsgRegisterInterchainAccount is a MsgType of type MsgRegisterInterchainAccount.
	MsgRegisterInterchainAccount MsgType = "MsgRegisterInterchainAccount"
	// MsgSendTx is a MsgType of type MsgSendTx.
	MsgSendTx MsgType = "MsgSendTx"
	// MsgRegisterPayee is a MsgType of type MsgRegisterPayee.
	MsgRegisterPayee MsgType = "MsgRegisterPayee"
	// MsgRegisterCounterpartyPayee is a MsgType of type MsgRegisterCounterpartyPayee.
	MsgRegisterCounterpartyPayee MsgType = "MsgRegisterCounterpartyPayee"
	// MsgPayPacketFee is a MsgType of type MsgPayPacketFee.
	MsgPayPacketFee MsgType = "MsgPayPacketFee"
	// MsgPayPacketFeeAsync is a MsgType of type MsgPayPacketFeeAsync.
	MsgPayPacketFeeAsync MsgType = "MsgPayPacketFeeAsync"
	// MsgTransfer is a MsgType of type MsgTransfer.
	MsgTransfer MsgType = "MsgTransfer"
	// MsgCreateClient is a MsgType of type MsgCreateClient.
	MsgCreateClient MsgType = "MsgCreateClient"
	// MsgUpdateClient is a MsgType of type MsgUpdateClient.
	MsgUpdateClient MsgType = "MsgUpdateClient"
	// MsgUpgradeClient is a MsgType of type MsgUpgradeClient.
	MsgUpgradeClient MsgType = "MsgUpgradeClient"
	// MsgSubmitMisbehaviour is a MsgType of type MsgSubmitMisbehaviour.
	MsgSubmitMisbehaviour MsgType = "MsgSubmitMisbehaviour"
	// MsgConnectionOpenInit is a MsgType of type MsgConnectionOpenInit.
	MsgConnectionOpenInit MsgType = "MsgConnectionOpenInit"
	// MsgConnectionOpenTry is a MsgType of type MsgConnectionOpenTry.
	MsgConnectionOpenTry MsgType = "MsgConnectionOpenTry"
	// MsgConnectionOpenAck is a MsgType of type MsgConnectionOpenAck.
	MsgConnectionOpenAck MsgType = "MsgConnectionOpenAck"
	// MsgConnectionOpenConfirm is a MsgType of type MsgConnectionOpenConfirm.
	MsgConnectionOpenConfirm MsgType = "MsgConnectionOpenConfirm"
	// MsgChannelOpenInit is a MsgType of type MsgChannelOpenInit.
	MsgChannelOpenInit MsgType = "MsgChannelOpenInit"
	// MsgChannelOpenTry is a MsgType of type MsgChannelOpenTry.
	MsgChannelOpenTry MsgType = "MsgChannelOpenTry"
	// MsgChannelOpenAck is a MsgType of type MsgChannelOpenAck.
	MsgChannelOpenAck MsgType = "MsgChannelOpenAck"
	// MsgChannelOpenConfirm is a MsgType of type MsgChannelOpenConfirm.
	MsgChannelOpenConfirm MsgType = "MsgChannelOpenConfirm"
	// MsgChannelCloseInit is a MsgType of type MsgChannelCloseInit.
	MsgChannelCloseInit MsgType = "MsgChannelCloseInit"
	// MsgChannelCloseConfirm is a MsgType of type MsgChannelCloseConfirm.
	MsgChannelCloseConfirm MsgType = "MsgChannelCloseConfirm"
	// MsgRecvPacket is a MsgType of type MsgRecvPacket.
	MsgRecvPacket MsgType = "MsgRecvPacket"
	// MsgTimeout is a MsgType of type MsgTimeout.
	MsgTimeout MsgType = "MsgTimeout"
	// MsgTimeoutOnClose is a MsgType of type MsgTimeoutOnClose.
	MsgTimeoutOnClose MsgType = "MsgTimeoutOnClose"
	// MsgAcknowledgement is a MsgType of type MsgAcknowledgement.
	MsgAcknowledgement MsgType = "MsgAcknowledgement"
)

var ErrInvalidMsgType = fmt.Errorf("not a valid MsgType, try [%s]", strings.Join(_MsgTypeNames, ", "))

var _MsgTypeNames = []string{
	string(MsgUnknown),
	string(MsgSetWithdrawAddress),
	string(MsgWithdrawDelegatorReward),
	string(MsgWithdrawValidatorCommission),
	string(MsgFundCommunityPool),
	string(MsgCreateValidator),
	string(MsgEditValidator),
	string(MsgDelegate),
	string(MsgBeginRedelegate),
	string(MsgUndelegate),
	string(MsgCancelUnbondingDelegation),
	string(MsgUnjail),
	string(MsgSend),
	string(MsgMultiSend),
	string(MsgCreateVestingAccount),
	string(MsgCreatePermanentLockedAccount),
	string(MsgCreatePeriodicVestingAccount),
	string(MsgPayForBlobs),
	string(MsgGrant),
	string(MsgExec),
	string(MsgRevoke),
	string(MsgGrantAllowance),
	string(MsgRevokeAllowance),
	string(MsgRegisterEVMAddress),
	string(MsgSubmitProposal),
	string(MsgExecLegacyContent),
	string(MsgVote),
	string(MsgVoteWeighted),
	string(MsgDeposit),
	string(IBCTransfer),
	string(MsgVerifyInvariant),
	string(MsgSubmitEvidence),
	string(MsgSendNFT),
	string(MsgCreateGroup),
	string(MsgUpdateGroupMembers),
	string(MsgUpdateGroupAdmin),
	string(MsgUpdateGroupMetadata),
	string(MsgCreateGroupPolicy),
	string(MsgUpdateGroupPolicyAdmin),
	string(MsgCreateGroupWithPolicy),
	string(MsgUpdateGroupPolicyDecisionPolicy),
	string(MsgUpdateGroupPolicyMetadata),
	string(MsgSubmitProposalGroup),
	string(MsgWithdrawProposal),
	string(MsgVoteGroup),
	string(MsgExecGroup),
	string(MsgLeaveGroup),
	string(MsgSoftwareUpgrade),
	string(MsgCancelUpgrade),
	string(MsgRegisterInterchainAccount),
	string(MsgSendTx),
	string(MsgRegisterPayee),
	string(MsgRegisterCounterpartyPayee),
	string(MsgPayPacketFee),
	string(MsgPayPacketFeeAsync),
	string(MsgTransfer),
	string(MsgCreateClient),
	string(MsgUpdateClient),
	string(MsgUpgradeClient),
	string(MsgSubmitMisbehaviour),
	string(MsgConnectionOpenInit),
	string(MsgConnectionOpenTry),
	string(MsgConnectionOpenAck),
	string(MsgConnectionOpenConfirm),
	string(MsgChannelOpenInit),
	string(MsgChannelOpenTry),
	string(MsgChannelOpenAck),
	string(MsgChannelOpenConfirm),
	string(MsgChannelCloseInit),
	string(MsgChannelCloseConfirm),
	string(MsgRecvPacket),
	string(MsgTimeout),
	string(MsgTimeoutOnClose),
	string(MsgAcknowledgement),
}

// MsgTypeNames returns a list of possible string values of MsgType.
func MsgTypeNames() []string {
	tmp := make([]string, len(_MsgTypeNames))
	copy(tmp, _MsgTypeNames)
	return tmp
}

// MsgTypeValues returns a list of the values for MsgType
func MsgTypeValues() []MsgType {
	return []MsgType{
		MsgUnknown,
		MsgSetWithdrawAddress,
		MsgWithdrawDelegatorReward,
		MsgWithdrawValidatorCommission,
		MsgFundCommunityPool,
		MsgCreateValidator,
		MsgEditValidator,
		MsgDelegate,
		MsgBeginRedelegate,
		MsgUndelegate,
		MsgCancelUnbondingDelegation,
		MsgUnjail,
		MsgSend,
		MsgMultiSend,
		MsgCreateVestingAccount,
		MsgCreatePermanentLockedAccount,
		MsgCreatePeriodicVestingAccount,
		MsgPayForBlobs,
		MsgGrant,
		MsgExec,
		MsgRevoke,
		MsgGrantAllowance,
		MsgRevokeAllowance,
		MsgRegisterEVMAddress,
		MsgSubmitProposal,
		MsgExecLegacyContent,
		MsgVote,
		MsgVoteWeighted,
		MsgDeposit,
		IBCTransfer,
		MsgVerifyInvariant,
		MsgSubmitEvidence,
		MsgSendNFT,
		MsgCreateGroup,
		MsgUpdateGroupMembers,
		MsgUpdateGroupAdmin,
		MsgUpdateGroupMetadata,
		MsgCreateGroupPolicy,
		MsgUpdateGroupPolicyAdmin,
		MsgCreateGroupWithPolicy,
		MsgUpdateGroupPolicyDecisionPolicy,
		MsgUpdateGroupPolicyMetadata,
		MsgSubmitProposalGroup,
		MsgWithdrawProposal,
		MsgVoteGroup,
		MsgExecGroup,
		MsgLeaveGroup,
		MsgSoftwareUpgrade,
		MsgCancelUpgrade,
		MsgRegisterInterchainAccount,
		MsgSendTx,
		MsgRegisterPayee,
		MsgRegisterCounterpartyPayee,
		MsgPayPacketFee,
		MsgPayPacketFeeAsync,
		MsgTransfer,
		MsgCreateClient,
		MsgUpdateClient,
		MsgUpgradeClient,
		MsgSubmitMisbehaviour,
		MsgConnectionOpenInit,
		MsgConnectionOpenTry,
		MsgConnectionOpenAck,
		MsgConnectionOpenConfirm,
		MsgChannelOpenInit,
		MsgChannelOpenTry,
		MsgChannelOpenAck,
		MsgChannelOpenConfirm,
		MsgChannelCloseInit,
		MsgChannelCloseConfirm,
		MsgRecvPacket,
		MsgTimeout,
		MsgTimeoutOnClose,
		MsgAcknowledgement,
	}
}

// String implements the Stringer interface.
func (x MsgType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x MsgType) IsValid() bool {
	_, err := ParseMsgType(string(x))
	return err == nil
}

var _MsgTypeValue = map[string]MsgType{
	"MsgUnknown":                         MsgUnknown,
	"MsgSetWithdrawAddress":              MsgSetWithdrawAddress,
	"MsgWithdrawDelegatorReward":         MsgWithdrawDelegatorReward,
	"MsgWithdrawValidatorCommission":     MsgWithdrawValidatorCommission,
	"MsgFundCommunityPool":               MsgFundCommunityPool,
	"MsgCreateValidator":                 MsgCreateValidator,
	"MsgEditValidator":                   MsgEditValidator,
	"MsgDelegate":                        MsgDelegate,
	"MsgBeginRedelegate":                 MsgBeginRedelegate,
	"MsgUndelegate":                      MsgUndelegate,
	"MsgCancelUnbondingDelegation":       MsgCancelUnbondingDelegation,
	"MsgUnjail":                          MsgUnjail,
	"MsgSend":                            MsgSend,
	"MsgMultiSend":                       MsgMultiSend,
	"MsgCreateVestingAccount":            MsgCreateVestingAccount,
	"MsgCreatePermanentLockedAccount":    MsgCreatePermanentLockedAccount,
	"MsgCreatePeriodicVestingAccount":    MsgCreatePeriodicVestingAccount,
	"MsgPayForBlobs":                     MsgPayForBlobs,
	"MsgGrant":                           MsgGrant,
	"MsgExec":                            MsgExec,
	"MsgRevoke":                          MsgRevoke,
	"MsgGrantAllowance":                  MsgGrantAllowance,
	"MsgRevokeAllowance":                 MsgRevokeAllowance,
	"MsgRegisterEVMAddress":              MsgRegisterEVMAddress,
	"MsgSubmitProposal":                  MsgSubmitProposal,
	"MsgExecLegacyContent":               MsgExecLegacyContent,
	"MsgVote":                            MsgVote,
	"MsgVoteWeighted":                    MsgVoteWeighted,
	"MsgDeposit":                         MsgDeposit,
	"IBCTransfer":                        IBCTransfer,
	"MsgVerifyInvariant":                 MsgVerifyInvariant,
	"MsgSubmitEvidence":                  MsgSubmitEvidence,
	"MsgSendNFT":                         MsgSendNFT,
	"MsgCreateGroup":                     MsgCreateGroup,
	"MsgUpdateGroupMembers":              MsgUpdateGroupMembers,
	"MsgUpdateGroupAdmin":                MsgUpdateGroupAdmin,
	"MsgUpdateGroupMetadata":             MsgUpdateGroupMetadata,
	"MsgCreateGroupPolicy":               MsgCreateGroupPolicy,
	"MsgUpdateGroupPolicyAdmin":          MsgUpdateGroupPolicyAdmin,
	"MsgCreateGroupWithPolicy":           MsgCreateGroupWithPolicy,
	"MsgUpdateGroupPolicyDecisionPolicy": MsgUpdateGroupPolicyDecisionPolicy,
	"MsgUpdateGroupPolicyMetadata":       MsgUpdateGroupPolicyMetadata,
	"MsgSubmitProposalGroup":             MsgSubmitProposalGroup,
	"MsgWithdrawProposal":                MsgWithdrawProposal,
	"MsgVoteGroup":                       MsgVoteGroup,
	"MsgExecGroup":                       MsgExecGroup,
	"MsgLeaveGroup":                      MsgLeaveGroup,
	"MsgSoftwareUpgrade":                 MsgSoftwareUpgrade,
	"MsgCancelUpgrade":                   MsgCancelUpgrade,
	"MsgRegisterInterchainAccount":       MsgRegisterInterchainAccount,
	"MsgSendTx":                          MsgSendTx,
	"MsgRegisterPayee":                   MsgRegisterPayee,
	"MsgRegisterCounterpartyPayee":       MsgRegisterCounterpartyPayee,
	"MsgPayPacketFee":                    MsgPayPacketFee,
	"MsgPayPacketFeeAsync":               MsgPayPacketFeeAsync,
	"MsgTransfer":                        MsgTransfer,
	"MsgCreateClient":                    MsgCreateClient,
	"MsgUpdateClient":                    MsgUpdateClient,
	"MsgUpgradeClient":                   MsgUpgradeClient,
	"MsgSubmitMisbehaviour":              MsgSubmitMisbehaviour,
	"MsgConnectionOpenInit":              MsgConnectionOpenInit,
	"MsgConnectionOpenTry":               MsgConnectionOpenTry,
	"MsgConnectionOpenAck":               MsgConnectionOpenAck,
	"MsgConnectionOpenConfirm":           MsgConnectionOpenConfirm,
	"MsgChannelOpenInit":                 MsgChannelOpenInit,
	"MsgChannelOpenTry":                  MsgChannelOpenTry,
	"MsgChannelOpenAck":                  MsgChannelOpenAck,
	"MsgChannelOpenConfirm":              MsgChannelOpenConfirm,
	"MsgChannelCloseInit":                MsgChannelCloseInit,
	"MsgChannelCloseConfirm":             MsgChannelCloseConfirm,
	"MsgRecvPacket":                      MsgRecvPacket,
	"MsgTimeout":                         MsgTimeout,
	"MsgTimeoutOnClose":                  MsgTimeoutOnClose,
	"MsgAcknowledgement":                 MsgAcknowledgement,
}

// ParseMsgType attempts to convert a string to a MsgType.
func ParseMsgType(name string) (MsgType, error) {
	if x, ok := _MsgTypeValue[name]; ok {
		return x, nil
	}
	return MsgType(""), fmt.Errorf("%s is %w", name, ErrInvalidMsgType)
}

// MarshalText implements the text marshaller method.
func (x MsgType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *MsgType) UnmarshalText(text []byte) error {
	tmp, err := ParseMsgType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errMsgTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *MsgType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = MsgType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseMsgType(v)
	case []byte:
		*x, err = ParseMsgType(string(v))
	case MsgType:
		*x = v
	case *MsgType:
		if v == nil {
			return errMsgTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errMsgTypeNilPtr
		}
		*x, err = ParseMsgType(*v)
	default:
		return errors.New("invalid type for MsgType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x MsgType) Value() (driver.Value, error) {
	return x.String(), nil
}

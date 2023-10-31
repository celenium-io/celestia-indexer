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
)

const (
	// EventTypeUnknown is a EventType of type unknown.
	EventTypeUnknown EventType = "unknown"
	// EventTypeCoinReceived is a EventType of type coin_received.
	EventTypeCoinReceived EventType = "coin_received"
	// EventTypeCoinbase is a EventType of type coinbase.
	EventTypeCoinbase EventType = "coinbase"
	// EventTypeCoinSpent is a EventType of type coin_spent.
	EventTypeCoinSpent EventType = "coin_spent"
	// EventTypeBurn is a EventType of type burn.
	EventTypeBurn EventType = "burn"
	// EventTypeMint is a EventType of type mint.
	EventTypeMint EventType = "mint"
	// EventTypeMessage is a EventType of type message.
	EventTypeMessage EventType = "message"
	// EventTypeProposerReward is a EventType of type proposer_reward.
	EventTypeProposerReward EventType = "proposer_reward"
	// EventTypeRewards is a EventType of type rewards.
	EventTypeRewards EventType = "rewards"
	// EventTypeCommission is a EventType of type commission.
	EventTypeCommission EventType = "commission"
	// EventTypeLiveness is a EventType of type liveness.
	EventTypeLiveness EventType = "liveness"
	// EventTypeTransfer is a EventType of type transfer.
	EventTypeTransfer EventType = "transfer"
	// EventTypeCelestiablobv1EventPayForBlobs is a EventType of type celestia.blob.v1.EventPayForBlobs.
	EventTypeCelestiablobv1EventPayForBlobs EventType = "celestia.blob.v1.EventPayForBlobs"
	// EventTypeRedelegate is a EventType of type redelegate.
	EventTypeRedelegate EventType = "redelegate"
	// EventTypeAttestationRequest is a EventType of type AttestationRequest.
	EventTypeAttestationRequest EventType = "AttestationRequest"
	// EventTypeWithdrawRewards is a EventType of type withdraw_rewards.
	EventTypeWithdrawRewards EventType = "withdraw_rewards"
	// EventTypeWithdrawCommission is a EventType of type withdraw_commission.
	EventTypeWithdrawCommission EventType = "withdraw_commission"
	// EventTypeSetWithdrawAddress is a EventType of type set_withdraw_address.
	EventTypeSetWithdrawAddress EventType = "set_withdraw_address"
	// EventTypeCreateValidator is a EventType of type create_validator.
	EventTypeCreateValidator EventType = "create_validator"
	// EventTypeDelegate is a EventType of type delegate.
	EventTypeDelegate EventType = "delegate"
	// EventTypeEditValidator is a EventType of type edit_validator.
	EventTypeEditValidator EventType = "edit_validator"
	// EventTypeUnbond is a EventType of type unbond.
	EventTypeUnbond EventType = "unbond"
	// EventTypeTx is a EventType of type tx.
	EventTypeTx EventType = "tx"
	// EventTypeUseFeegrant is a EventType of type use_feegrant.
	EventTypeUseFeegrant EventType = "use_feegrant"
	// EventTypeRevokeFeegrant is a EventType of type revoke_feegrant.
	EventTypeRevokeFeegrant EventType = "revoke_feegrant"
	// EventTypeSetFeegrant is a EventType of type set_feegrant.
	EventTypeSetFeegrant EventType = "set_feegrant"
	// EventTypeUpdateFeegrant is a EventType of type update_feegrant.
	EventTypeUpdateFeegrant EventType = "update_feegrant"
	// EventTypeSlash is a EventType of type slash.
	EventTypeSlash EventType = "slash"
	// EventTypeProposalVote is a EventType of type proposal_vote.
	EventTypeProposalVote EventType = "proposal_vote"
	// EventTypeProposalDeposit is a EventType of type proposal_deposit.
	EventTypeProposalDeposit EventType = "proposal_deposit"
	// EventTypeSubmitProposal is a EventType of type submit_proposal.
	EventTypeSubmitProposal EventType = "submit_proposal"
	// EventTypeCosmosauthzv1beta1EventGrant is a EventType of type cosmos.authz.v1beta1.EventGrant.
	EventTypeCosmosauthzv1beta1EventGrant EventType = "cosmos.authz.v1beta1.EventGrant"
	// EventTypeSendPacket is a EventType of type send_packet.
	EventTypeSendPacket EventType = "send_packet"
	// EventTypeIbcTransfer is a EventType of type ibc_transfer.
	EventTypeIbcTransfer EventType = "ibc_transfer"
	// EventTypeFungibleTokenPacket is a EventType of type fungible_token_packet.
	EventTypeFungibleTokenPacket EventType = "fungible_token_packet"
	// EventTypeAcknowledgePacket is a EventType of type acknowledge_packet.
	EventTypeAcknowledgePacket EventType = "acknowledge_packet"
	// EventTypeCreateClient is a EventType of type create_client.
	EventTypeCreateClient EventType = "create_client"
	// EventTypeUpdateClient is a EventType of type update_client.
	EventTypeUpdateClient EventType = "update_client"
	// EventTypeConnectionOpenTry is a EventType of type connection_open_try.
	EventTypeConnectionOpenTry EventType = "connection_open_try"
	// EventTypeConnectionOpenInit is a EventType of type connection_open_init.
	EventTypeConnectionOpenInit EventType = "connection_open_init"
	// EventTypeConnectionOpenConfirm is a EventType of type connection_open_confirm.
	EventTypeConnectionOpenConfirm EventType = "connection_open_confirm"
	// EventTypeChannelOpenTry is a EventType of type channel_open_try.
	EventTypeChannelOpenTry EventType = "channel_open_try"
	// EventTypeChannelOpenInit is a EventType of type channel_open_init.
	EventTypeChannelOpenInit EventType = "channel_open_init"
	// EventTypeChannelOpenConfirm is a EventType of type channel_open_confirm.
	EventTypeChannelOpenConfirm EventType = "channel_open_confirm"
	// EventTypeChannelOpenAck is a EventType of type channel_open_ack.
	EventTypeChannelOpenAck EventType = "channel_open_ack"
)

var ErrInvalidEventType = errors.New("not a valid EventType")

// EventTypeValues returns a list of the values for EventType
func EventTypeValues() []EventType {
	return []EventType{
		EventTypeUnknown,
		EventTypeCoinReceived,
		EventTypeCoinbase,
		EventTypeCoinSpent,
		EventTypeBurn,
		EventTypeMint,
		EventTypeMessage,
		EventTypeProposerReward,
		EventTypeRewards,
		EventTypeCommission,
		EventTypeLiveness,
		EventTypeTransfer,
		EventTypeCelestiablobv1EventPayForBlobs,
		EventTypeRedelegate,
		EventTypeAttestationRequest,
		EventTypeWithdrawRewards,
		EventTypeWithdrawCommission,
		EventTypeSetWithdrawAddress,
		EventTypeCreateValidator,
		EventTypeDelegate,
		EventTypeEditValidator,
		EventTypeUnbond,
		EventTypeTx,
		EventTypeUseFeegrant,
		EventTypeRevokeFeegrant,
		EventTypeSetFeegrant,
		EventTypeUpdateFeegrant,
		EventTypeSlash,
		EventTypeProposalVote,
		EventTypeProposalDeposit,
		EventTypeSubmitProposal,
		EventTypeCosmosauthzv1beta1EventGrant,
		EventTypeSendPacket,
		EventTypeIbcTransfer,
		EventTypeFungibleTokenPacket,
		EventTypeAcknowledgePacket,
		EventTypeCreateClient,
		EventTypeUpdateClient,
		EventTypeConnectionOpenTry,
		EventTypeConnectionOpenInit,
		EventTypeConnectionOpenConfirm,
		EventTypeChannelOpenTry,
		EventTypeChannelOpenInit,
		EventTypeChannelOpenConfirm,
		EventTypeChannelOpenAck,
	}
}

// String implements the Stringer interface.
func (x EventType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x EventType) IsValid() bool {
	_, err := ParseEventType(string(x))
	return err == nil
}

var _EventTypeValue = map[string]EventType{
	"unknown":                           EventTypeUnknown,
	"coin_received":                     EventTypeCoinReceived,
	"coinbase":                          EventTypeCoinbase,
	"coin_spent":                        EventTypeCoinSpent,
	"burn":                              EventTypeBurn,
	"mint":                              EventTypeMint,
	"message":                           EventTypeMessage,
	"proposer_reward":                   EventTypeProposerReward,
	"rewards":                           EventTypeRewards,
	"commission":                        EventTypeCommission,
	"liveness":                          EventTypeLiveness,
	"transfer":                          EventTypeTransfer,
	"celestia.blob.v1.EventPayForBlobs": EventTypeCelestiablobv1EventPayForBlobs,
	"redelegate":                        EventTypeRedelegate,
	"AttestationRequest":                EventTypeAttestationRequest,
	"withdraw_rewards":                  EventTypeWithdrawRewards,
	"withdraw_commission":               EventTypeWithdrawCommission,
	"set_withdraw_address":              EventTypeSetWithdrawAddress,
	"create_validator":                  EventTypeCreateValidator,
	"delegate":                          EventTypeDelegate,
	"edit_validator":                    EventTypeEditValidator,
	"unbond":                            EventTypeUnbond,
	"tx":                                EventTypeTx,
	"use_feegrant":                      EventTypeUseFeegrant,
	"revoke_feegrant":                   EventTypeRevokeFeegrant,
	"set_feegrant":                      EventTypeSetFeegrant,
	"update_feegrant":                   EventTypeUpdateFeegrant,
	"slash":                             EventTypeSlash,
	"proposal_vote":                     EventTypeProposalVote,
	"proposal_deposit":                  EventTypeProposalDeposit,
	"submit_proposal":                   EventTypeSubmitProposal,
	"cosmos.authz.v1beta1.EventGrant":   EventTypeCosmosauthzv1beta1EventGrant,
	"send_packet":                       EventTypeSendPacket,
	"ibc_transfer":                      EventTypeIbcTransfer,
	"fungible_token_packet":             EventTypeFungibleTokenPacket,
	"acknowledge_packet":                EventTypeAcknowledgePacket,
	"create_client":                     EventTypeCreateClient,
	"update_client":                     EventTypeUpdateClient,
	"connection_open_try":               EventTypeConnectionOpenTry,
	"connection_open_init":              EventTypeConnectionOpenInit,
	"connection_open_confirm":           EventTypeConnectionOpenConfirm,
	"channel_open_try":                  EventTypeChannelOpenTry,
	"channel_open_init":                 EventTypeChannelOpenInit,
	"channel_open_confirm":              EventTypeChannelOpenConfirm,
	"channel_open_ack":                  EventTypeChannelOpenAck,
}

// ParseEventType attempts to convert a string to a EventType.
func ParseEventType(name string) (EventType, error) {
	if x, ok := _EventTypeValue[name]; ok {
		return x, nil
	}
	return EventType(""), fmt.Errorf("%s is %w", name, ErrInvalidEventType)
}

// MarshalText implements the text marshaller method.
func (x EventType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *EventType) UnmarshalText(text []byte) error {
	tmp, err := ParseEventType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errEventTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *EventType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = EventType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseEventType(v)
	case []byte:
		*x, err = ParseEventType(string(v))
	case EventType:
		*x = v
	case *EventType:
		if v == nil {
			return errEventTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errEventTypeNilPtr
		}
		*x, err = ParseEventType(*v)
	default:
		return errors.New("invalid type for EventType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x EventType) Value() (driver.Value, error) {
	return x.String(), nil
}

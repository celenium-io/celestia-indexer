// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
)

type EventHandler func(ctx *context.Context, c *Cursor, msg *storage.Message) error

var eventHandlers = map[storageTypes.MsgType]EventHandler{
	storageTypes.MsgDelegate:                       handleDelegate,
	storageTypes.MsgBeginRedelegate:                handleRedelegate,
	storageTypes.MsgUndelegate:                     handleUndelegate,
	storageTypes.MsgCancelUnbondingDelegation:      handleCancelUnbonding,
	storageTypes.MsgExec:                           handleExec,
	storageTypes.MsgWithdrawValidatorCommission:    handleWithdrawValidatorCommission,
	storageTypes.MsgWithdrawDelegatorReward:        handleWithdrawDelegatorRewards,
	storageTypes.MsgUnjail:                         handleUnjail,
	storageTypes.MsgSubmitProposal:                 handleSubmitProposal,
	storageTypes.MsgDeposit:                        handleDeposit,
	storageTypes.MsgVote:                           handleVote,
	storageTypes.MsgVoteWeighted:                   handleVote,
	storageTypes.MsgCreateClient:                   handleCreateClient,
	storageTypes.MsgUpdateClient:                   handleUpdateClient,
	storageTypes.MsgConnectionOpenInit:             handleConnectionOpenInit,
	storageTypes.MsgConnectionOpenTry:              handleConnectionOpenInit,
	storageTypes.MsgConnectionOpenConfirm:          handleConnectionOpenConfirm,
	storageTypes.MsgConnectionOpenAck:              handleConnectionOpenConfirm,
	storageTypes.MsgChannelOpenInit:                handleChannelOpenInit,
	storageTypes.MsgChannelOpenTry:                 handleChannelOpenInit,
	storageTypes.MsgChannelOpenConfirm:             handleChannelOpenConfirm,
	storageTypes.MsgChannelOpenAck:                 handleChannelOpenConfirm,
	storageTypes.MsgChannelCloseInit:               handleChannelClose,
	storageTypes.MsgChannelCloseConfirm:            handleChannelClose,
	storageTypes.MsgAcknowledgement:                handleAcknowledgement,
	storageTypes.MsgRecvPacket:                     handleRecvPacket,
	storageTypes.MsgCreateMailbox:                  handleCreateMailbox,
	storageTypes.MsgSetMailbox:                     handleSetMailbox,
	storageTypes.MsgProcessMessage:                 handleHyperlaneProcessMessage,
	storageTypes.MsgRemoteTransfer:                 handleHyperlaneRemoteTransfer,
	storageTypes.MsgCreateCollateralToken:          handleCreateCollateralToken,
	storageTypes.MsgCreateSyntheticToken:           handleCreateSyntheticToken,
	storageTypes.MsgSetToken:                       handleSetToken,
	storageTypes.MsgSend:                           handleSend,
	storageTypes.MsgForward:                        handleForward,
	storageTypes.MsgCreateInterchainSecurityModule: handleCreateZkISM,
	storageTypes.MsgUpdateInterchainSecurityModule: handleUpdateZkISM,
	storageTypes.MsgSubmitMessages:                 handleSubmitZkISMMessages,
}

func handle(ctx *context.Context, c *Cursor, msg *storage.Message, eventHandlers map[storageTypes.MsgType]EventHandler, stopKey string) error {
	if handler, ok := eventHandlers[msg.Type]; ok {
		return handler(ctx, c, msg)
	}

	// if event handler is not found, skip events up to the next message
	// carrying an event of type "message" whose Data has a non-empty
	// stopKey (mirrors MsgEvents' boundary check, but additionally
	// requires EventTypeMessage — a plain stopKey match on some other
	// event type does not end the scan here).
	c.Next()
	for {
		event, ok := c.Peek()
		if !ok {
			return nil
		}
		if event.Type == storageTypes.EventTypeMessage && decoder.StringFromMap(event.Data, stopKey) != "" {
			return nil
		}
		c.Next()
	}
}

// Handle dispatches msg to the registered top-level handler for its type
// (or skips its events if none is registered), then advances c to the
// "action" boundary event that starts the next message in the
// transaction, ready for the next call to Handle.
func Handle(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if err := handle(ctx, c, msg, eventHandlers, "action"); err != nil {
		return err
	}
	c.SkipToNext("action")
	return nil
}

var ibcEventHandlers = map[storageTypes.MsgType]EventHandler{
	storageTypes.MsgDelegate:                       processDelegate,
	storageTypes.MsgBeginRedelegate:                processRedelegate,
	storageTypes.MsgUndelegate:                     processUndelegate,
	storageTypes.MsgCancelUnbondingDelegation:      processCancelUnbonding,
	storageTypes.MsgExec:                           processExec,
	storageTypes.MsgWithdrawValidatorCommission:    processWithdrawValidatorCommission,
	storageTypes.MsgWithdrawDelegatorReward:        processWithdrawDelegatorRewards,
	storageTypes.MsgUnjail:                         processUnjail,
	storageTypes.MsgSubmitProposal:                 processSubmitProposal,
	storageTypes.MsgDeposit:                        processDeposit,
	storageTypes.MsgVote:                           processVote,
	storageTypes.MsgVoteWeighted:                   processVote,
	storageTypes.MsgCreateClient:                   processCreateClient,
	storageTypes.MsgUpdateClient:                   processUpdateClient,
	storageTypes.MsgConnectionOpenInit:             processConnectionOpenInit,
	storageTypes.MsgConnectionOpenTry:              processConnectionOpenInit,
	storageTypes.MsgConnectionOpenConfirm:          processConnectionOpenConfirm,
	storageTypes.MsgConnectionOpenAck:              processConnectionOpenConfirm,
	storageTypes.MsgChannelOpenInit:                processChannelOpenInit,
	storageTypes.MsgChannelOpenTry:                 processChannelOpenInit,
	storageTypes.MsgChannelOpenConfirm:             processChannelOpenConfirm,
	storageTypes.MsgChannelOpenAck:                 processChannelOpenConfirm,
	storageTypes.MsgChannelCloseInit:               processChannelClose,
	storageTypes.MsgChannelCloseConfirm:            processChannelClose,
	storageTypes.MsgCreateMailbox:                  processCreateMailbox,
	storageTypes.MsgSetMailbox:                     processSetMailbox,
	storageTypes.MsgProcessMessage:                 processHyperlaneProcessMessage,
	storageTypes.MsgRemoteTransfer:                 processHyperlaneRemoteTransfer,
	storageTypes.MsgCreateCollateralToken:          processCreateCollateralToken,
	storageTypes.MsgCreateSyntheticToken:           processCreateSyntheticToken,
	storageTypes.MsgSetToken:                       processSetToken,
	storageTypes.MsgSend:                           processSend,
	storageTypes.MsgForward:                        processForward,
	storageTypes.MsgCreateInterchainSecurityModule: processCreateZkISM,
	storageTypes.MsgUpdateInterchainSecurityModule: processUpdateZkISM,
	storageTypes.MsgSubmitMessages:                 processSubmitZkISMMessages,
}

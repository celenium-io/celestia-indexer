// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
)

type EventHandler func(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error

var eventHandlers = map[storageTypes.MsgType]EventHandler{
	storageTypes.MsgDelegate:                    handleDelegate,
	storageTypes.MsgBeginRedelegate:             handleRedelegate,
	storageTypes.MsgUndelegate:                  handleUndelegate,
	storageTypes.MsgCancelUnbondingDelegation:   handleCancelUnbonding,
	storageTypes.MsgExec:                        handleExec,
	storageTypes.MsgWithdrawValidatorCommission: handleWithdrawValidatorCommission,
	storageTypes.MsgWithdrawDelegatorReward:     handleWithdrawDelegatorRewards,
	storageTypes.MsgUnjail:                      handleUnjail,
	storageTypes.MsgSubmitProposal:              handleSubmitProposal,
	storageTypes.MsgDeposit:                     handleDeposit,
	storageTypes.MsgVote:                        handleVote,
	storageTypes.MsgVoteWeighted:                handleVote,
	storageTypes.MsgCreateClient:                handleCreateClient,
	storageTypes.MsgUpdateClient:                handleUpdateClient,
	storageTypes.MsgConnectionOpenInit:          handleConnectionOpenInit,
	storageTypes.MsgConnectionOpenTry:           handleConnectionOpenInit,
	storageTypes.MsgConnectionOpenConfirm:       handleConnectionOpenConfirm,
	storageTypes.MsgConnectionOpenAck:           handleConnectionOpenConfirm,
	storageTypes.MsgChannelOpenInit:             handleChannelOpenInit,
	storageTypes.MsgChannelOpenTry:              handleChannelOpenInit,
	storageTypes.MsgChannelOpenConfirm:          handleChannelOpenConfirm,
	storageTypes.MsgChannelOpenAck:              handleChannelOpenConfirm,
	storageTypes.MsgChannelCloseInit:            handleChannelClose,
	storageTypes.MsgChannelCloseConfirm:         handleChannelClose,
	storageTypes.MsgAcknowledgement:             handleAcknowledgement,
	storageTypes.MsgRecvPacket:                  handleRecvPacket,
}

func handle(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int, eventHandlers map[storageTypes.MsgType]EventHandler) error {
	if handler, ok := eventHandlers[msg.Type]; ok {
		return handler(ctx, events, msg, idx)
	}

	// if event handler is not found list events to another action
	*idx++

	startIndex := *idx
	for _, event := range events[startIndex:] {
		if event.Type != storageTypes.EventTypeMessage {
			*idx++
			continue
		}
		if action := decoder.StringFromMap(event.Data, "action"); action == "" {
			*idx++
			continue
		}
		break
	}

	return nil
}

func Handle(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	return handle(ctx, events, msg, idx, eventHandlers)
}

var ibcEventHandlers = map[storageTypes.MsgType]EventHandler{
	storageTypes.MsgDelegate:                    processDelegate,
	storageTypes.MsgBeginRedelegate:             processRedelegate,
	storageTypes.MsgUndelegate:                  processUndelegate,
	storageTypes.MsgCancelUnbondingDelegation:   processCancelUnbonding,
	storageTypes.MsgExec:                        processExec,
	storageTypes.MsgWithdrawValidatorCommission: processWithdrawValidatorCommission,
	storageTypes.MsgWithdrawDelegatorReward:     processWithdrawDelegatorRewards,
	storageTypes.MsgUnjail:                      processUnjail,
	storageTypes.MsgSubmitProposal:              processSubmitProposal,
	storageTypes.MsgDeposit:                     processDeposit,
	storageTypes.MsgVote:                        processVote,
	storageTypes.MsgVoteWeighted:                processVote,
	storageTypes.MsgCreateClient:                processCreateClient,
	storageTypes.MsgUpdateClient:                processUpdateClient,
	storageTypes.MsgConnectionOpenInit:          processConnectionOpenInit,
	storageTypes.MsgConnectionOpenTry:           processConnectionOpenInit,
	storageTypes.MsgConnectionOpenConfirm:       processConnectionOpenConfirm,
	storageTypes.MsgConnectionOpenAck:           processConnectionOpenConfirm,
	storageTypes.MsgChannelOpenInit:             processChannelOpenInit,
	storageTypes.MsgChannelOpenTry:              processChannelOpenInit,
	storageTypes.MsgChannelOpenConfirm:          processChannelOpenConfirm,
	storageTypes.MsgChannelOpenAck:              processChannelOpenConfirm,
	storageTypes.MsgChannelCloseInit:            processChannelClose,
	storageTypes.MsgChannelCloseConfirm:         processChannelClose,
}

func toTheNextAction(events []storage.Event, idx *int) {
	if len(events)-1 <= *idx {
		return
	}
	var action = decoder.StringFromMap(events[*idx].Data, "action")
	for action == "" && len(events)-1 > *idx {
		*idx += 1
		action = decoder.StringFromMap(events[*idx].Data, "action")
	}
}

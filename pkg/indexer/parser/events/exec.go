// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleExec(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.authz.v1beta1.MsgExec" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1

	return processExec(ctx, events, msg, idx)
}

func processExec(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	for i := range msg.InternalMsgs {
		msgs, err := getInternalDataForExec(msg.Data, i)
		if err != nil {
			return err
		}

		internalMessage := &storage.Message{
			Height: msg.Height,
			Time:   msg.Time,
			Data:   msgs,
		}

		switch msg.InternalMsgs[i] {
		case "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation":
			if err := processCancelUnbonding(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgDelegate":
			if err := processDelegate(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			if err := processRedelegate(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			if err := processUndelegate(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case msgWithdrawValidatorCommission:
			if err := processWithdrawValidatorCommission(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			if err := processWithdrawDelegatorRewards(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/cosmos.slashing.v1beta1.MsgUnjail":
			if err := processUnjail(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		case "/celestia.signal.v1.MsgSignalVersion":
			data, err := getInternalDataForExec(msg.Data, i)
			if err != nil {
				return err
			}
			if err := processSignalVersion(ctx, events, msg, data, idx); err != nil {
				return err
			}
		case "/cosmos.gov.v1beta1.MsgVote", "/cosmos.gov.v1.MsgVote", "/cosmos.gov.v1.MsgVoteWeighted", "/cosmos.gov.v1beta1.MsgVoteWeighted":
			if err := processVote(ctx, events, internalMessage, idx); err != nil {
				return err
			}
		default:
			for j := *idx; j < len(events); j++ {
				authMsgIdxPtr, err := decoder.AuthMsgIndexFromMap(events[*idx].Data)
				if err != nil {
					return err
				}
				if authMsgIdxPtr == nil {
					break
				}

				if *authMsgIdxPtr != int64(i) {
					break
				}

				*idx = j + 1
			}

		}
	}

	return nil
}

func getInternalDataForExec(data map[string]any, idx int) (map[string]any, error) {
	msgsAny, ok := data["Msgs"]
	if !ok {
		return nil, errors.Errorf("can't find Msgs key in MsgExec: %##v", data)
	}
	msgsArr, ok := msgsAny.([]any)
	if !ok {
		return nil, errors.Errorf("Msgs is not an array in MsgExec: %T", msgsAny)
	}
	if idx < 0 || idx >= len(msgsArr) {
		return nil, errors.Errorf("Msgs index out of range in MsgExec: %d >= %d", idx, len(msgsArr))
	}
	msgs, ok := msgsArr[idx].(map[string]any)
	if !ok {
		return nil, errors.Errorf("Msgs invalid type in MsgExec: %T", msgsArr[idx])
	}
	return msgs, nil
}

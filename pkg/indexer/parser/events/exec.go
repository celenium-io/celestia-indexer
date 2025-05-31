// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
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
		switch msg.InternalMsgs[i] {
		case "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation":
			if err := processCancelUnbonding(ctx, events, msg, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgDelegate":
			if err := processDelegate(ctx, events, msg, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			msgsAny, ok := msg.Data["Msgs"]
			if !ok {
				return errors.Errorf("can't find Msgs key in MsgExec: %##v", msg.Data)
			}
			msgsArr, ok := msgsAny.([]any)
			if !ok {
				return errors.Errorf("Msgs is not an array in MsgExec: %T", msgsAny)
			}
			msgs, ok := msgsArr[i].(map[string]any)
			if !ok {
				return errors.Errorf("Msgs invalid type in MsgExec: %T", msgsArr[i])
			}
			if err := processRedelegate(ctx, events, &storage.Message{
				Data: msgs,
			}, idx); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			if err := processUndelegate(ctx, events, msg, idx); err != nil {
				return err
			}
		case "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission":
			if err := processWithdrawValidatorCommission(ctx, events, msg, idx); err != nil {
				return err
			}
		case "cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			if err := processWithdrawDelegatorRewards(ctx, events, msg, idx); err != nil {
				return err
			}
		case "/cosmos.slashing.v1beta1.MsgUnjail":
			if err := processUnjail(ctx, events, msg, idx); err != nil {
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

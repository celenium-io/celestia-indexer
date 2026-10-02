// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleExec(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	if c == nil {
		return errors.New("nil event cursor")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	event, _ := c.Peek()
	if action := decoder.StringFromMap(event.Data, "action"); action != "/cosmos.authz.v1beta1.MsgExec" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	c.Next()

	return processExec(ctx, c, msg)
}

func processExec(ctx *context.Context, c *Cursor, msg *storage.Message) error {
	for i := range msg.InternalMsgs {
		msgs, err := getInternalDataForExec(msg.Data, i)
		if err != nil {
			return err
		}

		// handlers read only the events of inner message i and can't overrun into the next ones
		sub, err := c.SubAuthz(i)
		if err != nil {
			return err
		}

		// inner messages have no rows of their own: entities point at the MsgExec
		internalMessage := &storage.Message{
			Id:     msg.Id,
			TxId:   msg.TxId,
			Height: msg.Height,
			Time:   msg.Time,
			Data:   msgs,
		}

		switch msg.InternalMsgs[i] {
		case "/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation":
			if err := processCancelUnbonding(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgDelegate":
			if err := processDelegate(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			if err := processRedelegate(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			if err := processUndelegate(ctx, sub, internalMessage); err != nil {
				return err
			}
		case msgWithdrawValidatorCommission:
			if err := processWithdrawValidatorCommission(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			if err := processWithdrawDelegatorRewards(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/cosmos.slashing.v1beta1.MsgUnjail":
			if err := processExecUnjail(ctx, internalMessage); err != nil {
				return err
			}
		case "/celestia.signal.v1.MsgSignalVersion":
			if err := processSignalVersion(ctx, sub, msg, internalMessage.Data); err != nil {
				return err
			}
		case "/celestia.signal.v1.MsgTryUpgrade":
			if err := processTryUpgrade(ctx, sub, msg, internalMessage.Data); err != nil {
				return err
			}
		case "/cosmos.gov.v1beta1.MsgVote", "/cosmos.gov.v1.MsgVote", "/cosmos.gov.v1.MsgVoteWeighted", "/cosmos.gov.v1beta1.MsgVoteWeighted":
			if err := processVote(ctx, sub, internalMessage); err != nil {
				return err
			}
		case "/celestia.forwarding.v1.MsgForward":
			if err := processForward(ctx, sub, internalMessage); err != nil {
				return err
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

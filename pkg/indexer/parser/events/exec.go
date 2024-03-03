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
		return errors.New("nil message in events hanler")
	}
	if action := decoder.StringFromMap(events[*idx].Data, "action"); action != "/cosmos.authz.v1beta1.MsgExec" {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1

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
			if err := processRedelegate(ctx, events, msg, idx); err != nil {
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
		case "/cosmos.slashing.v1beta1.MsgUnjail":
			if err := processUnjail(ctx, events, idx); err != nil {
				return err
			}
		}
	}

	return nil
}

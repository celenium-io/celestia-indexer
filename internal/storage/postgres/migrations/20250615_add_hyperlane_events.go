// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package migrations

import (
	"context"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(upAddEventHyperlane, downAddEventHyperlane)
}

func upAddEventHyperlane(ctx context.Context, db *bun.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorev1EventDispatch.String(), types.EventTypeUpdateClientProposal.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorev1EventProcess.String(), types.EventTypeHyperlanecorev1EventDispatch.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorev1EventCreateMailbox.String(), types.EventTypeHyperlanecorev1EventProcess.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorev1EventSetMailbox.String(), types.EventTypeHyperlanecorev1EventCreateMailbox.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventCreateSyntheticToken.String(), types.EventTypeHyperlanecorev1EventSetMailbox.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventCreateCollateralToken.String(), types.EventTypeHyperlanewarpv1EventCreateSyntheticToken.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventSetToken.String(), types.EventTypeHyperlanewarpv1EventCreateCollateralToken.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventEnrollRemoteRouter.String(), types.EventTypeHyperlanewarpv1EventSetToken.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventUnrollRemoteRouter.String(), types.EventTypeHyperlanewarpv1EventEnrollRemoteRouter.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventSendRemoteTransfer.String(), types.EventTypeHyperlanewarpv1EventUnrollRemoteRouter.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanewarpv1EventReceiveRemoteTransfer.String(), types.EventTypeHyperlanewarpv1EventSendRemoteTransfer.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventCreateMerkleTreeHook.String(), types.EventTypeHyperlanewarpv1EventReceiveRemoteTransfer.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventInsertedIntoTree.String(), types.EventTypeHyperlanecorepostDispatchv1EventCreateMerkleTreeHook.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventGasPayment.String(), types.EventTypeHyperlanecorepostDispatchv1EventInsertedIntoTree.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventCreateNoopHook.String(), types.EventTypeHyperlanecorepostDispatchv1EventGasPayment.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventCreateIgp.String(), types.EventTypeHyperlanecorepostDispatchv1EventCreateNoopHook.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventSetIgp.String(), types.EventTypeHyperlanecorepostDispatchv1EventCreateIgp.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventSetDestinationGasConfig.String(), types.EventTypeHyperlanecorepostDispatchv1EventSetIgp.String()); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TYPE event_type ADD VALUE IF NOT EXISTS ? AFTER ?`, types.EventTypeHyperlanecorepostDispatchv1EventClaimIgp.String(), types.EventTypeHyperlanecorepostDispatchv1EventSetDestinationGasConfig.String()); err != nil {
		return err
	}
	return nil
}

func downAddEventHyperlane(ctx context.Context, db *bun.DB) error {
	return nil
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	evidenceTypes "cosmossdk.io/x/evidence/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSubmitEvidence represents a message that supports submitting arbitrary
// Evidence of misbehavior such as equivocation or counterfactual signing.
func MsgSubmitEvidence(ctx *context.Context, msgId uint64, m *evidenceTypes.MsgSubmitEvidence) (storageTypes.MsgType, error) {
	msgType := storageTypes.MsgSubmitEvidence
	err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSubmitter, address: m.Submitter},
	}, ctx.Block.Height, msgId)
	return msgType, err
}

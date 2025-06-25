// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	evidenceTypes "cosmossdk.io/x/evidence/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
)

// MsgSubmitEvidence represents a message that supports submitting arbitrary
// Evidence of misbehavior such as equivocation or counterfactual signing.
func MsgSubmitEvidence(ctx *context.Context, m *evidenceTypes.MsgSubmitEvidence) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitEvidence
	addresses, err := createAddresses(ctx, addressesData{
		{t: storageTypes.MsgAddressTypeSubmitter, address: m.Submitter},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

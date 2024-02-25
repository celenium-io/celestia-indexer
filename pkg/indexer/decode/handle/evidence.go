// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handle

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
)

// MsgSubmitEvidence represents a message that supports submitting arbitrary
// Evidence of misbehavior such as equivocation or counterfactual signing.
func MsgSubmitEvidence(ctx *context.Context, m *evidenceTypes.MsgSubmitEvidence) (storageTypes.MsgType, []storage.AddressWithType, error) {
	msgType := storageTypes.MsgSubmitEvidence
	addresses, err := createAddresses(addressesData{
		{t: storageTypes.MsgAddressTypeSubmitter, address: m.Submitter},
	}, ctx.Block.Height)
	return msgType, addresses, err
}

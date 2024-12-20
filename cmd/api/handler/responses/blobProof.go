// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"github.com/celestiaorg/celestia-app/v3/pkg/proof"
)

type BlobProof struct {
	Start int32    `example:"0"  format:"integer" json:"start" swaggertype:"integer"`
	End   int32    `example:"16" format:"integer" json:"end"   swaggertype:"integer"`
	Nodes [][]byte `json:"nodes"`
}

func NewProofs(proofs []*proof.NMTProof) []BlobProof {
	result := make([]BlobProof, len(proofs))
	for i := range proofs {
		result[i] = BlobProof{
			Start: proofs[i].Start,
			End:   proofs[i].End,
			Nodes: proofs[i].Nodes,
		}
	}
	return result
}
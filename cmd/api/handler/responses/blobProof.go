// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"github.com/celestiaorg/celestia-app/v3/pkg/proof"
)

type BlobProof struct {
	Start int32    `example:"0"  format:"integer" json:"start" swaggertype:"integer"`
	End   int32    `example:"16" format:"integer" json:"end"   swaggertype:"integer"`
	Nodes []string `json:"nodes"`
}

func NewProofs(proofs []*proof.NMTProof) []BlobProof {
	result := make([]BlobProof, len(proofs))

	for i := range proofs {
		proofNodes := make([]string, len(proofs[i].Nodes))
		for j, node := range proofs[i].Nodes {
			proofNodes[j] = base64.StdEncoding.EncodeToString(node)
		}

		result[i] = BlobProof{
			Start: proofs[i].Start,
			End:   proofs[i].End,
			Nodes: proofNodes,
		}
	}
	return result
}

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

// BlobProof -
type BlobProof struct {
	ShareProof ShareProofData `json:"share_proof"`
}

// ShareProofData -
type ShareProofData struct {
	ShareProofs []ShareProof `json:"share_proofs"`
}

// ShareProof -
type ShareProof struct {
	Start int      `json:"start"`
	End   int      `json:"end"`
	Nodes []string `json:"nodes"`
}

// RowProof -
type RowProof struct {
	Proofs []Proof `json:"proofs"`
}

// Proof -
type Proof struct {
	LeafHash string `json:"leaf_hash"`
}

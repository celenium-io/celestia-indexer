// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

type SearchItem struct {
	// Result type which is in the result. Can be 'block', 'address', 'namespace', 'tx', 'validator', 'rollup'
	Type string `json:"type"`

	// Search result. Can be one of folowwing types: Block, Address, Namespace, Tx, Validator, Rollup
	Result any `json:"result" swaggertype:"object"`
}

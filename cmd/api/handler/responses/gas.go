// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

type GasPrice struct {
	Slow   string `example:"0.1234" format:"string" json:"slow"   swaggertype:"string"`
	Median string `example:"0.1234" format:"string" json:"median" swaggertype:"string"`
	Fast   string `example:"0.1234" format:"string" json:"fast"   swaggertype:"string"`

	ComputedBlocks []GasBlock `json:"computed_blocks"`
}

type GasBlock struct {
	Height       uint64 `example:"12345"      format:"int64"  json:"height"           swaggertype:"integer"`
	GasWanted    uint64 `example:"86756"      format:"int64"  json:"total_gas_wanted" swaggertype:"integer"`
	GasUsed      uint64 `example:"56789"      format:"int64"  json:"total_gas_used"   swaggertype:"integer"`
	TxCount      uint64 `example:"12"         format:"int64"  json:"tx_count"         swaggertype:"integer"`
	Fee          string `example:"1972367126" format:"string" json:"total_fee"        swaggertype:"string"`
	GasPrice     string `example:"0.12345"    format:"string" json:"avg_gas_price"    swaggertype:"string"`
	GasUsedRatio string `example:"0.12345"    format:"string" json:"gas_used_ratio"   swaggertype:"string"`

	Percentiles []string `json:"percentiles"`
}

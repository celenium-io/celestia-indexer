package types

// swagger:enum TxAddressType
/*
	ENUM(
		validatorAddress,
		delegatorAddress,
		validatorSrcAddress,
		validatorDstAddress,
		fromAddress,
		toAddress
	)
*/
//go:generate go-enum --marshal --sql --values
type TxAddressType string

package types

// swagger:enum TxAddressType
/*
	ENUM(
		validatorAddress,
		delegatorAddress,
		validatorSrcAddress,
		validatorDstAddress,
		fromAddress,
		toAddress,
		grantee,
		granter,
		signer,
		withdraw
	)
*/
//go:generate go-enum --marshal --sql --values
type TxAddressType string

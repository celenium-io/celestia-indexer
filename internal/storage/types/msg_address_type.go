package types

// swagger:enum MsgAddressType
/*
	ENUM(
		validator,
		delegator,
		depositor,

		validatorSrc,
		validatorDst,

		fromAddress,
		toAddress,
		input,
		output,

		grantee,
		granter,
		signer,
		withdraw,

		voter,
		proposer,
	)
*/
//go:generate go-enum --marshal --sql --values
type MsgAddressType string

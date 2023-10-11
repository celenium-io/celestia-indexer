// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

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
		authority,
	)
*/
//go:generate go-enum --marshal --sql --values
type MsgAddressType string

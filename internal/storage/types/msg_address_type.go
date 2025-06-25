// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
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

		sender,
		receiver,

		submitter,

		admin,
		newAdmin,
		groupPolicyAddress,
		executor,
		groupMember,

		owner,

		relayer,
		payee,
	)
*/
//go:generate go-enum --marshal --sql --values
type MsgAddressType string

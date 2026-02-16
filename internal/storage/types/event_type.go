// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

// swagger:enum EventType
/*
	ENUM(
		unknown,
		coin_received,
		coinbase,
		coin_spent,
		burn,
		mint,
		message,
		proposer_reward,
		rewards,
		commission,
		liveness,
		transfer,

		celestia.blob.v1.EventPayForBlobs,

		redelegate,
		AttestationRequest,
		withdraw_rewards,
		withdraw_commission,
		set_withdraw_address,
		create_validator,
		delegate,
		edit_validator,
		unbond,
		tx,
		complete_redelegation,
		complete_unbonding

		use_feegrant,
		revoke_feegrant,
		set_feegrant,
		update_feegrant,
		slash,

		proposal_vote,
		proposal_deposit,
		submit_proposal,

		cosmos.authz.v1beta1.EventGrant,

		send_packet,
		ibc_transfer,

		fungible_token_packet,
		acknowledge_packet,

		create_client,
		update_client,

		connection_open_try,
		connection_open_init,
		connection_open_confirm,
		connection_open_ack,

		channel_open_try,
		channel_open_init,
		channel_open_confirm,
		channel_open_ack,

		recv_packet,
		write_acknowledgement,

		timeout,
		timeout_packet,

		cosmos.authz.v1beta1.EventRevoke,
		cosmos.authz.v1.EventRevoke,
		cancel_unbonding_delegation,

		active_proposal,
		inactive_proposal,
		ics27_packet,
		channel_close_confirm,

		update_client_proposal,

		hyperlane.core.v1.EventDispatch,
		hyperlane.core.v1.EventProcess,
		hyperlane.core.v1.EventCreateMailbox,
		hyperlane.core.v1.EventSetMailbox,
		hyperlane.warp.v1.EventCreateSyntheticToken,
		hyperlane.warp.v1.EventCreateCollateralToken,
		hyperlane.warp.v1.EventSetToken,
		hyperlane.warp.v1.EventEnrollRemoteRouter,
		hyperlane.warp.v1.EventUnrollRemoteRouter,
		hyperlane.warp.v1.EventSendRemoteTransfer,
		hyperlane.warp.v1.EventReceiveRemoteTransfer,
		hyperlane.core.post_dispatch.v1.EventCreateMerkleTreeHook,
		hyperlane.core.post_dispatch.v1.EventInsertedIntoTree,
		hyperlane.core.post_dispatch.v1.EventGasPayment,
		hyperlane.core.post_dispatch.v1.EventCreateNoopHook,
		hyperlane.core.post_dispatch.v1.EventCreateIgp,
		hyperlane.core.post_dispatch.v1.EventSetIgp,
		hyperlane.core.post_dispatch.v1.EventSetDestinationGasConfig,
		hyperlane.core.post_dispatch.v1.EventClaimIgp,
		hyperlane.core.interchain_security.v1.EventCreateNoopIsm,
		hyperlane.core.interchain_security.v1.EventSetRoutingIsmDomain,
		hyperlane.core.interchain_security.v1.EventSetRoutingIsm,
		hyperlane.core.interchain_security.v1.EventCreateRoutingIsm,

		signal_version,
		ibccallbackerror-ics27_packet,

		celestia.forwarding.v1.EventTokenForwarded,
		celestia.forwarding.v1.EventForwardingComplete,

		celestia.zkism.v1.EventCreateInterchainSecurityModule,
		celestia.zkism.v1.EventUpdateInterchainSecurityModule,
		celestia.zkism.v1.EventSubmitMessages
	)
*/
//go:generate go-enum --marshal --sql --values --names
type EventType string

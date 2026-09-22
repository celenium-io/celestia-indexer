// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package types

import (
	"slices"
	"testing"

	icaTypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transferTypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clientTypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectionTypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channelTypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	coreTypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	"github.com/stretchr/testify/require"
)

// Every event of the IBC modules enabled in celestia-app (core, transfer, ICA host)
// must be a known EventType, otherwise it is silently stored as 'unknown'.
func TestEventType_CoversIbcGoEvents(t *testing.T) {
	ibcEvents := []string{
		clientTypes.EventTypeCreateClient,
		clientTypes.EventTypeUpdateClient,
		clientTypes.EventTypeUpgradeClient,
		clientTypes.EventTypeSubmitMisbehaviour,
		clientTypes.EventTypeRecoverClient,
		clientTypes.EventTypeScheduleIBCSoftwareUpgrade,
		clientTypes.EventTypeUpgradeChain,

		connectionTypes.EventTypeConnectionOpenInit,
		connectionTypes.EventTypeConnectionOpenTry,
		connectionTypes.EventTypeConnectionOpenAck,
		connectionTypes.EventTypeConnectionOpenConfirm,

		channelTypes.EventTypeSendPacket,
		channelTypes.EventTypeRecvPacket,
		channelTypes.EventTypeWriteAck,
		channelTypes.EventTypeAcknowledgePacket,
		channelTypes.EventTypeTimeoutPacket,
		channelTypes.EventTypeChannelOpenInit,
		channelTypes.EventTypeChannelOpenTry,
		channelTypes.EventTypeChannelOpenAck,
		channelTypes.EventTypeChannelOpenConfirm,
		channelTypes.EventTypeChannelCloseInit,
		channelTypes.EventTypeChannelCloseConfirm,
		channelTypes.EventTypeChannelClosed,
		channelTypes.EventTypeChannelUpgradeInit,
		channelTypes.EventTypeChannelUpgradeTry,
		channelTypes.EventTypeChannelUpgradeAck,
		channelTypes.EventTypeChannelUpgradeConfirm,
		channelTypes.EventTypeChannelUpgradeOpen,
		channelTypes.EventTypeChannelUpgradeTimeout,
		channelTypes.EventTypeChannelUpgradeCancel,
		channelTypes.EventTypeChannelUpgradeError,
		channelTypes.EventTypeChannelFlushComplete,

		transferTypes.EventTypeTimeout,
		transferTypes.EventTypePacket,
		transferTypes.EventTypeTransfer,
		transferTypes.EventTypeChannelClose,
		transferTypes.EventTypeDenomTrace,

		icaTypes.EventTypePacket,
	}

	// on an error ack, core re-emits the app's recv callback events with this prefix
	recvEvents := []string{
		transferTypes.EventTypePacket,
		transferTypes.EventTypeDenomTrace,
		icaTypes.EventTypePacket,
	}
	errorEvents := make([]string, 0, len(recvEvents))
	for _, e := range recvEvents {
		errorEvents = append(errorEvents, coreTypes.ErrorAttributeKeyPrefix+e)
	}

	for _, e := range slices.Concat(ibcEvents, errorEvents) {
		t.Run(e, func(t *testing.T) {
			typ, err := ParseEventType(e)
			require.NoError(t, err)
			require.NotEqual(t, EventTypeUnknown, typ)
		})
	}
}

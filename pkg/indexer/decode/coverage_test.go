// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"sort"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celestiaorg/celestia-app/v10/app"
	"github.com/celestiaorg/celestia-app/v10/app/encoding"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// knownUnhandledMsgs are messages registered by the app that Message() still
// falls through to MsgUnknown for. They are listed explicitly so a version bump
// surfaces new types instead of silently indexing them as unknown.
var knownUnhandledMsgs = map[string]struct{}{
	// Params updates of modules whose params the indexer does not track.
	"/cosmos.auth.v1beta1.MsgUpdateParams":                        {},
	"/cosmos.consensus.v1.MsgUpdateParams":                        {},
	"/ibc.applications.transfer.v1.MsgUpdateParams":               {},
	"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":          {},
	"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPool": {},
	"/cosmos.feegrant.v1beta1.MsgPruneAllowances":                 {},
	"/cosmos.gov.v1.MsgCancelProposal":                            {},

	// Hyperlane routing ISM management and native synthetic tokens.
	"/hyperlane.core.interchain_security.v1.MsgRemoveRoutingIsmDomain": {},
	"/hyperlane.core.interchain_security.v1.MsgSetRoutingIsmDomain":    {},
	"/hyperlane.core.interchain_security.v1.MsgUpdateRoutingIsmOwner":  {},
	"/hyperlane.warp.v1.MsgCreateNativeSyntheticToken":                 {},
	// Registered as a Msg implementation upstream although it is a response.
	"/hyperlane.warp.v1.MsgSetTokenResponse": {},

	// IBC channel upgrade handshake (ibc-go v8) and ack pruning.
	"/ibc.core.channel.v1.MsgChannelUpgradeAck":     {},
	"/ibc.core.channel.v1.MsgChannelUpgradeCancel":  {},
	"/ibc.core.channel.v1.MsgChannelUpgradeConfirm": {},
	"/ibc.core.channel.v1.MsgChannelUpgradeInit":    {},
	"/ibc.core.channel.v1.MsgChannelUpgradeOpen":    {},
	"/ibc.core.channel.v1.MsgChannelUpgradeTimeout": {},
	"/ibc.core.channel.v1.MsgChannelUpgradeTry":     {},
	"/ibc.core.channel.v1.MsgPruneAcknowledgements": {},
}

// TestMessageCoverage decodes an empty instance of every sdk.Msg the app
// registers and reports the ones Message() does not recognise. Empty messages
// make most handlers fail on address parsing, which is fine: the type is
// resolved by the switch before any handler runs.
func TestMessageCoverage(t *testing.T) {
	cfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	impls := cfg.InterfaceRegistry.ListImplementations(cosmosTypes.MsgInterfaceProtoName)
	require.NotEmpty(t, impls)

	block, now := testsuite.EmptyBlock()

	var unhandled, handledButListed []string
	for _, url := range impls {
		resolved, err := cfg.InterfaceRegistry.Resolve(url)
		require.NoErrorf(t, err, "cannot resolve %s", url)

		decodeCtx := context.NewContext()
		decodeCtx.Block = &storage.Block{Height: block.Height, Time: now}

		// Handlers may return an error on the empty message, but the type is
		// assigned by the switch regardless.
		dm, _ := Message(decodeCtx, resolved, 0, storageTypes.StatusSuccess, 0)

		_, listed := knownUnhandledMsgs[url]
		switch {
		case dm.Msg.Type == storageTypes.MsgUnknown && !listed:
			unhandled = append(unhandled, url)
		case dm.Msg.Type != storageTypes.MsgUnknown && listed:
			handledButListed = append(handledButListed, url)
		}
	}

	sort.Strings(unhandled)
	sort.Strings(handledButListed)
	require.Empty(t, unhandled, "messages decode to MsgUnknown; add a handler or list them in knownUnhandledMsgs")
	require.Empty(t, handledButListed, "messages are handled now; drop them from knownUnhandledMsgs")
}

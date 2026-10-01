// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle_test

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	decodeTestutil "github.com/celenium-io/celestia-indexer/pkg/indexer/decode/testutil"
	cosmosGovTypesV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
)

const cancelProposer = "celestia10d07y265gmmuvt4z0w9aw880jnsr700jtgz4v7"

func TestDecodeMsg_SuccessOnMsgCancelProposal(t *testing.T) {
	m := &cosmosGovTypesV1.MsgCancelProposal{
		ProposalId: 5,
		Proposer:   cancelProposer,
	}
	block, now := testsuite.EmptyBlock()
	position := 3

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, m, position, storageTypes.StatusSuccess, 0)
	require.NoError(t, err)

	msgExpected := decodeTestutil.CreateExpectations(
		t,
		block, now, m, position,
		storageTypes.MsgCancelProposal,
		m.Size(),
	)
	msgExpected.Id = dm.Msg.Id
	require.Equal(t, msgExpected, dm.Msg)
	require.EqualValues(t, 0, dm.BlobsSize)

	addresses := decodeCtx.AddressMessages.Values()
	require.Len(t, addresses, 1)
	require.Equal(t, storageTypes.MsgAddressTypeProposer, addresses[0].Type)
	require.Equal(t, cancelProposer, addresses[0].Address.Address)
	require.Equal(t, dm.Msg.Id, addresses[0].MsgId)

	// the status comes from the cancel_proposal event, not from the message
	require.Zero(t, decodeCtx.Proposals.Len())
}

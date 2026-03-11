// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handle

import (
	"testing"

	evidenceTypes "cosmossdk.io/x/evidence/types"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgSubmitEvidence

func createMsgSubmitEvidence() types.Msg {
	m := evidenceTypes.MsgSubmitEvidence{
		Submitter: "celestia1j33593mn9urzydakw06jdun8f37shlucmhr8p6",
	}

	return &m
}

func TestDecodeMsg_SuccessOnMsgSubmitEvidence(t *testing.T) {
	msg := createMsgSubmitEvidence()
	block, now := testsuite.EmptyBlock()
	position := 0

	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := decode.Message(decodeCtx, msg, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgSubmitEvidence,
		TxId:      0,
		Size:      49,
		Data:      structs.Map(msg),
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

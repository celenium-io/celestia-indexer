// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package decode

import (
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

// MsgUnknown

type UnknownMsgType struct{}

func (u *UnknownMsgType) Reset()                               {}
func (u *UnknownMsgType) String() string                       { return "unknown" }
func (u *UnknownMsgType) ProtoMessage()                        {}
func (u *UnknownMsgType) ValidateBasic() error                 { return nil }
func (u *UnknownMsgType) GetSigners() []cosmosTypes.AccAddress { return nil }

func createMsgUnknown() cosmosTypes.Msg {
	msgUnknown := UnknownMsgType{}
	return &msgUnknown
}

func TestDecodeMsg_MsgUnknown(t *testing.T) {
	msgUnknown := createMsgUnknown()
	block, now := testsuite.EmptyBlock()
	position := 0
	decodeCtx := context.NewContext()
	decodeCtx.Block = &storage.Block{
		Height: block.Height,
		Time:   block.Block.Time,
	}

	dm, err := Message(decodeCtx, msgUnknown, position, storageTypes.StatusSuccess, 0)

	msgExpected := storage.Message{
		Id:        1,
		Height:    block.Height,
		Time:      now,
		Position:  0,
		Type:      storageTypes.MsgUnknown,
		TxId:      0,
		Data:      structs.Map(msgUnknown),
		Namespace: nil,
	}

	require.NoError(t, err)
	require.Equal(t, int64(0), dm.BlobsSize)
	require.Equal(t, msgExpected, dm.Msg)
}

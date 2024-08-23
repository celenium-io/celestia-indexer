// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blobsaver

import (
	"context"
	"os"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/blob"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	sqBlob "github.com/celestiaorg/go-square/blob"
	"github.com/celestiaorg/go-square/inclusion"
	"github.com/celestiaorg/go-square/merkle"
	"github.com/celestiaorg/go-square/namespace"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	blobTypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

const (
	InputName  = "blobs"
	StopOutput = "stop"
)

type Msg struct {
	*blobTypes.Blob

	Height   pkgTypes.Level
	EndBlock bool
}

// Module - saves received from input block to storage.
//
//	                     |----------------|
//	                     |                |
//	-- storage.Block ->  |     MODULE     |
//	                     |                |
//	                     |----------------|
type Module struct {
	modules.BaseModule

	kind    string
	blocks  *sync.Map[pkgTypes.Level, *[]blob.Blob]
	blobs   *sync.Map[string, struct{}]
	storage blob.Storage
	head    pkgTypes.Level
}

var _ modules.Module = (*Module)(nil)

// NewModule -
func NewModule(
	kind string,
) (*Module, error) {
	m := Module{
		BaseModule: modules.New("blob_saver"),
		blocks:     sync.NewMap[pkgTypes.Level, *[]blob.Blob](),
		blobs:      sync.NewMap[string, struct{}](),
		kind:       kind,
	}

	m.CreateInputWithCapacity(InputName, 1024)
	m.CreateOutput(StopOutput)

	return &m, nil
}

// Start -
func (module *Module) Start(ctx context.Context) {
	if module.kind == "" {
		return
	}
	if err := module.init(ctx); err != nil {
		panic(err)
	}
	module.G.GoCtx(ctx, module.listen)
}

func (module *Module) init(ctx context.Context) error {
	switch module.kind {
	case "r2":
		r2 := blob.NewR2(blob.R2Config{
			BucketName:      os.Getenv("R2_BUCKET"),
			AccountId:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyId:     os.Getenv("R2_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("R2_ACCESS_KEY_SECRET"),
		})
		if err := r2.Init(ctx); err != nil {
			return errors.Wrap(err, "r2 initialization")
		}
		module.storage = r2
	case "mock":
	default:
		return errors.Errorf("unknown blob saver datasource: %s", module.kind)
	}
	head, err := module.storage.Head(ctx)
	if err != nil {
		return errors.Wrap(err, "can't receive head")
	}
	module.head = pkgTypes.Level(head)
	return nil
}

func (module *Module) listen(ctx context.Context) {
	module.Log.Info().Msg("module started")
	input := module.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				module.Log.Warn().Msg("can't read message from input")
				module.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
			message, ok := msg.(*Msg)
			if !ok {
				module.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			if err := module.processMessage(ctx, message); err != nil {
				module.Log.Err(err).Msg("blob processing")
			}
		}
	}
}

func (module *Module) processMessage(ctx context.Context, msg *Msg) error {
	if msg.Height <= module.head {
		return nil
	}

	if msg.EndBlock {
		return module.processEndOfBlock(ctx, msg.Height)
	}
	return module.processBlob(msg)
}

func (module *Module) processEndOfBlock(ctx context.Context, height pkgTypes.Level) error {
	if blobs, ok := module.blocks.Get(height); ok {
		if err := module.storage.SaveBulk(ctx, *blobs); err != nil {
			return errors.Wrap(err, "can't save blobs")
		}

		success := false
		for !success {
			timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			if err := module.storage.UpdateHead(timeoutCtx, uint64(height)); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					// try again
					module.Log.Err(err).Msg("update head")
					continue
				}
				return errors.Wrap(err, "can't update head")
			}

			success = true
		}
	}

	module.head = height
	module.blocks.Delete(height)
	module.blobs.Clear()
	return nil
}

func (module *Module) processBlob(msg *Msg) error {
	ns, err := namespace.New(uint8(msg.Blob.NamespaceVersion), msg.Blob.NamespaceId)
	if err != nil {
		return errors.Wrapf(err, "can't parse namespace: version=%d id=%x", msg.Blob.NamespaceVersion, msg.Blob.NamespaceId)
	}
	b := sqBlob.New(ns, msg.Blob.Data, uint8(msg.Blob.ShareVersion))
	commitment, err := inclusion.CreateCommitment(b, merkle.HashFromByteSlices, appconsts.SubtreeRootThreshold(uint64(msg.Blob.ShareVersion)))
	if err != nil {
		return errors.Wrap(err, "can't create commitment")
	}

	blb := blob.Blob{
		Commitment: commitment,
		Blob:       msg.Blob,
		Height:     uint64(msg.Height),
	}
	key := blb.String()

	// skip blobs with the same commitments in the current block.
	if _, ok := module.blobs.Get(key); ok {
		return nil
	}
	module.blobs.Set(key, struct{}{})

	if blobs, ok := module.blocks.Get(msg.Height); ok {
		*blobs = append(*blobs, blb)
	} else {
		module.blocks.Set(msg.Height, &[]blob.Blob{blb})
	}

	return nil
}

// Close -
func (module *Module) Close() error {
	module.Log.Info().Msg("closing module...")
	module.G.Wait()
	return nil
}

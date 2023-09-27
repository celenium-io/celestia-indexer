package websocket

import (
	"context"
	"strconv"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

type blockRepo struct {
	repo storage.IBlock
}

func newBlockRepo(repo storage.IBlock) blockRepo {
	return blockRepo{repo}
}

func (block blockRepo) GetById(ctx context.Context, id uint64) (storage.Block, error) {
	b, err := block.repo.ByIdWithRelations(ctx, id)
	if err != nil {
		return storage.Block{}, err
	}
	return b, nil
}

func HeadProcessor(ctx context.Context, payload string, repo identifiable[storage.Block]) (*responses.Block, error) {
	blockId, err := strconv.ParseUint(payload, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse block id")
	}

	b, err := repo.GetById(ctx, blockId)
	if err != nil {
		return nil, errors.Wrap(err, "receive block by id")
	}

	block := responses.NewBlock(b, true)
	return &block, nil
}

type txRepo struct {
	repo storage.ITx
}

func newTxRepo(repo storage.ITx) txRepo {
	return txRepo{repo}
}

func (tx txRepo) GetById(ctx context.Context, id uint64) (storage.Tx, error) {
	return tx.repo.ByIdWithRelations(ctx, id)
}

func TxProcessor(ctx context.Context, payload string, repo identifiable[storage.Tx]) (*responses.Tx, error) {
	txId, err := strconv.ParseUint(payload, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "parse tx id")
	}

	tx, err := repo.GetById(ctx, txId)
	if err != nil {
		return nil, errors.Wrap(err, "receive transaction by id")
	}

	response := responses.NewTx(tx)
	return &response, nil
}

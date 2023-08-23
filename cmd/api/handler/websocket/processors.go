package websocket

import (
	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func HeadProcessor(payload string) (responses.Block, error) {
	var b storage.Block
	if err := json.Unmarshal([]byte(payload), &b); err != nil {
		return responses.Block{}, errors.Errorf("block unmarhaling in processor: %s", err.Error())
	}

	return responses.NewBlock(b), nil
}

func TxProcessor(payload string) (responses.Tx, error) {
	var tx storage.Tx
	if err := json.Unmarshal([]byte(payload), &tx); err != nil {
		return responses.Tx{}, errors.Errorf("block unmarhaling in processor: %s", err.Error())
	}

	return responses.NewTx(tx), nil
}

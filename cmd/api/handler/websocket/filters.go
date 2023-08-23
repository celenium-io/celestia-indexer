package websocket

import (
	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

type Filterable[M any] interface {
	Filter(c *Client, msg M) bool
}

type HeadFilter struct{}

func (hf HeadFilter) Filter(c *Client, msg responses.Block) bool {
	return c.filters.head
}

type TxFilter struct{}

func (hf TxFilter) Filter(c *Client, msg responses.Tx) bool {
	fltr := c.filters.tx
	if fltr == nil {
		return false
	}
	if len(fltr.status) > 0 {
		if _, ok := fltr.status[msg.Status]; !ok {
			return false
		}
	}
	return true
}

type filters struct {
	head bool
	tx   *txFilters
}

func newFilters() *filters {
	return &filters{
		tx: newTxFilters(),
	}
}

type txFilters struct {
	status map[string]struct{}
}

func newTxFilters() *txFilters {
	return &txFilters{
		status: make(map[string]struct{}, 2),
	}
}

var (
	availiableStatus = map[string]struct{}{
		string(storage.StatusSuccess): {},
		string(storage.StatusFailed):  {},
	}
)

func (f *txFilters) Fill(msg TransactionFilters) error {
	for i := range msg.Status {
		if _, ok := availiableStatus[msg.Status[i]]; !ok {
			return errors.Wrapf(ErrUnavailiableFilter, "status %s", msg.Status[i])
		}
		f.status[msg.Status[i]] = struct{}{}
	}

	return nil
}

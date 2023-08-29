package websocket

import (
	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
)

type Filterable[M any] interface {
	Filter(c *Client, msg M) bool
}

type HeadFilter struct{}

func (hf HeadFilter) Filter(c *Client, msg responses.Block) bool {
	if c.filters == nil {
		return false
	}
	return c.filters.head
}

type TxFilter struct{}

func (hf TxFilter) Filter(c *Client, msg responses.Tx) bool {
	if c.filters == nil || c.filters.tx == nil {
		return false
	}

	fltr := c.filters.tx
	if len(fltr.status) > 0 {
		if _, ok := fltr.status[msg.Status]; !ok {
			return false
		}
	}
	if !fltr.msgs.Empty() {
		if !fltr.msgs.HasOne(msg.MsgTypeMask) {
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
	msgs   types.MsgTypeBits
}

func newTxFilters() *txFilters {
	return &txFilters{
		status: make(map[string]struct{}, 2),
		msgs:   types.NewMsgTypeBitMask(),
	}
}

func (f *txFilters) Fill(msg TransactionFilters) error {
	for i := range msg.Status {
		if !types.IsStatus(msg.Status[i]) {
			return errors.Wrapf(ErrUnavailiableFilter, "status %s", msg.Status[i])
		}
		f.status[msg.Status[i]] = struct{}{}
	}

	for i := range msg.Messages {
		f.msgs.SetBit(types.MsgType(msg.Messages[i]))
	}

	return nil
}

package websocket

import (
	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
)

type Filterable[M any] interface {
	Filter(c client, msg M) bool
}

type HeadFilter struct{}

func (hf HeadFilter) Filter(c client, msg *responses.Block) bool {
	if msg == nil {
		return false
	}
	fltrs := c.Filters()
	if fltrs == nil {
		return false
	}
	return fltrs.head
}

type TxFilter struct{}

func (hf TxFilter) Filter(c client, msg *responses.Tx) bool {
	if msg == nil {
		return false
	}
	fltrs := c.Filters()
	if fltrs == nil || fltrs.tx == nil {
		return false
	}

	fltr := fltrs.tx
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

type Filters struct {
	head bool
	tx   *txFilters
}

func newFilters() *Filters {
	return &Filters{
		tx: newTxFilters(),
	}
}

type txFilters struct {
	status map[types.Status]struct{}
	msgs   types.MsgTypeBits
}

func newTxFilters() *txFilters {
	return &txFilters{
		status: make(map[types.Status]struct{}, 2),
		msgs:   types.NewMsgTypeBitMask(),
	}
}

func (f *txFilters) Fill(msg TransactionFilters) error {
	for i := range msg.Status {
		status, err := types.ParseStatus(msg.Status[i])
		if err != nil {
			return errors.Wrapf(ErrUnavailableFilter, "status %s", msg.Status[i])
		}
		f.status[status] = struct{}{}
	}

	for i := range msg.Messages {
		f.msgs.SetBit(types.MsgType(msg.Messages[i]))
	}

	return nil
}

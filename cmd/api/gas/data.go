package gas

import (
	"sync"

	"github.com/shopspring/decimal"
)

type GasPrice struct {
	Slow   string
	Median string
	Fast   string
}

type info struct {
	Height       uint64
	Percentiles  []decimal.Decimal
	TxCount      int64
	GasUsed      int64
	GasWanted    int64
	Fee          decimal.Decimal
	GasUsedRatio decimal.Decimal
}

type queue struct {
	data     []info
	capacity int
	mx       *sync.RWMutex
}

func newQueue(capacity int) *queue {
	return &queue{
		data:     make([]info, 0),
		capacity: capacity,
		mx:       new(sync.RWMutex),
	}
}

func (q *queue) Push(item info) {
	q.mx.Lock()
	if len(q.data) == q.capacity {
		q.data = q.data[:len(q.data)-2]
	}
	q.data = append([]info{item}, q.data...)
	q.mx.Unlock()
}

func (q *queue) Range(handler func(item info) (bool, error)) error {
	if handler == nil {
		return nil
	}

	q.mx.RLock()
	defer q.mx.RUnlock()

	for i := range q.data {
		br, err := handler(q.data[i])
		if err != nil {
			return err
		}
		if br {
			return nil
		}
	}
	return nil
}

func (q *queue) Size() int {
	q.mx.RLock()
	defer q.mx.RUnlock()

	return len(q.data)
}

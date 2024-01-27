// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"sync"
	"testing"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const clientsTestCount = uint64(1000)

type testHeadClient struct {
	id    uint64
	fltrs *Filters
	ch    chan any
	wg    *sync.WaitGroup
}

func newTestHeadClient(id uint64) *testHeadClient {
	return &testHeadClient{
		id: id,
		fltrs: &Filters{
			head: true,
		},
		ch: make(chan any, 1024),
		wg: new(sync.WaitGroup),
	}
}

func (c *testHeadClient) Id() uint64 {
	return c.id
}

func (c *testHeadClient) Filters() *Filters {
	return c.fltrs
}

func (c *testHeadClient) ApplyFilters(msg Subscribe) error {
	return nil
}

func (c *testHeadClient) DetachFilters(msg Unsubscribe) error {
	return nil
}

func (c *testHeadClient) Notify(msg any) {
	c.ch <- msg
}

func (c *testHeadClient) WriteMessages(ctx context.Context, ws *websocket.Conn, log echo.Logger) {
	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.ch:
			continue
		}
	}
}

func (c *testHeadClient) ReadMessages(ctx context.Context, ws *websocket.Conn, sub *Client, log echo.Logger) {
}

func (c *testHeadClient) Close() error {
	c.wg.Wait()
	close(c.ch)
	return nil
}

func BenchmarkProcessingMessage(b *testing.B) {
	channel := NewChannel[storage.Block, *responses.Block](
		blockProcessor,
		BlockFilter{},
	)

	ctx, cancel := context.WithCancel(context.Background())

	for id := uint64(0); id < clientsTestCount; id++ {
		client := newTestHeadClient(id)
		channel.clients.Set(id, client)
		go client.WriteMessages(ctx, nil, nil)
	}

	b.Run("websocket_process_message", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			height := uint64(i)

			err := channel.processMessage(storage.Block{
				Id:           height,
				Height:       types.Level(height),
				MessageTypes: storageTypes.NewMsgTypeBits(),
			})
			require.NoError(b, err)
		}
	})

	cancel()

	err := channel.clients.Range(func(_ uint64, value client) (error, bool) {
		return value.Close(), false
	})
	require.NoError(b, err)
}

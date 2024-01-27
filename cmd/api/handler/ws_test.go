// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	ws "github.com/celenium-io/celestia-indexer/cmd/api/handler/websocket"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWebsocket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listenerFactory := mock.NewMockListenerFactory(ctrl)
	listener := mock.NewMockListener(ctrl)

	listenerFactory.EXPECT().CreateListener().Return(listener).Times(1)

	headChannel := make(chan *pq.Notification, 10)
	listener.EXPECT().Listen().Return(headChannel).AnyTimes()
	listener.EXPECT().Subscribe(gomock.Any(), storage.ChannelHead).Return(nil).Times(1)
	listener.EXPECT().Close().Return(nil).MaxTimes(1)

	ctx, cancel := context.WithCancel(context.Background())

	blockMock := mock.NewMockIBlock(ctrl)
	dispatcher, err := bus.NewDispatcher(listenerFactory, blockMock)
	require.NoError(t, err)
	dispatcher.Start(ctx)
	observer := dispatcher.Observe(storage.ChannelHead, storage.ChannelBlock)

	for i := 0; i < 10; i++ {
		hash := make([]byte, 32)
		_, err := rand.Read(hash)
		require.NoError(t, err)

		blockMock.EXPECT().ByIdWithRelations(ctx, uint64(i)).Return(storage.Block{
			Id:           uint64(i),
			Height:       types.Level(i),
			Time:         time.Now(),
			Hash:         hash,
			MessageTypes: storageTypes.NewMsgTypeBits(),
			Stats:        testBlock.Stats,
		}, nil).MaxTimes(1)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var id uint64

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				id++

				headChannel <- &pq.Notification{
					Channel: storage.ChannelBlock,
					Extra:   strconv.FormatUint(id, 10),
				}
			}
		}
	}()
	manager := ws.NewManager(observer)
	manager.Start(ctx)

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			e := echo.New()
			c := e.NewContext(r, w)
			err := manager.Handle(c)
			require.NoError(t, err, "handle")
			<-ctx.Done()
		},
	))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	dialed, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "dial")

	body, err := json.Marshal(ws.Subscribe{
		Channel: ws.ChannelBlocks,
	})
	require.NoError(t, err, "marshal subscribe")

	err = dialed.WriteJSON(ws.Message{
		Method: ws.MethodSubscribe,
		Body:   body,
	})
	require.NoError(t, err, "send subscribe message")

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			body, err := json.Marshal(ws.Unsubscribe{
				Channel: ws.ChannelHead,
			})
			require.NoError(t, err, "marshal unsubscribe")

			err = dialed.WriteJSON(ws.Message{
				Method: ws.MethodUnsubscribe,
				Body:   body,
			})
			require.NoError(t, err, "send unsubscribe message")

			err = dialed.Close()
			require.NoError(t, err, "closing connection")

			time.Sleep(time.Second)
			cancel()

			err = manager.Close()
			require.NoError(t, err, "closing manager")

			close(headChannel)
			return
		default:
			err := dialed.SetReadDeadline(time.Now().Add(time.Second * 3))
			require.NoError(t, err)

			_, msg, err := dialed.ReadMessage()
			require.NoError(t, err, err)

			var block responses.Block
			err = json.Unmarshal(msg, &block)
			require.NoError(t, err, err)

			require.Greater(t, block.Id, uint64(0))
			require.Greater(t, block.Height, uint64(0))
			require.False(t, block.Time.IsZero())
			require.Len(t, block.Hash, 32)

			log.Info().
				Uint64("height", block.Height).
				Time("block_time", block.Time).
				Msg("new block")
		}
	}
}

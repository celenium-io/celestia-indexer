package handler

import (
	"context"
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	ws "github.com/dipdup-io/celestia-indexer/cmd/api/handler/websocket"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
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
	headListener := mock.NewMockListener(ctrl)
	txListener := mock.NewMockListener(ctrl)

	listenerFactory.EXPECT().CreateListener().Return(headListener).MaxTimes(1)
	listenerFactory.EXPECT().CreateListener().Return(txListener).MaxTimes(1)

	ctx, cancel := context.WithCancel(context.Background())

	headChannel := make(chan *pq.Notification, 10)
	headListener.EXPECT().Listen().Return(headChannel).AnyTimes()
	txChannel := make(chan *pq.Notification, 10)
	txListener.EXPECT().Listen().Return(txChannel).AnyTimes()

	headListener.EXPECT().Subscribe(gomock.Any(), storage.ChannelHead).Return(nil).MaxTimes(1)
	txListener.EXPECT().Subscribe(gomock.Any(), storage.ChannelTx).Return(nil).MaxTimes(1)

	headListener.EXPECT().Close().Return(nil).MaxTimes(1)
	txListener.EXPECT().Close().Return(nil).MaxTimes(1)

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

				hash := make([]byte, 32)
				_, err := rand.Read(hash)
				require.NoError(t, err)

				payload, err := json.Marshal(storage.Block{
					Id:     id,
					Height: id,
					Time:   time.Now(),
					Hash:   hash,
				})
				require.NoError(t, err)

				headChannel <- &pq.Notification{
					Channel: storage.ChannelHead,
					Extra:   string(payload),
				}
			}
		}
	}()
	manager := ws.NewManager(listenerFactory)
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
		Channel: ws.ChannelHead,
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
			close(txChannel)
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
			require.Len(t, block.Hash, 64)

			log.Info().
				Uint64("height", block.Height).
				Time("block_time", block.Time).
				Msg("new block")
		}
	}
}

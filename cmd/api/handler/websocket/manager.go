// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/dipdup-io/workerpool"
	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	upgrader websocket.Upgrader
	clientId *atomic.Uint64
	clients  *sdkSync.Map[uint64, *Client]
	observer *bus.Observer

	blocks *Channel[storage.Block, *responses.Block]
	head   *Channel[storage.State, *responses.State]

	g workerpool.Group
}

func NewManager(observer *bus.Observer) *Manager {
	manager := &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		observer: observer,
		clientId: new(atomic.Uint64),
		clients:  sdkSync.NewMap[uint64, *Client](),
		g:        workerpool.NewGroup(),
	}

	manager.blocks = NewChannel[storage.Block, *responses.Block](
		blockProcessor,
		BlockFilter{},
	)

	manager.head = NewChannel[storage.State, *responses.State](
		headProcessor,
		HeadFilter{},
	)

	return manager
}

func (manager *Manager) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-manager.observer.Blocks():
			if err := manager.blocks.processMessage(*block); err != nil {
				log.Err(err).Msg("handle block")
			}
		case state := <-manager.observer.Head():
			if err := manager.head.processMessage(*state); err != nil {
				log.Err(err).Msg("handle state")
			}
		}
	}
}

// Handle godoc
//
//	@Summary				Websocket API
//	@Description.markdown	websocket
//	@Tags					websocket
//	@ID						websocket
//	@Produce				json
//	@Router					/v1/ws [get]
func (manager *Manager) Handle(c echo.Context) error {
	ws, err := manager.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	ws.SetReadLimit(1024 * 10) // 10KB

	sId := manager.clientId.Add(1)
	sub := newClient(sId, manager)

	manager.clients.Set(sId, sub)

	ctx, cancel := context.WithCancel(c.Request().Context())
	sub.WriteMessages(ctx, ws, c.Logger())
	sub.ReadMessages(ctx, ws, sub, c.Logger())
	manager.clients.Delete(sId)
	cancel()

	if err := sub.Close(); err != nil {
		return err
	}

	return ws.Close()
}

func (manager *Manager) Start(ctx context.Context) {
	manager.g.GoCtx(ctx, manager.listen)
}

func (manager *Manager) Close() error {
	manager.g.Wait()

	return manager.clients.Range(func(_ uint64, value *Client) (error, bool) {
		if err := value.Close(); err != nil {
			return err, false
		}
		return nil, false
	})
}

func (manager *Manager) AddClientToChannel(channel string, client *Client) {
	switch channel {
	case ChannelHead:
		manager.head.AddClient(client)
	case ChannelBlocks:
		manager.blocks.AddClient(client)
	default:
		log.Error().Str("channel", channel).Msg("unknown channel name")
	}
}

func (manager *Manager) RemoveClientFromChannel(channel string, client *Client) {
	switch channel {
	case ChannelHead:
		manager.head.RemoveClient(client.id)
	case ChannelBlocks:
		manager.blocks.RemoveClient(client.id)
	default:
		log.Error().Str("channel", channel).Msg("unknown channel name")
	}
}

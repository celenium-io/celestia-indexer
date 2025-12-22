// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	upgrader websocket.Upgrader
	clientId *atomic.Uint64
	clients  *sync.Map[uint64, *Client]
	observer *bus.Observer

	ips                   *sync.Map[string, int]
	websocketClientsPerIp int

	blocks   *Channel[storage.Block, *responses.Block]
	head     *Channel[storage.State, *responses.State]
	gasPrice *Channel[gas.GasPrice, *responses.GasPrice]

	g workerpool.Group
}

func NewManager(observer *bus.Observer, opts ...ManagerOption) *Manager {
	manager := &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		observer:              observer,
		clientId:              new(atomic.Uint64),
		clients:               sync.NewMap[uint64, *Client](),
		g:                     workerpool.NewGroup(),
		ips:                   sync.NewMap[string, int](),
		websocketClientsPerIp: 10,
	}

	manager.blocks = NewChannel(
		blockProcessor,
		BlockFilter{},
	)

	manager.head = NewChannel(
		headProcessor,
		HeadFilter{},
	)

	manager.gasPrice = NewChannel(
		gasPriceProcessor,
		GasPriceFilter{},
	)

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

func (manager *Manager) listenBlocks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-manager.observer.Blocks():
			if err := manager.blocks.processMessage(*block); err != nil {
				log.Err(err).Msg("handle block")
			}
		}
	}
}

func (manager *Manager) listenHead(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case state := <-manager.observer.Head():
			if err := manager.head.processMessage(*state); err != nil {
				log.Err(err).Msg("handle state")
			}
		}
	}
}

func (manager *Manager) countClientsByIp(ip string, value int) error {
	if count, ok := manager.ips.Get(ip); ok {
		if count >= manager.websocketClientsPerIp {
			return ErrTooManyClients
		}
		manager.ips.Set(ip, count+value)
	} else {
		manager.ips.Set(ip, value)
	}
	return nil
}

// Handle godoc
//
//	@Summary				Websocket API
//	@Description.markdown	websocket
//	@Tags					websocket
//	@ID						websocket
//	@x-internal				true
//	@Produce				json
//	@Router					/ws [get]
func (manager *Manager) Handle(c echo.Context) error {
	ws, err := manager.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		wsErrors.WithLabelValues("upgrade").Inc()
		return err
	}
	ws.SetReadLimit(1024 * 10) // 10KB

	if err := manager.countClientsByIp(c.RealIP(), 1); err != nil {
		wsConnectionsTotal.WithLabelValues("rejected").Inc()
		return err
	}
	defer func() {
		if err := manager.countClientsByIp(c.RealIP(), -1); err != nil {
			log.Err(err).Msg("decrease client count by ip")
		}
	}()

	wsConnectionsTotal.WithLabelValues("accepted").Inc()
	wsActiveConnections.Inc()

	startTime := time.Now()
	defer func() {
		wsActiveConnections.Dec()
		wsConnectionDuration.Observe(time.Since(startTime).Seconds())
	}()

	sId := manager.clientId.Add(1)
	sub := newClient(sId, manager.AddClientToChannel, manager.RemoveClientFromChannel)

	manager.clients.Set(sId, sub)

	ctx, cancel := context.WithCancel(c.Request().Context())
	sub.WriteMessages(ctx, ws, c.Logger())
	sub.ReadMessages(ctx, ws, c.Logger())
	manager.clients.Delete(sId)
	cancel()

	if err := sub.Close(); err != nil {
		return err
	}

	return ws.Close()
}

func (manager *Manager) Start(ctx context.Context) {
	manager.g.GoCtx(ctx, manager.listenHead)
	manager.g.GoCtx(ctx, manager.listenBlocks)
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
		wsSubscriptions.WithLabelValues(channel).Inc()
	case ChannelBlocks:
		manager.blocks.AddClient(client)
		wsSubscriptions.WithLabelValues(channel).Inc()
	case ChannelGasPrice:
		manager.gasPrice.AddClient(client)
		wsSubscriptions.WithLabelValues(channel).Inc()
	default:
		log.Error().Str("channel", channel).Msg("unknown channel name")
		wsErrors.WithLabelValues("unknown_channel").Inc()
	}
}

func (manager *Manager) RemoveClientFromChannel(channel string, client *Client) {
	switch channel {
	case ChannelHead:
		manager.head.RemoveClient(client.id)
		wsSubscriptions.WithLabelValues(channel).Dec()
	case ChannelBlocks:
		manager.blocks.RemoveClient(client.id)
		wsSubscriptions.WithLabelValues(channel).Dec()
	case ChannelGasPrice:
		manager.gasPrice.RemoveClient(client.id)
		wsSubscriptions.WithLabelValues(channel).Dec()
	default:
		log.Error().Str("channel", channel).Msg("unknown channel name")
	}
}

func (manager *Manager) GasTrackerHandler(_ context.Context, state gas.GasPrice) error {
	return manager.gasPrice.processMessage(state)
}

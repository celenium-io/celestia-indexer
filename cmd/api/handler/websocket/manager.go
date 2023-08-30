package websocket

import (
	"context"
	"sync/atomic"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	upgrader websocket.Upgrader
	clientId *atomic.Uint64
	clients  Map[uint64, *Client]

	head    *Channel[responses.Block]
	tx      *Channel[responses.Tx]
	factory storage.ListenerFactory
}

func NewManager(factory storage.ListenerFactory) *Manager {
	manager := &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clientId: new(atomic.Uint64),
		clients:  NewMap[uint64, *Client](),
		factory:  factory,
	}

	manager.head = NewChannel[responses.Block](
		storage.ChannelHead,
		HeadProcessor,
		HeadFilter{},
	)

	manager.tx = NewChannel[responses.Tx](
		storage.ChannelTx,
		TxProcessor,
		TxFilter{},
	)

	return manager
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
	sub := newClient(sId, ws, manager)

	manager.clients.Set(sId, sub)

	go sub.ReadMessages(c.Request().Context(), ws, sub, c.Logger())
	go sub.WriteMessages(c.Request().Context(), c.Logger())

	c.Logger().Infof("client %d connected", sId)
	return nil
}

func (manager *Manager) Start(ctx context.Context) {
	manager.head.Start(ctx, manager.factory)
	manager.tx.Start(ctx, manager.factory)
}

func (manager *Manager) Close() error {
	if err := manager.clients.Range(func(_ uint64, value *Client) (error, bool) {
		if err := value.Close(); err != nil {
			return err, false
		}
		return nil, false
	}); err != nil {
		return err
	}
	if err := manager.head.Close(); err != nil {
		return err
	}
	if err := manager.tx.Close(); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) AddClientToChannel(channel string, client *Client) {
	switch channel {
	case ChannelHead:
		manager.head.AddClient(client)
	case ChannelTx:
		manager.tx.AddClient(client)
	default:
		log.Error().Str("channel", channel).Msg("unknwon channel name")
	}
}

func (manager *Manager) RemoveClientFromChannel(channel string, client *Client) {
	switch channel {
	case ChannelHead:
		manager.head.RemoveClient(client.id)
	case ChannelTx:
		manager.tx.RemoveClient(client.id)
	default:
		log.Error().Str("channel", channel).Msg("unknwon channel name")
	}
}

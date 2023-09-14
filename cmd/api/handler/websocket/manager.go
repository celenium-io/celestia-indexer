package websocket

import (
	"context"
	"net/http"
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

	head    *Channel[storage.Block, *responses.Block]
	tx      *Channel[storage.Tx, *responses.Tx]
	factory storage.ListenerFactory
}

func NewManager(factory storage.ListenerFactory, blockRepo storage.IBlock, txRepo storage.ITx) *Manager {
	manager := &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clientId: new(atomic.Uint64),
		clients:  NewMap[uint64, *Client](),
		factory:  factory,
	}

	manager.head = NewChannel[storage.Block, *responses.Block](
		storage.ChannelHead,
		HeadProcessor,
		newBlockRepo(blockRepo),
		HeadFilter{},
	)

	manager.tx = NewChannel[storage.Tx, *responses.Tx](
		storage.ChannelTx,
		TxProcessor,
		newTxRepo(txRepo),
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
	sub := newClient(sId, manager)

	manager.clients.Set(sId, sub)

	ctx, cancel := context.WithCancel(c.Request().Context())
	sub.WriteMessages(ctx, ws, c.Logger())
	sub.ReadMessages(ctx, ws, sub, c.Logger())
	cancel()

	if err := sub.Close(); err != nil {
		return err
	}

	return ws.Close()
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

package websocket

import (
	"context"
	"io"
	"time"

	"github.com/dipdup-io/workerpool"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type client interface {
	Id() uint64
	ApplyFilters(msg Subscribe) error
	DetachFilters(msg Unsubscribe) error
	Notify(msg any)
	WriteMessages(ctx context.Context, ws *websocket.Conn, log echo.Logger)
	ReadMessages(ctx context.Context, ws *websocket.Conn, sub *Client, log echo.Logger)
	Filters() *Filters

	io.Closer
}

const (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	id      uint64
	manager *Manager
	filters *Filters
	ch      chan any
	g       workerpool.Group
}

func newClient(id uint64, manager *Manager) *Client {
	return &Client{
		id:      id,
		manager: manager,
		ch:      make(chan any, 1024),
		g:       workerpool.NewGroup(),
	}
}

func (c *Client) Id() uint64 {
	return c.id
}

func (c *Client) Filters() *Filters {
	return c.filters
}

func (c *Client) ApplyFilters(msg Subscribe) error {
	if c.filters == nil {
		c.filters = &Filters{}
	}
	switch msg.Channel {
	case ChannelHead:
		c.filters.head = true
	case ChannelTx:
		if c.filters.tx == nil {
			c.filters.tx = newTxFilters()
		}
		var fltr TransactionFilters
		if err := json.Unmarshal(msg.Filters, &fltr); err != nil {
			return err
		}
		if err := c.filters.tx.Fill(fltr); err != nil {
			return err
		}
	default:
		return errors.Wrap(ErrUnknownChannel, msg.Channel)
	}
	return nil
}

func (c *Client) DetachFilters(msg Unsubscribe) error {
	if c.filters == nil {
		return nil
	}
	switch msg.Channel {
	case ChannelHead:
		c.filters.head = false
	case ChannelTx:
		c.filters.tx = nil
	default:
		return errors.Wrap(ErrUnknownChannel, msg.Channel)
	}
	return nil
}

func (c *Client) Notify(msg any) {
	c.ch <- msg
}

func (c *Client) Close() error {
	c.g.Wait()
	close(c.ch)
	return nil
}

func (c *Client) writeThread(ctx context.Context, ws *websocket.Conn, log echo.Logger) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("writemsg: %s", err)
				return
			}

		case msg, ok := <-c.ch:
			if !ok {
				if err := ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Errorf("send close message: %s", err)
				}
				return
			}

			if err := ws.WriteJSON(msg); err != nil {
				log.Errorf("send client message: %s", err)
			}
		}
	}
}

func (c *Client) WriteMessages(ctx context.Context, ws *websocket.Conn, log echo.Logger) {
	c.g.GoCtx(ctx, func(ctx context.Context) {
		c.writeThread(ctx, ws, log)
	})
}

func (c *Client) ReadMessages(ctx context.Context, ws *websocket.Conn, sub *Client, log echo.Logger) {
	defer c.manager.clients.Delete(sub.id)

	if err := ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Error(err)
		return
	}
	ws.SetPongHandler(func(_ string) error {
		return ws.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.read(ctx, ws); err != nil {
				switch {
				case err == io.EOF:
					return
				case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway):
					return
				}
				log.Errorf("read websocket message: %s", err.Error())
			}
		}
	}
}

func (c *Client) read(ctx context.Context, ws *websocket.Conn) error {
	var msg Message
	if err := ws.ReadJSON(&msg); err != nil {
		return err
	}

	switch msg.Method {
	case MethodSubscribe:
		return c.handleSubscribeMessage(ctx, msg)
	case MethodUnsubscribe:
		return c.handleUnsubscribeMessage(ctx, msg)
	default:
		return errors.Wrap(ErrUnknownMethod, msg.Method)
	}
}

func (c *Client) handleSubscribeMessage(ctx context.Context, msg Message) error {
	var subscribeMsg Subscribe
	if err := json.Unmarshal(msg.Body, &subscribeMsg); err != nil {
		return err
	}

	if err := c.ApplyFilters(subscribeMsg); err != nil {
		return err
	}

	c.manager.AddClientToChannel(subscribeMsg.Channel, c)
	return nil
}

func (c *Client) handleUnsubscribeMessage(ctx context.Context, msg Message) error {
	var unsubscribeMsg Unsubscribe
	if err := json.Unmarshal(msg.Body, &unsubscribeMsg); err != nil {
		return err
	}
	if err := c.DetachFilters(unsubscribeMsg); err != nil {
		return err
	}
	c.manager.RemoveClientFromChannel(unsubscribeMsg.Channel, c)
	return nil
}

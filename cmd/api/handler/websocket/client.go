// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"io"
	"net"
	"sync/atomic"
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
	ReadMessages(ctx context.Context, ws *websocket.Conn, log echo.Logger)
	Filters() *Filters

	io.Closer
}

const (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is because otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type ClientHandler func(string, *Client)

type Client struct {
	id      uint64
	filters *Filters
	ch      chan any
	g       workerpool.Group

	subscribeHandler   ClientHandler
	unsubscribeHandler ClientHandler

	closed *atomic.Bool
}

func newClient(id uint64, subscribeHandler, unsubscribeHandler ClientHandler) *Client {
	closed := new(atomic.Bool)
	closed.Store(false)
	return &Client{
		id:                 id,
		ch:                 make(chan any, 128),
		g:                  workerpool.NewGroup(),
		subscribeHandler:   subscribeHandler,
		unsubscribeHandler: unsubscribeHandler,
		closed:             closed,
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
	case ChannelBlocks:
		c.filters.blocks = true
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
	case ChannelBlocks:
		c.filters.blocks = false
	default:
		return errors.Wrap(ErrUnknownChannel, msg.Channel)
	}
	return nil
}

func (c *Client) Notify(msg any) {
	if c.closed.Load() {
		return
	}
	c.ch <- msg
}

func (c *Client) Close() error {
	c.g.Wait()
	c.closed.Store(true)
	close(c.ch)
	c.subscribeHandler = nil
	c.unsubscribeHandler = nil
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

func (c *Client) ReadMessages(ctx context.Context, ws *websocket.Conn, log echo.Logger) {
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
			if err := c.read(ws); err != nil {
				timeoutErr, ok := err.(net.Error)

				switch {
				case err == io.EOF:
					return
				case errors.Is(err, websocket.ErrCloseSent):
					return
				case ok && timeoutErr.Timeout():
					return
				case websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseAbnormalClosure,
					websocket.CloseNoStatusReceived,
					websocket.CloseGoingAway):
					if c.unsubscribeHandler != nil {
						c.unsubscribeHandler(ChannelHead, c)
						c.unsubscribeHandler(ChannelBlocks, c)
					}
					return
				}
				log.Errorf("read websocket message: %s", err.Error())
			}
		}
	}
}

func (c *Client) read(ws *websocket.Conn) error {
	var msg Message
	if err := ws.ReadJSON(&msg); err != nil {
		return err
	}

	switch msg.Method {
	case MethodSubscribe:
		return c.handleSubscribeMessage(msg)
	case MethodUnsubscribe:
		return c.handleUnsubscribeMessage(msg)
	default:
		return errors.Wrap(ErrUnknownMethod, msg.Method)
	}
}

func (c *Client) handleSubscribeMessage(msg Message) error {
	var subscribeMsg Subscribe
	if err := json.Unmarshal(msg.Body, &subscribeMsg); err != nil {
		return err
	}

	if err := c.ApplyFilters(subscribeMsg); err != nil {
		return err
	}

	if c.unsubscribeHandler != nil {
		c.subscribeHandler(subscribeMsg.Channel, c)
	}
	return nil
}

func (c *Client) handleUnsubscribeMessage(msg Message) error {
	var unsubscribeMsg Unsubscribe
	if err := json.Unmarshal(msg.Body, &unsubscribeMsg); err != nil {
		return err
	}
	if err := c.DetachFilters(unsubscribeMsg); err != nil {
		return err
	}
	if c.unsubscribeHandler != nil {
		c.unsubscribeHandler(unsubscribeMsg.Channel, c)
	}
	return nil
}

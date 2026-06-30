// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	gosync "sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

const testIp = "10.0.0.1"

func dialWS(t *testing.T, srv *httptest.Server) *websocket.Conn {
	t.Helper()
	wsUrl := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, resp, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	require.NoError(t, err)
	return conn
}

func subscribe(t *testing.T, conn *websocket.Conn, channel string) {
	t.Helper()
	require.NoError(t, conn.WriteJSON(Message{
		Method: MethodSubscribe,
		Body:   json.RawMessage(`{"channel":"` + channel + `","filters":{}}`),
	}))
}

func unsubscribe(t *testing.T, conn *websocket.Conn, channel string) {
	t.Helper()
	require.NoError(t, conn.WriteJSON(Message{
		Method: MethodUnsubscribe,
		Body:   json.RawMessage(`{"channel":"` + channel + `"}`),
	}))
}

// TestHandleUnsubscribeMetrics verifies that unsubscribe requests are accounted
// with success/error status symmetrically with subscribe.
func TestHandleUnsubscribeMetrics(t *testing.T) {
	manager := NewManager(nil)
	e := echo.New()
	e.GET("/ws", manager.Handle)

	srv := httptest.NewServer(e)
	defer srv.Close()

	successBase := testutil.ToFloat64(wsUnsubscribeRequests.WithLabelValues(ChannelHead, "success"))
	unknownBase := testutil.ToFloat64(wsUnsubscribeRequests.WithLabelValues("unknown", "error"))

	conn := dialWS(t, srv)
	defer conn.Close()

	subscribe(t, conn, ChannelHead)
	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 1
	}, time.Second, 10*time.Millisecond)

	// valid unsubscribe -> success
	unsubscribe(t, conn, ChannelHead)
	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 0
	}, time.Second, 10*time.Millisecond, "client must be removed from the head channel")

	require.Eventually(t, func() bool {
		return testutil.ToFloat64(wsUnsubscribeRequests.WithLabelValues(ChannelHead, "success")) == successBase+1
	}, time.Second, 10*time.Millisecond, "successful unsubscribe must be counted with status=success")

	// unknown channel -> error, the client is notified and the connection stays open
	unsubscribe(t, conn, "unknown")
	var errMsg ErrorMessage
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(time.Second)))
	require.NoError(t, conn.ReadJSON(&errMsg))
	require.Equal(t, ChannelError, errMsg.Channel)
	require.Equal(t, ErrCodeUnknownChannel, errMsg.Body.Code)

	require.Eventually(t, func() bool {
		return testutil.ToFloat64(wsUnsubscribeRequests.WithLabelValues("unknown", "error")) == unknownBase+1
	}, time.Second, 10*time.Millisecond, "unsubscribe from unknown channel must be counted with status=error")
}

// TestHandleCleansUpSubscriptionsOnAbruptDisconnect verifies that an abrupt
// connection drop (no close handshake) terminates the read loop and releases the
// client from both the manager and the channels it was subscribed to. Before the
// fix such errors fell through the switch and the loop span until gorilla panicked
// with "repeated read on failed websocket connection".
func TestHandleCleansUpSubscriptionsOnAbruptDisconnect(t *testing.T) {
	manager := NewManager(nil)
	e := echo.New()
	e.GET("/ws", manager.Handle)

	srv := httptest.NewServer(e)
	defer srv.Close()

	conn := dialWS(t, srv)
	subscribe(t, conn, ChannelHead)

	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 1 && manager.clients.Len() == 1
	}, time.Second, 10*time.Millisecond, "client must be registered in the head channel")

	// drop the underlying tcp connection without a websocket close handshake
	require.NoError(t, conn.UnderlyingConn().Close())

	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 0 && manager.clients.Len() == 0
	}, 2*time.Second, 10*time.Millisecond, "subscriptions must be cleaned up after an abrupt disconnect")
}

// TestHandleDoesNotTouchUnsubscribedChannels verifies that disconnecting only
// releases the channels the client actually subscribed to: the subscription gauge
// of channels the client never subscribed to must stay untouched.
func TestHandleDoesNotTouchUnsubscribedChannels(t *testing.T) {
	manager := NewManager(nil)
	e := echo.New()
	e.GET("/ws", manager.Handle)

	srv := httptest.NewServer(e)
	defer srv.Close()

	// wsSubscriptions is a package-global gauge, capture baselines and assert deltas
	headBase := testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelHead))
	blocksBase := testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelBlocks))
	gasBase := testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelGasPrice))

	conn := dialWS(t, srv)
	subscribe(t, conn, ChannelHead)

	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 1
	}, time.Second, 10*time.Millisecond)

	require.Equal(t, 0, manager.blocks.clients.Len(), "client never subscribed to blocks")
	require.Equal(t, 0, manager.gasPrice.clients.Len(), "client never subscribed to gas price")

	require.NoError(t, conn.UnderlyingConn().Close())

	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 0
	}, 2*time.Second, 10*time.Millisecond, "head subscription must be released after disconnect")

	// the head gauge must return to baseline (inc on subscribe, dec on disconnect)
	require.Equal(t, headBase, testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelHead)))
	// gauges of channels the client never subscribed to must be untouched
	require.Equal(t, blocksBase, testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelBlocks)), "blocks subscription gauge must not drift")
	require.Equal(t, gasBase, testutil.ToFloat64(wsSubscriptions.WithLabelValues(ChannelGasPrice)), "gas price subscription gauge must not drift")
}

// TestHandleSurvivesApplicationLevelErrors verifies that malformed payloads and
// unknown methods are logged but keep the connection alive: a subsequent valid
// subscribe must still take effect.
func TestHandleSurvivesApplicationLevelErrors(t *testing.T) {
	manager := NewManager(nil)
	e := echo.New()
	e.GET("/ws", manager.Handle)

	srv := httptest.NewServer(e)
	defer srv.Close()

	conn := dialWS(t, srv)
	defer conn.Close()

	readError := func(wantCode int) {
		t.Helper()
		var errMsg ErrorMessage
		require.NoError(t, conn.SetReadDeadline(time.Now().Add(time.Second)))
		require.NoError(t, conn.ReadJSON(&errMsg))
		require.Equal(t, ChannelError, errMsg.Channel)
		require.Equal(t, wantCode, errMsg.Body.Code)
		require.NotEmpty(t, errMsg.Body.Message)
	}

	// malformed json: handle fails to decode, must not break the connection
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte("{not a json")))
	readError(ErrCodeInvalidMessage)

	// unknown method: handle returns an error, but the connection must stay open
	require.NoError(t, conn.WriteJSON(Message{Method: "unknown", Body: json.RawMessage(`{}`)}))
	readError(ErrCodeUnknownMethod)

	// unknown channel: handle returns an error, but the connection must stay open
	require.NoError(t, conn.WriteJSON(Message{Method: MethodSubscribe, Body: json.RawMessage(`{"channel":"unknown","filters":{}}`)}))
	readError(ErrCodeUnknownChannel)

	// the connection must still be usable
	subscribe(t, conn, ChannelHead)

	require.Eventually(t, func() bool {
		return manager.head.clients.Len() == 1
	}, time.Second, 10*time.Millisecond, "connection must survive application-level errors")
}

func TestIpsCheckAndSetLimit(t *testing.T) {
	ips := NewIps(3)

	for i := 0; i < 3; i++ {
		require.NoError(t, ips.CheckAndSet(testIp))
	}

	count, ok := ips.Get(testIp)
	require.True(t, ok)
	require.Equal(t, 3, count)

	err := ips.CheckAndSet(testIp)
	require.ErrorIs(t, err, ErrTooManyClients)

	count, ok = ips.Get(testIp)
	require.True(t, ok)
	require.Equal(t, 3, count, "rejected increment must not change counter")
}

func TestIpsDecrementAtLimitUnblocks(t *testing.T) {
	ips := NewIps(2)

	require.NoError(t, ips.CheckAndSet(testIp))
	require.NoError(t, ips.CheckAndSet(testIp))
	require.ErrorIs(t, ips.CheckAndSet(testIp), ErrTooManyClients)

	ips.Decrement(testIp)

	count, ok := ips.Get(testIp)
	require.True(t, ok)
	require.Equal(t, 1, count)

	require.NoError(t, ips.CheckAndSet(testIp), "ip must be able to connect again after disconnect")
}

func TestIpsDecrementDeletesEntryAtZero(t *testing.T) {
	ips := NewIps(2)

	require.NoError(t, ips.CheckAndSet(testIp))
	ips.Decrement(testIp)

	_, ok := ips.Get(testIp)
	require.False(t, ok, "entry must be removed from map when counter reaches zero")
}

func TestIpsDecrementUnknownIp(t *testing.T) {
	ips := NewIps(2)

	ips.Decrement(testIp)

	_, ok := ips.Get(testIp)
	require.False(t, ok, "decrement of unknown ip must not create an entry")
}

func TestIpsIndependentCounters(t *testing.T) {
	ips := NewIps(1)

	require.NoError(t, ips.CheckAndSet("10.0.0.1"))
	require.NoError(t, ips.CheckAndSet("10.0.0.2"), "limit is per ip, not global")
	require.ErrorIs(t, ips.CheckAndSet("10.0.0.1"), ErrTooManyClients)
}

func TestIpsConcurrentCheckAndSet(t *testing.T) {
	const (
		max     = 10
		workers = 100
	)
	ips := NewIps(max)

	var (
		wg            gosync.WaitGroup
		mx            gosync.Mutex
		acceptedCount int
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := ips.CheckAndSet(testIp); err == nil {
				mx.Lock()
				acceptedCount++
				mx.Unlock()
			}
		}()
	}
	wg.Wait()

	require.Equal(t, max, acceptedCount, "exactly max connections must be accepted")

	count, ok := ips.Get(testIp)
	require.True(t, ok)
	require.Equal(t, max, count)

	// concurrent disconnects release all slots and remove the entry
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ips.Decrement(testIp)
		}()
	}
	wg.Wait()

	_, ok = ips.Get(testIp)
	require.False(t, ok, "entry must be removed after all clients disconnect")
}

func TestNewManagerIpsLimit(t *testing.T) {
	manager := NewManager(nil, WithWebsocketClientsPerIp(1))
	require.NoError(t, manager.ips.CheckAndSet(testIp))
	require.ErrorIs(t, manager.ips.CheckAndSet(testIp), ErrTooManyClients)

	defaulted := NewManager(nil, WithWebsocketClientsPerIp(0))
	require.Equal(t, 10, defaulted.ips.max, "non-positive limit must fall back to default")
}

func TestHandleRejectionDoesNotDecrement(t *testing.T) {
	manager := NewManager(nil, WithWebsocketClientsPerIp(1))
	e := echo.New()

	require.NoError(t, manager.ips.CheckAndSet(testIp))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXRealIP, testIp)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := manager.Handle(c)
	require.ErrorIs(t, err, ErrTooManyClients)

	count, ok := manager.ips.Get(testIp)
	require.True(t, ok)
	require.Equal(t, 1, count, "rejected request must not decrement counter of accepted connections")
}

func TestHandleDecrementsAfterFailedUpgrade(t *testing.T) {
	manager := NewManager(nil, WithWebsocketClientsPerIp(2))
	e := echo.New()

	// plain GET without websocket headers: passes the ip limit check, fails on upgrade
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXRealIP, testIp)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := manager.Handle(c)
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTooManyClients)

	_, ok := manager.ips.Get(testIp)
	require.False(t, ok, "counter must be released after failed upgrade")
}

func TestHandleLimitDisconnectReconnect(t *testing.T) {
	manager := NewManager(nil, WithWebsocketClientsPerIp(2))
	e := echo.New()
	e.GET("/ws", manager.Handle)

	srv := httptest.NewServer(e)
	defer srv.Close()

	wsUrl := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	ip := "127.0.0.1"

	dial := func() (*websocket.Conn, error) {
		conn, resp, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return conn, err
	}

	first, err := dial()
	require.NoError(t, err)
	defer first.Close()

	second, err := dial()
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		count, ok := manager.ips.Get(ip)
		return ok && count == 2
	}, time.Second, 10*time.Millisecond)

	_, err = dial()
	require.Error(t, err, "third connection from the same ip must be rejected")

	count, ok := manager.ips.Get(ip)
	require.True(t, ok)
	require.Equal(t, 2, count, "rejected connection must not change counter")

	require.NoError(t, second.Close())

	require.Eventually(t, func() bool {
		count, ok := manager.ips.Get(ip)
		return ok && count == 1
	}, time.Second, 10*time.Millisecond, "counter must be decremented after disconnect")

	reconnect, err := dial()
	require.NoError(t, err, "ip must be able to connect again after disconnect")
	defer reconnect.Close()
}

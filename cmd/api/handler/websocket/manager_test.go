// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	gosync "sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const testIp = "10.0.0.1"

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

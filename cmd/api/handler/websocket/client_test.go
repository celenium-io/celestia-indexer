// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestNotifyClosedClient(t *testing.T) {
	client := newClient(10, nil, nil)
	err := client.Close()
	require.NoError(t, err, "closing client")
	client.Notify("test")
}

func BenchmarkHandle(b *testing.B) {
	e := echo.New()
	manager := NewManager(nil)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = manager.Handle(c)
	}

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	log.Println("Alloc", rtm.Alloc)
	log.Println("Frees", rtm.Frees)
	log.Println("Heap alloc", rtm.HeapAlloc)
	log.Println("Heap in use", rtm.HeapInuse)
	log.Println("last GC", rtm.LastGC)
}

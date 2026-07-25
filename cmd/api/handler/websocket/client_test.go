// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/gorilla/websocket"
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
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
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

// fakeNetError is a minimal net.Error implementation used to exercise the
// errors.As(err, &timeoutErr) branch of isExpectedDisconnect without
// depending on a real network timeout.
type fakeNetError struct {
	msg     string
	timeout bool
}

func (e *fakeNetError) Error() string   { return e.msg }
func (e *fakeNetError) Timeout() bool   { return e.timeout }
func (e *fakeNetError) Temporary() bool { return false }

func TestIsExpectedDisconnect(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "io.EOF",
			err:  io.EOF,
			want: true,
		},
		{
			name: "websocket close sent",
			err:  websocket.ErrCloseSent,
			want: true,
		},
		{
			name: "expected close code",
			err:  &websocket.CloseError{Code: websocket.CloseNormalClosure, Text: "bye"},
			want: true,
		},
		{
			name: "unexpected close code",
			err:  &websocket.CloseError{Code: websocket.CloseProtocolError, Text: "oops"},
			want: false,
		},
		{
			name: "timeout net.Error",
			err:  &fakeNetError{msg: "i/o timeout", timeout: true},
			want: true,
		},
		{
			name: "wrapped timeout net.Error",
			err:  fmt.Errorf("read message: %w", &fakeNetError{msg: "i/o timeout", timeout: true}),
			want: true,
		},
		{
			name: "non-timeout net.Error",
			err:  &fakeNetError{msg: "connection refused", timeout: false},
			want: false,
		},
		{
			name: "unrelated error does not panic",
			err:  errors.New("boom"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			require.NotPanics(t, func() {
				got = isExpectedDisconnect(tt.err)
			})
			require.Equal(t, tt.want, got)
		})
	}
}

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_observableCacheSkipper(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name   string
		path   string
		method string
		want   bool
	}{
		{
			name:   "test 1",
			path:   "/v1/ws",
			method: http.MethodGet,
			want:   true,
		}, {
			name:   "test 2",
			path:   "/metrics",
			method: http.MethodGet,
			want:   true,
		}, {
			// 	name:   "test 3",
			// 	path:   "/v1/block/:height",
			// 	method: http.MethodGet,
			// 	want:   true,
			// }, {
			// 	name:   "test 4",
			// 	path:   "/v1/tx/:hash",
			// 	method: http.MethodGet,
			// 	want:   true,
			// }, {
			name:   "test 5",
			path:   "/v1/some_post",
			method: http.MethodPost,
			want:   true,
		}, {
			name:   "test 6",
			path:   "/v1/valid",
			method: http.MethodGet,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			got := observableCacheSkipper(c)
			require.Equal(t, tt.want, got)
		})
	}
}

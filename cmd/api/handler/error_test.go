// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type fakeNoRows struct {
	isNoRows bool
}

func (f fakeNoRows) IsNoRows(err error) bool {
	return f.isNoRows
}

func TestHandleError(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		err        error
		noRows     NoRows
		wantStatus int
	}{
		{
			name:   "nil error",
			err:    nil,
			noRows: fakeNoRows{},
		},
		{
			name:       "deadline exceeded",
			err:        context.DeadlineExceeded,
			noRows:     fakeNoRows{},
			wantStatus: http.StatusRequestTimeout,
		},
		{
			name:       "wrapped deadline exceeded",
			err:        fmt.Errorf("query: %w", context.DeadlineExceeded),
			noRows:     fakeNoRows{},
			wantStatus: http.StatusRequestTimeout,
		},
		{
			name:       "context canceled",
			err:        context.Canceled,
			noRows:     fakeNoRows{},
			wantStatus: http.StatusBadGateway,
		},
		{
			name:       "no rows",
			err:        errors.New("not found"),
			noRows:     fakeNoRows{isNoRows: true},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid address",
			err:        errInvalidAddress,
			noRows:     fakeNoRows{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown address",
			err:        errUnknownAddress,
			noRows:     fakeNoRows{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unrelated error does not panic",
			err:        errors.New("boom"),
			noRows:     fakeNoRows{},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var handleErr error
			require.NotPanics(t, func() {
				handleErr = handleError(c, tt.err, tt.noRows)
			})
			require.NoError(t, handleErr)

			if tt.err == nil {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Empty(t, rec.Body.String())
				return
			}

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

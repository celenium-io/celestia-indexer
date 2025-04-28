// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/mock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_isAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "test 1",
			address: "celestia12y6fchaufs4tmn8e8wlk3rtrrftpqp6vr228a7",
			want:    true,
		}, {
			name:    "test 2",
			address: "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
			want:    false,
		}, {
			name:    "test 3",
			address: "invalid",
			want:    false,
		}, {
			name:    "test 4",
			address: "celestiavaloper111111111111111111111111111111111111111",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAddress(tt.address)
			require.Equal(t, tt.want, got, tt.name)
		})
	}
}

func Test_isValoperAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "test 1",
			address: "celestia12y6fchaufs4tmn8e8wlk3rtrrftpqp6vr228a7",
			want:    false,
		}, {
			name:    "test 2",
			address: "celestiavaloper1qycj0ymu9fqvwgyw4xz93p3n4a83jjk7sm2wzh",
			want:    true,
		}, {
			name:    "test 3",
			address: "invalid",
			want:    false,
		}, {
			name:    "test 4",
			address: "celestiavaloper111111111111111111111111111111111111111",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValoperAddress(tt.address)
			require.Equal(t, tt.want, got, tt.name)
		})
	}
}

func TestKeyValidator_Validate(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errChecker := mock.NewMockIDelegation(ctrl)
		apiKeys := mock.NewMockIApiKey(ctrl)
		kv := NewKeyValidator(apiKeys, errChecker)

		apiKeys.EXPECT().
			Get(gomock.Any(), "valid").
			Return(storage.ApiKey{
				Key:         "valid",
				Description: "descr",
			}, nil).
			Times(1)

		ok, err := kv.Validate("valid", ctx)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("invalid key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errChecker := mock.NewMockIDelegation(ctrl)
		apiKeys := mock.NewMockIApiKey(ctrl)
		kv := NewKeyValidator(apiKeys, errChecker)

		apiKeys.EXPECT().
			Get(gomock.Any(), "invalid").
			Return(storage.ApiKey{}, sql.ErrNoRows).
			Times(1)

		errChecker.EXPECT().
			IsNoRows(sql.ErrNoRows).
			Return(true).
			Times(1)

		ok, err := kv.Validate("invalid", ctx)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("unexpected error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(req, rec)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errChecker := mock.NewMockIDelegation(ctrl)
		apiKeys := mock.NewMockIApiKey(ctrl)
		kv := NewKeyValidator(apiKeys, errChecker)

		unexpectedErr := errors.New("unexpected")

		apiKeys.EXPECT().
			Get(gomock.Any(), "invalid").
			Return(storage.ApiKey{}, unexpectedErr).
			Times(1)

		errChecker.EXPECT().
			IsNoRows(unexpectedErr).
			Return(false).
			Times(1)

		ok, err := kv.Validate("invalid", ctx)
		require.Error(t, err)
		require.False(t, ok)
	})
}

func TestValidateRollupProvider(t *testing.T) {
	t.Run("all fields are filled", func(t *testing.T) {
		req := createRollupRequest{
			Name:        "name",
			Description: "description",
			Logo:        "https://celenium.io/",
			Providers: []rollupProvider{
				{
					Namespace: "AAAAAAAAAAAAAAAAAAAAAAAAALt7GEYP9l+FgiU=",
					Address:   "celestia1q58cnwjk6mftzh48hw76wdf27zs5vf5mys9ujz",
				},
			},
		}

		v := NewCelestiaApiValidator()
		err := v.Validate(req)
		require.NoError(t, err)
	})

	t.Run("only address is filled", func(t *testing.T) {
		req := createRollupRequest{
			Name:        "name",
			Description: "description",
			Logo:        "https://celenium.io/",
			Providers: []rollupProvider{
				{
					Address: "celestia1q58cnwjk6mftzh48hw76wdf27zs5vf5mys9ujz",
				},
			},
		}

		v := NewCelestiaApiValidator()
		err := v.Validate(req)
		require.NoError(t, err)
	})

	t.Run("only namespace is filled", func(t *testing.T) {
		req := createRollupRequest{
			Name:        "name",
			Description: "description",
			Logo:        "https://celenium.io/",
			Providers: []rollupProvider{
				{
					Namespace: "AAAAAAAAAAAAAAAAAAAAAAAAALt7GEYP9l+FgiU=",
				},
			},
		}

		v := NewCelestiaApiValidator()
		err := v.Validate(req)
		require.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		req := createRollupRequest{
			Name:        "name",
			Description: "description",
			Logo:        "https://celenium.io/",
			Providers: []rollupProvider{
				{},
			},
		}

		v := NewCelestiaApiValidator()
		err := v.Validate(req)
		require.Error(t, err)
	})
}

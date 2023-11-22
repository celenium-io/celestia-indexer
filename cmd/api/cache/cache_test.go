// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache_SetGet(t *testing.T) {
	t.Run("set and get key from cache", func(t *testing.T) {
		c := NewCache(Config{MaxEntitiesCount: 2})
		c.Set("test", []byte{0, 1, 2, 3})

		got, ok := c.Get("test")
		require.True(t, ok)
		require.Equal(t, []byte{0, 1, 2, 3}, got)

		_, exists := c.Get("unknown")
		require.False(t, exists)
	})

	t.Run("overflow set queue", func(t *testing.T) {
		c := NewCache(Config{MaxEntitiesCount: 2})
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}

		require.Len(t, c.queue, 2)
		require.Len(t, c.m, 2)

		got, ok := c.Get("99")
		require.True(t, ok)
		require.Equal(t, []byte{99}, got)

		got1, ok1 := c.Get("98")
		require.True(t, ok1)
		require.Equal(t, []byte{98}, got1)

		_, exists := c.Get("0")
		require.False(t, exists)
	})
}

func TestCache_Clear(t *testing.T) {
	t.Run("set and get key from cache", func(t *testing.T) {
		c := NewCache(Config{MaxEntitiesCount: 100})
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}
		c.Clear()

		require.Len(t, c.queue, 100)
		for i := 0; i < 100; i++ {
			require.EqualValues(t, c.queue[i], "")
		}
		require.Len(t, c.m, 0)
	})
}

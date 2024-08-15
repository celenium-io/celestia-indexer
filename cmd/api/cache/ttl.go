// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
)

type TTLCache struct {
	db         *badger.DB
	expiration time.Duration
}

func NewTTLCache(expiration time.Duration) (*TTLCache, error) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return nil, err
	}

	return &TTLCache{
		db:         db,
		expiration: expiration,
	}, nil
}

func (c *TTLCache) Close() error {
	return c.db.Close()
}

func (c *TTLCache) Get(key string) (data []byte, found bool) {
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		if err := item.Value(func(val []byte) error {
			data = make([]byte, len(val))
			copy(data, val)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, false
	}
	return data, true
}

func (c *TTLCache) Set(key string, data []byte) {
	keyBytes := []byte(key)
	err := c.db.Update(func(txn *badger.Txn) error {
		e := badger.
			NewEntry(keyBytes, data).
			WithTTL(c.expiration)
		return txn.SetEntry(e)
	})
	if err != nil {
		log.Err(err).Msgf("set %s to TTL cache", key)
	}
}

func (c *TTLCache) Clear() {
	if err := c.db.DropAll(); err != nil {
		log.Err(err).Msg("clear ttl cache")
	}
}

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
)

type TTLCache struct {
	db              *badger.DB
	expiration      time.Duration
	gcCollectPeriod time.Duration
	wg              *sync.WaitGroup
}

const dir = "/tmp/badger"

func NewTTLCache(expiration time.Duration, inMemory bool) (*TTLCache, error) {
	path := ""
	if !inMemory {
		path = dir
	}

	db, err := badger.Open(badger.DefaultOptions(path).WithInMemory(inMemory))
	if err != nil {
		return nil, err
	}

	return &TTLCache{
		db:              db,
		expiration:      expiration,
		gcCollectPeriod: time.Minute * 5,
		wg:              new(sync.WaitGroup),
	}, nil
}

func (c *TTLCache) Start(ctx context.Context) {
	c.wg.Add(1)
	go c.gcCollect(ctx)
}

func (c *TTLCache) gcCollect(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.gcCollectPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.db.RunValueLogGC(0.5); err != nil {
				log.Err(err).Msg("ttl cache garbage collection error")
			}
		}
	}
}

func (c *TTLCache) Close() error {
	c.wg.Wait()
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

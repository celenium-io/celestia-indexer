// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"context"
	"sync"

	"github.com/celenium-io/celestia-indexer/cmd/api/bus"
	"github.com/dipdup-io/workerpool"
)

type ObservableCache struct {
	maxEntitiesCount int
	observer         *bus.Observer

	m     map[string][]byte
	queue []string
	mx    *sync.RWMutex
	g     workerpool.Group
}

func NewObservableCache(cfg Config, observer *bus.Observer) *ObservableCache {
	return &ObservableCache{
		maxEntitiesCount: cfg.MaxEntitiesCount,
		observer:         observer,
		m:                make(map[string][]byte),
		queue:            make([]string, cfg.MaxEntitiesCount),
		mx:               new(sync.RWMutex),
		g:                workerpool.NewGroup(),
	}
}

func (c *ObservableCache) Start(ctx context.Context) {
	c.g.GoCtx(ctx, c.listen)
}

func (c *ObservableCache) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.observer.Head():
			c.Clear()
		}
	}
}

func (c *ObservableCache) Close() error {
	c.g.Wait()
	return nil
}

func (c *ObservableCache) Get(key string) ([]byte, bool) {
	c.mx.RLock()
	data, ok := c.m[key]
	c.mx.RUnlock()
	return data, ok
}

func (c *ObservableCache) Set(key string, data []byte) {
	c.mx.Lock()
	queueIdx := len(c.m)

	if _, ok := c.m[key]; ok {
		c.m[key] = data
	} else {
		c.m[key] = data
		if queueIdx == c.maxEntitiesCount {
			keyForRemove := c.queue[queueIdx-1]
			c.queue = append([]string{key}, c.queue[:queueIdx-1]...)
			delete(c.m, keyForRemove)
		} else {
			c.queue[c.maxEntitiesCount-queueIdx-1] = key
		}
	}
	c.mx.Unlock()
}

func (c *ObservableCache) Clear() {
	c.mx.Lock()
	for key := range c.m {
		delete(c.m, key)
	}
	c.queue = make([]string, c.maxEntitiesCount)
	c.mx.Unlock()
}

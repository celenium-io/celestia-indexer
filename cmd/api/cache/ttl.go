// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"sync"
	"time"
)

type ttlItem struct {
	data      []byte
	expiredAt time.Time
}

type TTLCache struct {
	maxEntitiesCount int
	m                map[string]*ttlItem
	queue            []string
	mx               *sync.RWMutex
	expiration       time.Duration
}

func NewTTLCache(cfg Config, expiration time.Duration) *TTLCache {
	return &TTLCache{
		maxEntitiesCount: cfg.MaxEntitiesCount,
		m:                make(map[string]*ttlItem),
		queue:            make([]string, cfg.MaxEntitiesCount),
		mx:               new(sync.RWMutex),
		expiration:       expiration,
	}
}

func (c *TTLCache) Get(key string) ([]byte, bool) {
	c.mx.RLock()
	item, ok := c.m[key]
	c.mx.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiredAt) {
		c.mx.Lock()
		{
			copying := false
			if len(c.queue) > len(c.m) {
				for i := len(c.queue) - 1; i > 0; i-- {
					if copying = copying || c.queue[i] == key; copying {
						c.queue[i] = c.queue[i-1]
					}
				}
				c.queue[0] = ""
			} else {
				for i := 0; i < len(c.queue)-1; i++ {
					if copying = copying || c.queue[i] == key; copying {
						c.queue[i] = c.queue[i+1]
					}
				}
				c.queue[len(c.queue)-1] = ""
			}
			delete(c.m, key)
		}
		c.mx.Unlock()
		return nil, false
	}
	return item.data, true
}

func (c *TTLCache) Set(key string, data []byte) {
	c.mx.Lock()
	queueIdx := len(c.m)
	item := &ttlItem{
		data:      data,
		expiredAt: time.Now().Add(c.expiration),
	}

	if _, ok := c.m[key]; ok {
		c.m[key] = item
	} else {
		c.m[key] = item
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

func (c *TTLCache) Clear() {
	c.mx.Lock()
	for key := range c.m {
		delete(c.m, key)
	}
	c.queue = make([]string, c.maxEntitiesCount)
	c.mx.Unlock()
}

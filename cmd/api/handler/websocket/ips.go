// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import "sync"

type Ips struct {
	ips map[string]int
	mx  sync.RWMutex
	max int
}

func NewIps(max int) *Ips {
	return &Ips{
		ips: make(map[string]int),
		max: max,
	}
}

func (ips *Ips) Get(ip string) (int, bool) {
	ips.mx.RLock()
	defer ips.mx.RUnlock()
	if count, ok := ips.ips[ip]; ok {
		return count, true
	}
	return 0, false
}

func (ips *Ips) CheckAndSet(ip string) error {
	ips.mx.Lock()
	defer ips.mx.Unlock()
	if count, ok := ips.ips[ip]; ok {
		if count >= ips.max {
			return ErrTooManyClients
		}
		ips.ips[ip] = count + 1
	} else {
		ips.ips[ip] = 1
	}
	return nil
}

func (ips *Ips) Decrement(ip string) {
	ips.mx.Lock()
	defer ips.mx.Unlock()
	if count, ok := ips.ips[ip]; ok {
		if count > 1 {
			ips.ips[ip] = count - 1
		} else {
			delete(ips.ips, ip)
		}
	}
}

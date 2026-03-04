// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Notificator struct {
	pool *pgxpool.Pool
	ch   chan pgconn.Notification
	once sync.Once
	done chan struct{}
	wg   sync.WaitGroup
}

func NewNotificator(pool *pgxpool.Pool) *Notificator {
	return &Notificator{
		pool: pool,
		ch:   make(chan pgconn.Notification, 16),
		done: make(chan struct{}),
	}
}

func (n *Notificator) Notify(ctx context.Context, channel string, payload string) error {
	conn, err := n.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "acquire connection for notify")
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, "SELECT pg_notify($1, $2)", channel, payload)
	return err
}

func (n *Notificator) Subscribe(ctx context.Context, channels ...string) error {
	conn, err := n.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "acquire connection for listen")
	}
	for _, ch := range channels {
		if _, err := conn.Exec(ctx, "LISTEN "+ch); err != nil {
			conn.Release()
			return errors.Wrapf(err, "listen %s", ch)
		}
	}
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer conn.Release()
		for {
			notification, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				return
			}
			select {
			case n.ch <- *notification:
			case <-n.done:
				return
			}
		}
	}()
	return nil
}

func (n *Notificator) Listen() <-chan pgconn.Notification {
	return n.ch
}

func (n *Notificator) Close() error {
	n.once.Do(func() { close(n.done) })
	n.wg.Wait()
	return nil
}

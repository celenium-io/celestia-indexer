// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

const (
	connectionName       = "celestia_notifications"
	minReconnectInterval = 10 * time.Second
	maxReconnectInterval = time.Minute
)

type Notificator struct {
	db *bun.DB
	l  *pq.Listener
}

func NewNotificator(cfg config.Database, db *bun.DB) *Notificator {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	return &Notificator{
		l: pq.NewListener(
			connStr,
			minReconnectInterval,
			maxReconnectInterval,
			nil,
		),
		db: db,
	}
}

func (n *Notificator) Notify(ctx context.Context, channel string, payload string) error {
	_, err := n.db.ExecContext(ctx, "NOTIFY ?, ?", bun.Ident(channel), payload)
	return err
}

func (n *Notificator) Listen() chan *pq.Notification {
	return n.l.Notify
}

func (n *Notificator) Subscribe(ctx context.Context, channels ...string) error {
	for i := range channels {
		if err := n.l.Listen(channels[i]); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notificator) Close() error {
	return n.l.Close()
}

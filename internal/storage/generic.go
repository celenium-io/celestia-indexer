package storage

import (
	"context"
	"io"

	"github.com/lib/pq"
)

var Models = []any{
	&State{},
	&Address{},
	&Block{},
	&Tx{},
	&Message{},
	&Event{},
	&Namespace{},
	&NamespaceMessage{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	io.Closer

	Subscribe(ctx context.Context, channels ...string) error
	Listen() chan *pq.Notification
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ListenerFactory interface {
	CreateListener() Listener
}

const (
	ChannelHead = "head"
	ChannelTx   = "tx"
)

package websocket

import (
	"context"

	sdkSync "github.com/dipdup-net/indexer-sdk/pkg/sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/workerpool"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type identifiable[M any] interface {
	GetById(ctx context.Context, id uint64) (M, error)
}

type processor[I, M any] func(ctx context.Context, payload string, repo identifiable[I]) (M, error)

type Channel[I, M any] struct {
	storageChannelName string
	clients            *sdkSync.Map[uint64, client]
	listener           storage.Listener
	log                zerolog.Logger
	processor          processor[I, M]
	filters            Filterable[M]
	repo               identifiable[I]

	g workerpool.Group
}

func NewChannel[I, M any](storageChannelName string, processor processor[I, M], repo identifiable[I], filters Filterable[M]) *Channel[I, M] {
	return &Channel[I, M]{
		storageChannelName: storageChannelName,
		clients:            sdkSync.NewMap[uint64, client](),
		processor:          processor,
		filters:            filters,
		repo:               repo,
		log:                log.With().Str("channel", storageChannelName).Logger(),
		g:                  workerpool.NewGroup(),
	}
}

func (channel *Channel[I, M]) AddClient(c client) {
	channel.clients.Set(c.Id(), c)
}

func (channel *Channel[I, M]) RemoveClient(id uint64) {
	channel.clients.Delete(id)
}

func (channel *Channel[I, M]) String() string {
	return channel.storageChannelName
}

func (channel *Channel[I, M]) Start(ctx context.Context, factory storage.ListenerFactory) {
	if channel.processor == nil {
		channel.log.Panic().Msg("nil processor in channel")
		return
	}
	if channel.filters == nil {
		channel.log.Panic().Msg("nil filters in channel")
		return
	}
	if factory == nil {
		channel.log.Panic().Msg("nil listener factory in channel")
		return
	}

	channel.listener = factory.CreateListener()

	if err := channel.listener.Subscribe(ctx, channel.storageChannelName); err != nil {
		channel.log.Panic().Err(err).Msg("subscribe on storage channel")
		return
	}

	channel.g.GoCtx(ctx, channel.waitMessage)
}

func (channel *Channel[I, M]) waitMessage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-channel.listener.Listen():
			if !ok {
				return
			}
			if msg.Channel != channel.storageChannelName {
				channel.log.Error().
					Str("msg", msg.Channel).
					Msg("unexpected channel message")
				continue
			}

			if channel.clients.Len() == 0 {
				continue
			}

			if err := channel.processMessage(ctx, msg); err != nil {
				log.Err(err).
					Str("msg", msg.Channel).
					Str("payload", msg.Extra).
					Msg("processing channel message")
			}

		}
	}
}

func (channel *Channel[I, M]) processMessage(ctx context.Context, msg *pq.Notification) error {
	log.Trace().
		Str("channel", msg.Channel).
		Str("payload", msg.Extra).
		Msg("message received")

	notification, err := channel.processor(ctx, msg.Extra, channel.repo)
	if err != nil {
		return errors.Wrap(err, "processing channel message")
	}

	if err := channel.clients.Range(func(_ uint64, value client) (error, bool) {
		if channel.filters.Filter(value, notification) {
			value.Notify(notification)
		}
		return nil, false
	}); err != nil {
		return errors.Wrap(err, "write message to client")
	}

	return nil
}

func (channel *Channel[I, M]) Close() error {
	channel.g.Wait()

	if channel.listener != nil {
		if err := channel.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}

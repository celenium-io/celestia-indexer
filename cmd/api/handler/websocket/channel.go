package websocket

import (
	"context"
	"sync"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type identifiable[M any] interface {
	GetById(ctx context.Context, id uint64) (M, error)
}

type processor[I, M any] func(ctx context.Context, payload string, repo identifiable[I]) (M, error)

type Channel[I, M any] struct {
	storageChannelName string
	clients            Map[uint64, *Client]
	listener           storage.Listener
	log                zerolog.Logger
	processor          processor[I, M]
	filters            Filterable[M]
	repo               identifiable[I]

	wg *sync.WaitGroup
}

func NewChannel[I, M any](storageChannelName string, processor processor[I, M], repo identifiable[I], filters Filterable[M]) *Channel[I, M] {
	return &Channel[I, M]{
		storageChannelName: storageChannelName,
		clients:            NewMap[uint64, *Client](),
		processor:          processor,
		filters:            filters,
		repo:               repo,
		log:                log.With().Str("channel", storageChannelName).Logger(),
		wg:                 new(sync.WaitGroup),
	}
}

func (channel *Channel[I, M]) AddClient(c *Client) {
	channel.clients.Set(c.id, c)
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
	channel.wg.Add(1)
	go channel.waitMessage(ctx)
}

func (channel *Channel[I, M]) waitMessage(ctx context.Context) {
	defer channel.wg.Done()

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

			log.Debug().
				Str("channel", msg.Channel).
				Str("payload", msg.Extra).
				Msg("message received")

			notification, err := channel.processor(ctx, msg.Extra, channel.repo)
			if err != nil {
				channel.log.Err(err).Msg("processing channel message")
				continue
			}

			if err := channel.clients.Range(func(_ uint64, value *Client) (error, bool) {
				if channel.filters.Filter(value, notification) {
					value.Notify(notification)
				}
				return nil, false
			}); err != nil {
				channel.log.Err(err).Msg("write message to client")
			}
		}
	}
}

func (channel *Channel[I, M]) Close() error {
	channel.wg.Wait()

	if channel.listener != nil {
		if err := channel.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}

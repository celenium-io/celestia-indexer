// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package receiver

import (
	"context"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	tendermint "github.com/cometbft/cometbft/types"
	"github.com/pkg/errors"
)

func (r *Module) sync(ctx context.Context) {
	var blocksCtx context.Context
	blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)

	if err := r.readBlocks(blocksCtx); err != nil {
		r.Log.Err(err).Msg("while reading blocks")
		r.stopAll()
		return
	}

	if ctx.Err() != nil {
		return
	}

	if r.ws != nil {
		if err := r.live(blocksCtx); err != nil {
			r.Log.Err(err).Msg("while reading blocks")
			r.stopAll()
			return
		}
	} else {
		ticker := time.NewTicker(time.Second * time.Duration(r.cfg.BlockPeriod))
		defer ticker.Stop()

		for {
			r.rollbackSync.Wait()

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blocksCtx, r.cancelReadBlocks = context.WithCancel(ctx)
				if err := r.readBlocks(blocksCtx); err != nil && !errors.Is(err, context.Canceled) {
					r.Log.Err(err).Msg("while reading blocks by timer")
					r.stopAll()
					return
				}
			}
		}
	}
}

func (r *Module) live(ctx context.Context) error {
	if err := r.ws.Start(); err != nil {
		return err
	}
	r.Log.Info().Msg("websocket started")

	ch, err := r.ws.Subscribe(ctx, "test", "tm.event = 'NewBlockHeader'")
	if err != nil {
		return err
	}
	r.Log.Info().Msg("websocket was subscribed on block header events")

	for {
		r.rollbackSync.Wait()

		select {
		case <-ctx.Done():
			return nil
		case block := <-ch:
			if block.Data == nil {
				continue
			}
			blockHeader, ok := block.Data.(tendermint.EventDataNewBlockHeader)
			if !ok {
				r.Log.Error().Msgf("unexpected block type: %T", block.Data)
				continue
			}
			r.Log.Info().Int64("height", blockHeader.Header.Height).Msg("new block received")
			r.passBlocks(ctx, types.Level(blockHeader.Header.Height))
		}
	}
}

func (r *Module) readBlocks(ctx context.Context) error {
	for {
		headLevel, err := r.api.CurrentHead(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		if level, _ := r.Level(); level == headLevel {
			time.Sleep(time.Second)
			continue
		}

		r.passBlocks(ctx, headLevel)
		return nil
	}
}

func (r *Module) passBlocks(ctx context.Context, head types.Level) {
	fetchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		batch    []types.Level
		maxLevel types.Level
	)

	for level := r.receivedLevel + 1; level <= head; level++ {
		for r.queueBlock.Load() {
			select {
			case <-fetchCtx.Done():
				r.fetchWg.Wait()
				return
			case <-time.After(100 * time.Millisecond):
			}
		}

		batch = append(batch, level)
		bulkSize := int(r.bulkSize.Load())
		if len(batch) >= bulkSize || level == head {
			levels := batch
			batch = nil

			last := levels[len(levels)-1]
			if last > maxLevel {
				maxLevel = last
			}

			// Acquire semaphore slot before spawning — blocks when at concurrency limit.
			select {
			case r.fetchSem <- struct{}{}:
			case <-fetchCtx.Done():
				r.fetchWg.Wait()
				return
			}

			r.fetchWg.Add(1)
			go func(lvls []types.Level) {
				defer r.fetchWg.Done()
				defer func() { <-r.fetchSem }()
				r.fetchBatch(fetchCtx, lvls)
			}(levels)
		}
	}

	r.fetchWg.Wait()
	if fetchCtx.Err() == nil && maxLevel > r.receivedLevel {
		r.receivedLevel = maxLevel
	}
}

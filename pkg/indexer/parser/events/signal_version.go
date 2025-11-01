// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func processSignalVersion(ctx *context.Context, _ []storage.Event, msg *storage.Message, data map[string]any, idx *int) error {
	version, err := decoder.Uint64(data, "Version")
	if err != nil {
		return errors.Wrap(err, "get signal version in exec")
	}

	val := storage.EmptyValidator()
	val.Address = decoder.StringFromMap(data, "ValidatorAddress")
	val.Version = version

	msg.SignalVersion = &storage.SignalVersion{
		Height:    msg.Height,
		Time:      msg.Time,
		Version:   version,
		Validator: &val,
	}
	ctx.AddValidator(*msg.SignalVersion.Validator)
	ctx.AddUpgrade(storage.Upgrade{
		Version:      version,
		SignalsCount: 1,
		Height:       msg.Height,
		Time:         msg.Time,
	})
	*idx += 1
	return nil
}

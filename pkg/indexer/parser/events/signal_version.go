// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func processSignalVersion(_ *context.Context, events []storage.Event, msg *storage.Message, data map[string]any, idx *int) error {
	version, err := decoder.Uint64(data, "Version")
	if err != nil {
		return errors.Wrap(err, "get signal version in exec")
	}
	validator := decoder.StringFromMap(data, "ValidatorAddress")
	msg.SignalVersion = &storage.SignalVersion{
		Height:  msg.Height,
		Time:    msg.Time,
		Version: version,
		Validator: &storage.Validator{
			Address: validator,
		},
	}
	toTheNextAction(events, idx)
	return nil
}

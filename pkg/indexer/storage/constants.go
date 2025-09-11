package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

func (module *Module) saveConstantUpdates(
	ctx context.Context,
	tx storage.Transaction,
	consts *sync.Map[string, *storage.Constant],
) error {

	newConstants := make([]storage.Constant, 0)
	err := consts.Range(func(key string, value *storage.Constant) (error, bool) {
		switch value.Name {
		case "evidence_max_age_num_blocks":
			if value.Value != module.maxAgeNumBlocks {
				newConstants = append(newConstants, *value)
				module.maxAgeNumBlocks = value.Value
			}
		case "evidence_max_age_duration":
			if value.Value != module.maxAgeDuration {
				newConstants = append(newConstants, *value)
				module.maxAgeDuration = value.Value
			}
		}
		return nil, false
	})

	if err != nil {
		return errors.Wrap(err, "range")
	}

	if len(newConstants) == 0 {
		return nil
	}

	return tx.SaveConstants(ctx, newConstants...)
}

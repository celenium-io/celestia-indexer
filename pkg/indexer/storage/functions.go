// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"strconv"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/pkg/errors"
)

func setNamespacesFromMessage(msg storage.Message, namespaces map[string]*storage.Namespace) {
	for i := range msg.Namespace {
		key := msg.Namespace[i].String()
		if ns, ok := namespaces[key]; !ok {
			msg.Namespace[i].PfbCount = 1
			namespaces[key] = &msg.Namespace[i]
		} else {
			ns.PfbCount += 1
			ns.Size += msg.Namespace[i].Size
		}
	}
}

func getMaxValidatorsCount(ctx context.Context, constants storage.IConstant) (int, error) {
	maxValsConsts, err := constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return 0, errors.Wrap(err, "get max validators value")
	}
	return strconv.Atoi(maxValsConsts.Value)
}

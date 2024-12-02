// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package rpc

import (
	"context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"strconv"
)

const pathBlobProofs = "prove_shares_v2"

func (api *API) BlobProofs(ctx context.Context, level pkgTypes.Level, startShare, endShare int) (pkgTypes.BlobProof, error) {
	if startShare < 0 || endShare < 0 {
		return pkgTypes.BlobProof{}, errors.New("params 'startShare' and 'endShare' should not be lower than 0")
	}
	if level <= 0 {
		return pkgTypes.BlobProof{}, errors.New("param 'level' should be greater than 0")
	}

	args := make(map[string]string)
	args["height"] = strconv.FormatInt(int64(level), 10)
	args["startShare"] = strconv.FormatInt(int64(startShare), 10)
	args["endShare"] = strconv.FormatInt(int64(endShare), 10)

	var proof types.Response[pkgTypes.BlobProof]
	if err := api.get(ctx, pathBlobProofs, args, &proof); err != nil {
		return pkgTypes.BlobProof{}, errors.Wrap(err, "api.get")
	}

	if proof.Error != nil {
		return pkgTypes.BlobProof{}, errors.Wrapf(types.ErrRequest, "request %d error: %s", proof.Id, proof.Error.Error())
	}

	return proof.Result, nil
}

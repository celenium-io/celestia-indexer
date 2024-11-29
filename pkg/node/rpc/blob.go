package rpc

import (
	"context"
	"github.com/celenium-io/celestia-indexer/pkg/node/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/pkg/errors"
	"strconv"
)

const pathBlockProofs = "prove_shares_v2"

func (api *API) BlobProofs(ctx context.Context, level pkgTypes.Level, startShare, endShare int) (pkgTypes.BlobProof, error) {
	args := make(map[string]string)
	if level != 0 {
		args["height"] = strconv.FormatInt(int64(level), 10)
		args["startShare"] = strconv.FormatInt(int64(startShare), 10)
		args["endShare"] = strconv.FormatInt(int64(endShare), 10)
	}

	var proof types.Response[pkgTypes.BlobProof]
	if err := api.get(ctx, pathBlockProofs, args, &proof); err != nil {
		return pkgTypes.BlobProof{}, errors.Wrap(err, "api.get")
	}

	if proof.Error != nil {
		return pkgTypes.BlobProof{}, errors.Wrapf(types.ErrRequest, "request %d error: %s", proof.Id, proof.Error.Error())
	}

	return proof.Result, nil
}

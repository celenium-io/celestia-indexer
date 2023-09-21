package dal

import (
	"context"

	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/pkg/errors"
)

// Blobs - returns all blobs under the given namespaces and height.
func (node *Node) Blobs(ctx context.Context, height uint64, namespaces ...string) ([]types.Blob, error) {
	if len(namespaces) == 0 {
		return nil, nil
	}

	var response types.Response[[]types.Blob]
	if err := node.post(ctx, "blob.GetAll", []any{height, namespaces}, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.Wrapf(types.ErrRequest, "request %d error: %s", response.Id, response.Error.Error())
	}
	return response.Result, nil
}

// Blob - retrieves the blob by commitment under the given namespace and height.
func (node *Node) Blob(ctx context.Context, height uint64, namespace, commitment string) (types.Blob, error) {
	var response types.Response[types.Blob]
	if err := node.post(ctx, "blob.Get", []any{height, namespace, commitment}, &response); err != nil {
		return response.Result, err
	}

	if response.Error != nil {
		return response.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", response.Id, response.Error.Error())
	}
	return response.Result, nil
}

// Proofs - retrieves proofs in the given namespaces at the given height by commitment.
func (node *Node) Proofs(ctx context.Context, height uint64, namespace, commitment string) ([]types.Proof, error) {
	var response types.Response[[]types.Proof]
	if err := node.post(ctx, "blob.GetProof", []any{height, namespace, commitment}, &response); err != nil {
		return response.Result, err
	}

	if response.Error != nil {
		return response.Result, errors.Wrapf(types.ErrRequest, "request %d error: %s", response.Id, response.Error.Error())
	}
	return response.Result, nil
}

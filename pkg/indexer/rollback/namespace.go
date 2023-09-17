package rollback

import (
	"context"
	"encoding/base64"
	"encoding/hex"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/postgres"
	"github.com/pkg/errors"
)

var errInvalidPayForBlob = errors.New("invalid MsgPayForBlob content")

func (module *Module) rollbackNamespaces(
	ctx context.Context,
	tx postgres.Transaction,
	nsMsgs []storage.NamespaceMessage,
	deletedNs []storage.Namespace,
	deletedMsgs []storage.Message,
) error {
	if len(nsMsgs) == 0 {
		return nil
	}
	deleted := make(map[uint64]struct{}, len(deletedNs))
	for i := range deletedNs {
		deleted[deletedNs[i].Id] = struct{}{}
	}
	deletedMessages := make(map[uint64]storage.Message, len(deletedMsgs))
	for i := range deletedMsgs {
		deletedMessages[deletedMsgs[i].Id] = deletedMsgs[i]
	}

	diffs := make(map[uint64]*storage.Namespace)
	for i := range nsMsgs {
		nsId := nsMsgs[i].NamespaceId
		if _, ok := deleted[nsId]; ok {
			continue
		}
		msgId := nsMsgs[i].MsgId
		msg, ok := deletedMessages[msgId]
		if !ok {
			return errors.Errorf("unknown message: %d", msgId)
		}
		nsSize, err := newNamespaceSize(msg.Data)
		if err != nil {
			return err
		}

		ns, err := tx.Namespace(ctx, nsId)
		if err != nil {
			return err
		}

		size, ok := nsSize[ns.String()]
		if !ok {
			return errors.Errorf("message does not contain info about namespace: ns_id=%d msg_id=%d", nsId, msgId)
		}

		if diff, ok := diffs[nsId]; ok {
			diff.PfbCount -= 1
			diff.Size -= size
		} else {
			ns.PfbCount -= 1
			ns.Size -= size
			diffs[nsId] = &ns
		}
	}

	namespaces := make([]*storage.Namespace, 0, len(diffs))
	for key := range diffs {
		namespaces = append(namespaces, diffs[key])
	}

	return tx.SaveNamespaces(ctx, namespaces...)
}

type namespaceSize map[string]uint64

func newNamespaceSize(data map[string]any) (namespaceSize, error) {
	sizesRaw, ok := data["blob_sizes"]
	if !ok {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}
	sizesAny, ok := sizesRaw.([]any)
	if !ok {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}
	if len(sizesAny) == 0 {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}

	nsRaw, ok := data["namespaces"]
	if !ok {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}
	nsAny, ok := nsRaw.([]any)
	if !ok {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}
	if len(nsAny) != len(sizesAny) {
		return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
	}

	size := make(namespaceSize)
	for i := range nsAny {
		nsString, ok := nsAny[i].(string)
		if !ok {
			return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
		}
		nsSize, ok := sizesAny[i].(int)
		if !ok {
			return nil, errors.Wrapf(errInvalidPayForBlob, "%##v", data)
		}
		data, err := base64.StdEncoding.DecodeString(nsString)
		if err != nil {
			return nil, errors.Wrap(err, nsString)
		}

		size[hex.EncodeToString(data)] = uint64(nsSize)
	}

	return size, nil
}

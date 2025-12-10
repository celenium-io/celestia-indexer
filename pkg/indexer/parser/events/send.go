package events

import (
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/context"
	"github.com/celenium-io/celestia-indexer/pkg/indexer/decode/decoder"
	"github.com/pkg/errors"
)

func handleSend(ctx *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	if idx == nil {
		return errors.New("nil event index")
	}
	if msg == nil {
		return errors.New("nil message in events handler")
	}
	action := decoder.StringFromMap(events[*idx].Data, "action")
	isValidMsg := action == "/cosmos.bank.v1beta1.MsgSend"
	if !isValidMsg {
		return errors.Errorf("unexpected event action %s for message type %s", action, msg.Type.String())
	}
	*idx += 1
	return processSend(ctx, events, msg, idx)
}

func processSend(_ *context.Context, events []storage.Event, msg *storage.Message, idx *int) error {
	msgIdx := decoder.StringFromMap(events[*idx].Data, "msg_index")
	newFormat := msgIdx != ""

	if newFormat {
		*idx += 4
	} else {
		*idx += 5
	}

	return nil
}

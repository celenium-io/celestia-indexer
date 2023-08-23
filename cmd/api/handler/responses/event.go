package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Event struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggettype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggettype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggettype:"string"`
	Position uint64    `example:"1"                         format:"int64"     json:"position"        swaggettype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggettype:"integer"`

	Type string `enums:"coin_received,coinbase,coin_spent,burn,mint,message,proposer_reward,rewards,commission,liveness,attestation_request,transfer,pay_for_blobs,redelegate,withdraw_rewards,withdraw_commission,create_validator,delegate,edit_validator,unbond,tx,unknown" example:"commission" format:"string" json:"type" swaggettype:"string"`

	Data map[string]any `json:"data"`
}

func NewEvent(event storage.Event) Event {
	result := Event{
		Id:       event.Id,
		Height:   event.Height,
		Time:     event.Time,
		Position: event.Position,
		Type:     string(event.Type),
		Data:     event.Data,
	}

	if event.TxId != nil {
		result.TxId = *event.TxId
	}

	return result
}

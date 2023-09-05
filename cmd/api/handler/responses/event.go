package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Event struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggertype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggertype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggertype:"string"`
	Position uint64    `example:"1"                         format:"int64"     json:"position"        swaggertype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggertype:"integer"`

	Type string `enums:"coin_received,coinbase,coin_spent,burn,mint,message,proposer_reward,rewards,commission,liveness,attestation_request,transfer,celestia.blob.v1.EventPayForBlobs,redelegate,withdraw_rewards,withdraw_commission,create_validator,delegate,edit_validator,unbond,tx,use_feegrant,revoke_feegrant,set_feegrant,update_feegrant,unknown" example:"commission" format:"string" json:"type" swaggertype:"string"`

	Data map[string]any `json:"data"`
}

func NewEvent(event storage.Event) Event {
	result := Event{
		Id:       event.Id,
		Height:   uint64(event.Height),
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

package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Message struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggettype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggettype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggettype:"string"`
	Position uint64    `example:"2"                         format:"int64"     json:"position"        swaggettype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggettype:"integer"`

	Type string `enums:"WithdrawValidatorCommission,WithdrawDelegatorReward,EditValidator,BeginRedelegate,CreateValidator,Delegate,Undelegate,Unjail,Send,CreateVestingAccount,CreatePeriodicVestingAccount,PayForBlobs" example:"CreatePeriodicVestingAccount" format:"string" json:"type" swaggettype:"string"`

	Data map[string]any `json:"data"`
}

func NewMessage(msg storage.Message) Message {
	return Message{
		Id:       msg.Id,
		Height:   msg.Height,
		Time:     msg.Time,
		Position: msg.Position,
		Type:     string(msg.Type),
		TxId:     msg.TxId,
		Data:     msg.Data,
	}
}

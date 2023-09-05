package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Message struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"              swaggertype:"integer"`
	Height   uint64    `example:"100"                       format:"int64"     json:"height"          swaggertype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"            swaggertype:"string"`
	Position uint64    `example:"2"                         format:"int64"     json:"position"        swaggertype:"integer"`
	TxId     uint64    `example:"11"                        format:"int64"     json:"tx_id,omitempty" swaggertype:"integer"`

	Type string `enums:"MsgWithdrawValidatorCommission,MsgWithdrawDelegatorReward,MsgEditValidator,MsgBeginRedelegate,MsgCreateValidator,MsgDelegate,MsgUndelegate,MsgUnjail,MsgSend,MsgCreateVestingAccount,MsgCreatePeriodicVestingAccount,MsgPayForBlobs,MsgGrantAllowance" example:"MsgCreatePeriodicVestingAccount" format:"string" json:"type" swaggertype:"string"`

	Data map[string]any `json:"data"`
}

func NewMessage(msg storage.Message) Message {
	return Message{
		Id:       msg.Id,
		Height:   uint64(msg.Height),
		Time:     msg.Time,
		Position: msg.Position,
		Type:     string(msg.Type),
		TxId:     msg.TxId,
		Data:     msg.Data,
	}
}

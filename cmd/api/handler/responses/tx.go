package responses

import (
	"encoding/hex"
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type Tx struct {
	Id            uint64    `example:"321"                                                              format:"int64"     json:"id"                  swaggettype:"integer"`
	Height        uint64    `example:"100"                                                              format:"int64"     json:"height"              swaggettype:"integer"`
	Position      uint64    `example:"11"                                                               format:"int64"     json:"position"            swaggettype:"integer"`
	GasWanted     uint64    `example:"9348"                                                             format:"int64"     json:"gas_wanted"          swaggettype:"integer"`
	GasUsed       uint64    `example:"4253"                                                             format:"int64"     json:"gas_used"            swaggettype:"integer"`
	TimeoutHeight uint64    `example:"0"                                                                format:"int64"     json:"timeout_height"      swaggettype:"integer"`
	EventsCount   uint64    `example:"2"                                                                format:"int64"     json:"events_count"        swaggettype:"integer"`
	MessagesCount uint64    `example:"1"                                                                format:"int64"     json:"messages_count"      swaggettype:"integer"`
	Hash          string    `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"hash"                swaggettype:"string"`
	Fee           string    `example:"9348"                                                             format:"int64"     json:"fee"                 swaggettype:"string"`
	Error         string    `example:""                                                                 format:"string"    json:"error,omitempty"     swaggettype:"string"`
	Codespace     string    `example:"sdk"                                                              format:"string"    json:"codespace,omitempty" swaggettype:"string"`
	Memo          string    `example:"Transfer to private account"                                      format:"string"    json:"memo,omitempty"      swaggettype:"string"`
	Time          time.Time `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                swaggettype:"string"`

	Status string `enums:"success,failed" example:"success" format:"string" json:"status" swaggettype:"string"`
}

func NewTx(tx storage.Tx) Tx {
	return Tx{
		Id:            tx.Id,
		Height:        tx.Height,
		Time:          tx.Time,
		Position:      tx.Position,
		GasWanted:     tx.GasWanted,
		GasUsed:       tx.GasUsed,
		TimeoutHeight: tx.TimeoutHeight,
		EventsCount:   tx.EventsCount,
		MessagesCount: tx.MessagesCount,
		Fee:           tx.Fee.String(),
		Status:        string(tx.Status),
		Error:         tx.Error,
		Codespace:     tx.Codespace,
		Hash:          hex.EncodeToString(tx.Hash),
		Memo:          tx.Memo,
	}
}

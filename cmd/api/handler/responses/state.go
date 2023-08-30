package responses

import (
	"time"

	"github.com/dipdup-io/celestia-indexer/internal/storage"
)

type State struct {
	Id                 uint64    `example:"321"                       format:"int64"     json:"id"                   swaggettype:"integer"`
	Name               string    `example:"indexer"                   format:"string"    json:"name"                 swaggettype:"string"`
	LastHeight         uint64    `example:"100"                       format:"int64"     json:"last_height"          swaggettype:"integer"`
	LastTime           time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"last_time"            swaggettype:"string"`
	TotalTx            uint64    `example:"23456"                     format:"int64"     json:"total_tx"             swaggettype:"integer"`
	TotalAccounts      uint64    `example:"43"                        format:"int64"     json:"total_accounts"       swaggettype:"integer"`
	TotalFee           string    `example:"312"                       format:"string"    json:"total_fee"            swaggettype:"string"`
	TotalNamespaceSize uint64    `example:"56789"                     format:"int64"     json:"total_namespace_size" swaggettype:"integer"`
}

func NewState(state storage.State) State {
	return State{
		Id:                 state.ID,
		Name:               state.Name,
		LastHeight:         uint64(state.LastHeight),
		LastTime:           state.LastTime,
		TotalTx:            state.TotalTx,
		TotalAccounts:      state.TotalAccounts,
		TotalFee:           state.TotalFee.String(),
		TotalNamespaceSize: state.TotalNamespaceSize,
	}
}

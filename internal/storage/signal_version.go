package storage

import (
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ISignalVersion interface {
}

type SignalVersion struct {
	bun.BaseModel `bun:"signal_version" comment:"Table with signal version"`

	Id          uint64          `bun:"id,pk"                       comment:"Unique identity"`
	Height      pkgTypes.Level  `bun:"height"                      comment:"The number (height) of this block"`
	ValidatorId uint64          `bun:"validator_id"                comment:"Validator address identity"`
	Time        time.Time       `bun:"time,pk,notnull"             comment:"The time of signal"`
	VotingPower decimal.Decimal `bun:"voting_power,type:numeric"   comment:"Voting power"`
	Version     uint64          `bun:"version"                     comment:"Version"`
	MsgId       uint64          `bun:"msg_id,notnull"              comment:"Message internal identity"`
	TxId        uint64          `bun:"tx_id,notnull"               comment:"Transaction internal identity"`

	Validator *Validator `bun:"rel:belongs-to,join:validator_id=id"`
}

// TableName -
func (SignalVersion) TableName() string {
	return "signal_version"
}

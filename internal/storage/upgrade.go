package storage

import (
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IUpgrade interface {
}

type Upgrade struct {
	bun.BaseModel `bun:"upgrade" comment:"Table with upgrades"`

	Id       uint64         `bun:"id,pk"                       comment:"Unique identity"`
	Height   pkgTypes.Level `bun:"height"                      comment:"The number (height) of this block"`
	SignerId uint64         `bun:"signer_id"                   comment:"Signer internal identity"`
	Time     time.Time      `bun:"time,pk,notnull"             comment:"The time of upgrade"`
	Version  uint64         `bun:"version"                     comment:"Version"`
	MsgId    uint64         `bun:"msg_id,notnull"              comment:"Message internal identity"`
	TxId     uint64         `bun:"tx_id,notnull"               comment:"Transaction internal identity"`

	Signer *Address `bun:"rel:belongs-to,join:signer_id=id"`
}

// TableName -
func (Upgrade) TableName() string {
	return "upgrade"
}

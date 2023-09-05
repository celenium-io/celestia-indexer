package storage

import (
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/uptrace/bun"
)

type TxAddress struct {
	bun.BaseModel `bun:"tx_address" comment:"Table with relation tx to address"`

	AddressId uint64              `bun:"address_id,pk"         comment:"Address internal id"`
	TxId      uint64              `bun:"tx_id,pk"              comment:"Transaction internal id"`
	Type      types.TxAddressType `bun:",type:tx_address_type" comment:"The reason why address link to transaction"`

	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
	Tx      *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (TxAddress) TableName() string {
	return "tx_address"
}

type AddressWithType struct {
	Address

	Type types.TxAddressType
}

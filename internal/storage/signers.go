package storage

import (
	"github.com/uptrace/bun"
)

type Signer struct {
	bun.BaseModel `bun:"signer" comment:"Table with signers tx"`

	AddressId uint64 `bun:"address_id,pk" comment:"Address internal id"`
	TxId      uint64 `bun:"tx_id,pk"      comment:"Transaction internal id"`

	Address *Address `bun:"rel:belongs-to,join:address_id=id"`
	Tx      *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (Signer) TableName() string {
	return "signer"
}

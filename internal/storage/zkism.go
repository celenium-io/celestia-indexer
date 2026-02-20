// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/hex"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// ZkISMFilter holds query filters for listing ZK ISM entities.
type ZkISMFilter struct {
	CreatorId *uint64
	TxId      *uint64
	Sort      sdk.SortOrder
	Limit     int
	Offset    int
}

// ZkISMUpdatesFilter holds query filters for listing ZK ISM updates and messages.
type ZkISMUpdatesFilter struct {
	SignerId *uint64
	TxId     *uint64
	From     time.Time
	To       time.Time
	Sort     sdk.SortOrder
	Limit    int
	Offset   int
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IZkISM interface {
	List(ctx context.Context, filter ZkISMFilter) ([]ZkISM, error)
	ById(ctx context.Context, id uint64) (ZkISM, error)
	Updates(ctx context.Context, id uint64, filter ZkISMUpdatesFilter) ([]ZkISMUpdate, error)
	Messages(ctx context.Context, id uint64, filter ZkISMUpdatesFilter) ([]ZkISMMessage, error)
}

// ZkISM represents a ZK Interchain Security Module created via MsgCreateInterchainSecurityModule.
type ZkISM struct {
	bun.BaseModel `bun:"zk_ism" comment:"Table with ZK Interchain Security Modules"`

	Id                  uint64         `bun:"id,pk,autoincrement"              comment:"Internal identity"`
	Height              pkgTypes.Level `bun:"height,notnull"                   comment:"Block height of creation"`
	Time                time.Time      `bun:"time,notnull"                     comment:"Block time of creation"`
	TxId                uint64         `bun:"tx_id"                            comment:"Creation transaction id"`
	CreatorId           uint64         `bun:"creator_id"                       comment:"Creator address identity"`
	ExternalId          []byte         `bun:"external_id,unique,type:bytea"    comment:"Chain-assigned ISM id"`
	State               []byte         `bun:"state,type:bytea"                 comment:"Current trusted state"`
	StateRoot           []byte         `bun:"state_root,type:bytea"            comment:"Current state root (first 32 bytes of state)"`
	MerkleTreeAddress   []byte         `bun:"merkle_tree_address,type:bytea"   comment:"External chain merkle tree address (32 bytes)"`
	Groth16VKey         []byte         `bun:"groth16_vkey,type:bytea"          comment:"On-chain Groth16 verifier key"`
	StateTransitionVKey []byte         `bun:"state_transition_vkey,type:bytea" comment:"State transition verifier key commitment (32 bytes)"`
	StateMembershipVKey []byte         `bun:"state_membership_vkey,type:bytea" comment:"State membership verifier key commitment (32 bytes)"`

	Creator *Address `bun:"rel:belongs-to,join:creator_id=id"`
	Tx      *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (z *ZkISM) TableName() string {
	return "zk_ism"
}

func (z *ZkISM) String() string {
	return hex.EncodeToString(z.StateRoot)
}

func (z *ZkISM) ExternalIdString() string {
	return hex.EncodeToString(z.ExternalId)
}

// ZkISMUpdate represents a state transition recorded via MsgUpdateInterchainSecurityModule.
type ZkISMUpdate struct {
	bun.BaseModel `bun:"zk_ism_update" comment:"Table with ZK ISM state updates"`

	Id           uint64         `bun:"id,pk,autoincrement"       comment:"Internal identity"`
	Height       pkgTypes.Level `bun:"height,notnull"            comment:"Block height"`
	Time         time.Time      `bun:"time,pk,notnull"           comment:"Block time"`
	TxId         uint64         `bun:"tx_id"                     comment:"Transaction id"`
	ZkISMId      uint64         `bun:"zk_ism_id,notnull"         comment:"ZK ISM internal id"`
	SignerId     uint64         `bun:"signer_id"                 comment:"Signer address identity"`
	NewStateRoot []byte         `bun:"new_state_root,type:bytea" comment:"New state root after update"`
	NewState     []byte         `bun:"new_state,type:bytea"      comment:"New full state after update"`

	// Temporary field used during parsing to hold the chain-level ISM id before DB resolution.
	ZkISMExternalId []byte `bun:"-"`

	ZkISM  *ZkISM   `bun:"rel:belongs-to,join:zk_ism_id=id"`
	Signer *Address `bun:"rel:belongs-to,join:signer_id=id"`
	Tx     *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (z *ZkISMUpdate) TableName() string {
	return "zk_ism_update"
}

func (z *ZkISMUpdate) ExternalId() string {
	return hex.EncodeToString(z.ZkISMExternalId)
}

// ZkISMMessage represents an authorized message submitted via MsgSubmitMessages.
type ZkISMMessage struct {
	bun.BaseModel `bun:"zk_ism_message" comment:"Table with ZK ISM authorized messages"`

	Id        uint64         `bun:"id,pk,autoincrement"   comment:"Internal identity"`
	Height    pkgTypes.Level `bun:"height,notnull"        comment:"Block height"`
	Time      time.Time      `bun:"time,pk,notnull"       comment:"Block time"`
	TxId      uint64         `bun:"tx_id"                 comment:"Transaction id"`
	ZkISMId   uint64         `bun:"zk_ism_id,notnull"     comment:"ZK ISM internal id"`
	SignerId  uint64         `bun:"signer_id"             comment:"Signer address identity"`
	StateRoot []byte         `bun:"state_root,type:bytea" comment:"State root at time of authorization"`
	MessageId []byte         `bun:"message_id,type:bytea" comment:"Authorized Hyperlane message id"`

	// Temporary field used during parsing to hold the chain-level ISM id before DB resolution.
	ZkISMExternalId []byte `bun:"-"`

	ZkISM  *ZkISM   `bun:"rel:belongs-to,join:zk_ism_id=id"`
	Signer *Address `bun:"rel:belongs-to,join:signer_id=id"`
	Tx     *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
}

func (z *ZkISMMessage) TableName() string {
	return "zk_ism_message"
}

func (z *ZkISMMessage) ExternalId() string {
	return hex.EncodeToString(z.ZkISMExternalId)
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"io"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

const (
	columnNameHeight   = "height"
	columnNameTime     = "time"
	columnNamePosition = "position"
	columnNameTxId     = "tx_id"
	columnNameType     = "type"
)

var Models = []any{
	&State{},
	&Constant{},
	&DenomMetadata{},
	&Balance{},
	&Address{},
	&VestingAccount{},
	&VestingPeriod{},
	&Block{},
	&BlockStats{},
	&BlockSignature{},
	&Tx{},
	&Message{},
	&Event{},
	&Namespace{},
	&NamespaceMessage{},
	&Signer{},
	&MsgAddress{},
	&MsgValidator{},
	&Validator{},
	&Delegation{},
	&Redelegation{},
	&Undelegation{},
	&StakingLog{},
	&Jail{},
	&BlobLog{},
	&Rollup{},
	&RollupProvider{},
	&Grant{},
	&ApiKey{},
	&celestials.Celestial{},
	&celestials.CelestialState{},
	&Proposal{},
	&Vote{},
	&IbcClient{},
	&IbcConnection{},
	&IbcChannel{},
	&IbcTransfer{},
	&HLMailbox{},
	&HLToken{},
	&HLTransfer{},
	&SignalVersion{},
	&Upgrade{},
	&HLIGP{},
	&HLIGPConfig{},
	&HLGasPayment{},
	&Forwarding{},
	&ZkISM{},
	&ZkISMUpdate{},
	&ZkISMMessage{},
	&ValidatorBondUpdate{},
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Notificator interface {
	Notify(ctx context.Context, channel string, payload string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Listener interface {
	io.Closer

	Subscribe(ctx context.Context, channels ...string) error
	Listen() <-chan pgconn.Notification
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ListenerFactory interface {
	CreateListener() Listener
}

const (
	ChannelHead  = "head"
	ChannelBlock = "block"
)

type Signal struct {
	VotingPower types.Numeric `bun:"voting_power"`
	Version     uint64        `bun:"version"`
}

type SearchResult struct {
	Id    uint64 `bun:"id"`
	Value string `bun:"value"`
	Type  string `bun:"type"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type ISearch interface {
	Search(ctx context.Context, query []byte) ([]SearchResult, error)
	SearchText(ctx context.Context, text string) ([]SearchResult, error)
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type Export interface {
	ToCsv(ctx context.Context, writer io.Writer, query *bun.SelectQuery) error
}

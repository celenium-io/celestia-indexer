// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"time"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IIbcClient interface {
	ById(ctx context.Context, id string) (IbcClient, error)
	List(ctx context.Context, limit, offset int, sort sdk.SortOrder) ([]IbcClient, error)
}

type IbcClient struct {
	bun.BaseModel `bun:"ibc_client" comment:"Table with IBC clients."`

	Id                    string         `bun:"id,pk,notnull"           comment:"Client identity"`
	ChainId               string         `bun:"chain_id"                comment:"Chain id"`
	Type                  string         `bun:"type"                    comment:"Light client type"`
	CreatedAt             time.Time      `bun:"created_at"              comment:"Creation time"`
	UpdatedAt             time.Time      `bun:"updated_at"              comment:"Time of last update message"`
	Height                pkgTypes.Level `bun:"height"                  comment:"Creation height"`
	TxId                  uint64         `bun:"tx_id"                   comment:"Internal transaction identity where client was created"`
	CreatorId             uint64         `bun:"creator_id"              comment:"Creator internal identity"`
	LatestRevisionHeight  uint64         `bun:"latest_revision_height"  comment:"Latest height the client was updated to"`
	LatestRevisionNumber  uint64         `bun:"latest_revision_number"  comment:"Revision number the client was updated to"`
	FrozenRevisionHeight  uint64         `bun:"frozen_revision_height"  comment:"Block height when the client was frozen due to a misbehaviour"`
	FrozenRevisionNumber  uint64         `bun:"frozen_revision_number"  comment:"Revision number when the client was frozen due to a misbehaviour"`
	TrustingPeriod        time.Duration  `bun:"trusting_period"         comment:"Duration of the period since the LastestTimestamp during which the submitted headers are valid for upgrade"`
	UnbondingPeriod       time.Duration  `bun:"unbonding_period"        comment:"Duration of the staking unbonding period"`
	MaxClockDrift         time.Duration  `bun:"max_clock_drift"         comment:"Defines how much new (untrusted) header's time can drift into the future"`
	TrustLevelDenominator uint64         `bun:"trust_level_denominator" comment:"Denominator of trust level"`
	TrustLevelNumerator   uint64         `bun:"trust_level_numerator"   comment:"Numerator of trust level"`
	ConnectionCount       uint64         `bun:"connection_count"        comment:"Count of connections which is associated with client"`

	Tx      *Tx      `bun:"rel:belongs-to,join:tx_id=id"`
	Creator *Address `bun:"rel:belongs-to,join:creator_id=id"`
}

func (IbcClient) TableName() string {
	return "ibc_client"
}

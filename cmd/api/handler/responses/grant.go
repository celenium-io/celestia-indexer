// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Grant struct {
	Authorization string     `example:"/cosmos.staking.v1beta1.MsgDelegate" json:"authorization"           swaggertype:"string"`
	Expiration    *time.Time `example:"2023-07-04T03:10:57+00:00"           json:"expiration,omitempty"    swaggertype:"string"`
	Revoked       bool       `example:"true"                                json:"revoked"                 swaggertype:"boolean"`
	RevokeHeight  uint64     `example:"123123"                              json:"revoke_height,omitempty" swaggertype:"integer"`
	Height        uint64     `example:"123123"                              json:"height"                  swaggertype:"integer"`
	Time          time.Time  `example:"2023-07-04T03:10:57+00:00"           json:"time"                    swaggertype:"string"`

	Params  map[string]any `json:"params"`
	Granter *ShortAddress  `json:"granter,omitempty"`
	Grantee *ShortAddress  `json:"grantee,omitempty"`
}

func NewGrant(g storage.Grant) Grant {
	grant := Grant{
		Height:        uint64(g.Height),
		Authorization: g.Authorization,
		Expiration:    g.Expiration,
		Revoked:       g.Revoked,
		Params:        g.Params,
		Time:          g.Time,
		Granter:       NewShortAddress(g.Granter),
		Grantee:       NewShortAddress(g.Grantee),
	}

	if g.RevokeHeight != nil {
		grant.RevokeHeight = uint64(*g.RevokeHeight)
	}

	return grant
}

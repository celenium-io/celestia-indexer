// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
)

type Grant struct {
	Granter       string     `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"granter,omitempty"       swaggertype:"string"`
	Grantee       string     `example:"celestia1jc92qdnty48pafummfr8ava2tjtuhfdw774w60" json:"grantee,omitempty"       swaggertype:"string"`
	Authorization string     `example:"/cosmos.staking.v1beta1.MsgDelegate"             json:"authorization"           swaggertype:"string"`
	Expiration    *time.Time `example:"2023-07-04T03:10:57+00:00"                       json:"expiration,omitempty"    swaggertype:"string"`
	Revoked       bool       `example:"true"                                            json:"revoked"                 swaggertype:"boolean"`
	RevokeHeight  uint64     `example:"123123"                                          json:"revoke_height,omitempty" swaggertype:"integer"`
	Height        uint64     `example:"123123"                                          json:"height"                  swaggertype:"integer"`
	Time          time.Time  `example:"2023-07-04T03:10:57+00:00"                       json:"time"                    swaggertype:"string"`

	Params map[string]any `json:"params"`
}

func NewGrant(g storage.Grant) Grant {
	grant := Grant{
		Height:        uint64(g.Height),
		Authorization: g.Authorization,
		Expiration:    g.Expiration,
		Revoked:       g.Revoked,
		Params:        g.Params,
		Time:          g.Time,
	}

	if g.Grantee != nil {
		grant.Grantee = g.Grantee.Address
	}
	if g.Granter != nil {
		grant.Granter = g.Granter.Address
	}
	if g.RevokeHeight != nil {
		grant.RevokeHeight = uint64(*g.RevokeHeight)
	}

	return grant
}

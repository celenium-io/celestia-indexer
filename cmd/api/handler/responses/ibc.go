// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type IbcClient struct {
	Id                    string         `example:"client-1"                                                         format:"string"    json:"id"                      swaggertype:"string"`
	Type                  string         `example:"client"                                                           format:"string"    json:"type"                    swaggertype:"string"`
	CreatedAt             time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"created_at"              swaggertype:"string"`
	UpdatedAt             time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"updated_at"              swaggertype:"string"`
	Height                pkgTypes.Level `example:"100"                                                              format:"integer"   json:"height"                  swaggertype:"integer"`
	ChainId               string         `example:"osmosis-1"                                                        format:"binary"    json:"chain_id"                swaggertype:"string"`
	TxHash                string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash"                 swaggertype:"string"`
	LatestRevisionHeight  uint64         `example:"100"                                                              format:"integer"   json:"latest_revision_height"  swaggertype:"integer"`
	LatestRevisionNumber  uint64         `example:"100"                                                              format:"integer"   json:"latest_revision_number"  swaggertype:"integer"`
	FrozenRevisionHeight  uint64         `example:"100"                                                              format:"integer"   json:"frozen_revision_height"  swaggertype:"integer"`
	FrozenRevisionNumber  uint64         `example:"100"                                                              format:"integer"   json:"frozen_revision_number"  swaggertype:"integer"`
	TrustingPeriod        time.Duration  `example:"100"                                                              format:"integer"   json:"trusting_period"         swaggertype:"integer"`
	UnbondingPeriod       time.Duration  `example:"100"                                                              format:"integer"   json:"unbonding_period"        swaggertype:"integer"`
	MaxClockDrift         time.Duration  `example:"100"                                                              format:"integer"   json:"max_clock_drift"         swaggertype:"integer"`
	TrustLevelDenominator uint64         `example:"100"                                                              format:"integer"   json:"trust_level_denominator" swaggertype:"integer"`
	TrustLevelNumerator   uint64         `example:"100"                                                              format:"integer"   json:"trust_level_numerator"   swaggertype:"integer"`
	ConnectionCount       uint64         `example:"100"                                                              format:"integer"   json:"connection_count"        swaggertype:"integer"`

	Creator *ShortAddress `json:"creator,omitempty"`
}

func NewIbcClient(client storage.IbcClient) IbcClient {
	response := IbcClient{
		Id:                    client.Id,
		Type:                  client.Type,
		CreatedAt:             client.CreatedAt,
		UpdatedAt:             client.UpdatedAt,
		Height:                client.Height,
		LatestRevisionHeight:  client.LatestRevisionHeight,
		LatestRevisionNumber:  client.LatestRevisionNumber,
		FrozenRevisionHeight:  client.FrozenRevisionHeight,
		FrozenRevisionNumber:  client.FrozenRevisionNumber,
		TrustingPeriod:        client.TrustingPeriod,
		UnbondingPeriod:       client.UnbondingPeriod,
		MaxClockDrift:         client.MaxClockDrift,
		TrustLevelDenominator: client.TrustLevelDenominator,
		TrustLevelNumerator:   client.TrustLevelNumerator,
		ConnectionCount:       client.ConnectionCount,
		ChainId:               client.ChainId,
		Creator:               NewShortAddress(client.Creator),
	}

	if client.Tx != nil {
		response.TxHash = hex.EncodeToString(client.Tx.Hash)
	}

	return response
}

type ShortIbcClient struct {
	Id      string `example:"client-1"  format:"string" json:"id"       swaggertype:"string"`
	Type    string `example:"client"    format:"string" json:"type"     swaggertype:"string"`
	ChainId string `example:"osmosis-1" format:"binary" json:"chain_id" swaggertype:"string"`
}

func NewShortIbcClient(client *storage.IbcClient) *ShortIbcClient {
	if client == nil {
		return nil
	}
	return &ShortIbcClient{
		ChainId: client.ChainId,
		Type:    client.Type,
		Id:      client.Id,
	}
}

type IbcConnection struct {
	Id                   string         `example:"connection-1"                                                     format:"string"    json:"id"                         swaggertype:"string"`
	CounterpartyConnId   string         `example:"connection-1"                                                     format:"string"    json:"counterparty_connection_id" swaggertype:"string"`
	CounterpartyClientId string         `example:"client-1"                                                         format:"string"    json:"counterparty_client_id"     swaggertype:"string"`
	CreatedAt            time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"created_at"                 swaggertype:"string"`
	ConnectedAt          time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"connected_at"               swaggertype:"string"`
	Height               pkgTypes.Level `example:"100"                                                              format:"integer"   json:"height"                     swaggertype:"integer"`
	ConnectedHeight      pkgTypes.Level `example:"100"                                                              format:"integer"   json:"connected_height"           swaggertype:"integer"`
	CreatedTxHash        string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"created_tx_hash"            swaggertype:"string"`
	ConnectedTxHash      string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"connected_tx_hash"          swaggertype:"string"`
	ChannelsCount        int64          `example:"100"                                                              format:"integer"   json:"channels_count"             swaggertype:"integer"`

	Client *ShortIbcClient `json:"client,omitempty"`
}

func NewIbcConnection(conn storage.IbcConnection) IbcConnection {
	conn.Client.Id = conn.ClientId
	response := IbcConnection{
		Id:                   conn.ConnectionId,
		CounterpartyConnId:   conn.CounterpartyConnectionId,
		CounterpartyClientId: conn.CounterpartyClientId,
		CreatedAt:            conn.CreatedAt,
		ConnectedAt:          conn.ConnectedAt,
		Height:               conn.Height,
		ConnectedHeight:      conn.ConnectionHeight,
		Client:               NewShortIbcClient(conn.Client),
		ChannelsCount:        conn.ChannelsCount,
	}
	if conn.CreateTx != nil {
		response.CreatedTxHash = hex.EncodeToString(conn.CreateTx.Hash)
	}
	if conn.ConnectionTx != nil {
		response.ConnectedTxHash = hex.EncodeToString(conn.ConnectionTx.Hash)
	}

	return response
}

type IbcChannel struct {
	Id                    string         `example:"channel-1"                                                        format:"string"    json:"id"                             swaggertype:"string"`
	ConnectionId          string         `example:"connection-1"                                                     format:"string"    json:"connection_id"                  swaggertype:"string"`
	PortId                string         `example:"transfer"                                                         format:"string"    json:"port_id"                        swaggertype:"string"`
	CounterpartyPortId    string         `example:"transfer"                                                         format:"string"    json:"counterparty_port_id"           swaggertype:"string"`
	CounterpartyChannelId string         `example:"channel-1"                                                        format:"string"    json:"counterparty_channel_id"        swaggertype:"string"`
	Version               string         `example:"ics20-1"                                                          format:"string"    json:"version"                        swaggertype:"string"`
	CreatedAt             time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"created_at"                     swaggertype:"string"`
	ConfirmedAt           *time.Time     `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"confirmed_at,omitempty"         swaggertype:"string"`
	Height                pkgTypes.Level `example:"100"                                                              format:"integer"   json:"height"                         swaggertype:"integer"`
	ConfirmationHeight    pkgTypes.Level `example:"100"                                                              format:"integer"   json:"confirmation_height,omitempty"  swaggertype:"integer"`
	CreatedTxHash         string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"created_tx_hash"                swaggertype:"string"`
	ConfirmationTxHash    string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"confirmation_tx_hash,omitempty" swaggertype:"string"`
	Ordering              bool           `example:"false"                                                            format:"boolean"   json:"ordering"                       swaggertype:"boolean"`
	Status                string         `example:"opened"                                                           format:"string"    json:"status"                         swaggertype:"string"`

	Client  *ShortIbcClient `json:"client,omitempty"`
	Creator *ShortAddress   `json:"creator,omitempty"`
}

func NewIbcChannel(channel storage.IbcChannel) IbcChannel {
	channel.Client.Id = channel.ClientId
	response := IbcChannel{
		Id:                    channel.Id,
		ConnectionId:          channel.ConnectionId,
		PortId:                channel.PortId,
		CounterpartyPortId:    channel.CounterpartyPortId,
		CounterpartyChannelId: channel.CounterpartyChannelId,
		Version:               channel.Version,
		CreatedAt:             channel.CreatedAt,
		Height:                channel.Height,
		ConfirmationHeight:    channel.ConfirmationHeight,
		Ordering:              channel.Ordering,
		Status:                channel.Status.String(),
		Client:                NewShortIbcClient(channel.Client),
		Creator:               NewShortAddress(channel.Creator),
	}

	if channel.CreateTx != nil {
		response.CreatedTxHash = hex.EncodeToString(channel.CreateTx.Hash)
	}
	if channel.ConfirmationTx != nil {
		response.ConfirmationTxHash = hex.EncodeToString(channel.ConfirmationTx.Hash)
	}
	if !channel.ConfirmedAt.IsZero() {
		response.ConfirmedAt = &channel.ConfirmedAt
	}

	return response
}

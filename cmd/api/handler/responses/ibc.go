// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	ChainId               string         `example:"osmosis-1"                                                        format:"string"    json:"chain_id"                swaggertype:"string"`
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

type IbcTransfer struct {
	Id            uint64         `example:"123456"                                                           format:"integer"   json:"id"                       swaggertype:"integer"`
	Time          time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                     swaggertype:"string"`
	Height        pkgTypes.Level `example:"100"                                                              format:"integer"   json:"height"                   swaggertype:"integer"`
	ChannelId     string         `example:"channel-1"                                                        format:"string"    json:"channel_id"               swaggertype:"string"`
	ConnectionId  string         `example:"connection-1"                                                     format:"string"    json:"connection_id"            swaggertype:"string"`
	Port          string         `example:"transfer"                                                         format:"string"    json:"port"                     swaggertype:"string"`
	Amount        string         `example:"123445"                                                           format:"string"    json:"amount"                   swaggertype:"string"`
	Denom         string         `example:"utia"                                                             format:"string"    json:"denom"                    swaggertype:"string"`
	Memo          string         `example:"memo"                                                             format:"string"    json:"memo,omitempty"           swaggertype:"string"`
	Timeout       *time.Time     `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"timeout,omitempty"        swaggertype:"string"`
	TimeoutHeight uint64         `example:"100"                                                              format:"integer"   json:"timeout_height,omitempty" swaggertype:"integer"`
	TxHash        string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash"                  swaggertype:"string"`
	Sequence      uint64         `example:"123456"                                                           format:"integer"   json:"sequence"                 swaggertype:"integer"`
	ChainId       string         `example:"osmosis-1"                                                        format:"binary"    json:"chain_id"                 swaggertype:"string"`

	Sender   *ShortAddress `json:"sender,omitempty"`
	Receiver *ShortAddress `json:"receiver,omitempty"`
	Relayer  *Relayer      `json:"relayer,omitempty"`
}

func NewIbcTransfer(transfer storage.IbcTransfer) IbcTransfer {
	response := IbcTransfer{
		Id:            transfer.Id,
		Time:          transfer.Time,
		Height:        transfer.Height,
		ChannelId:     transfer.ChannelId,
		ConnectionId:  transfer.ConnectionId,
		Port:          transfer.Port,
		Amount:        transfer.Amount.String(),
		Denom:         transfer.Denom,
		Memo:          transfer.Memo,
		Timeout:       transfer.Timeout,
		TimeoutHeight: transfer.HeightTimeout,
		Sequence:      transfer.Sequence,
		Sender:        NewShortAddress(transfer.Sender),
		Receiver:      NewShortAddress(transfer.Receiver),
	}

	if transfer.ReceiverAddress != nil {
		response.Receiver = &ShortAddress{
			Hash: *transfer.ReceiverAddress,
		}
	}

	if transfer.SenderAddress != nil {
		response.Sender = &ShortAddress{
			Hash: *transfer.SenderAddress,
		}
	}

	if transfer.Tx != nil {
		response.TxHash = hex.EncodeToString(transfer.Tx.Hash)
	}

	if transfer.Connection != nil && transfer.Connection.Client != nil {
		response.ChainId = transfer.Connection.Client.ChainId
	}

	return response
}

func NewIbcTransferWithRelayer(transfer storage.IbcTransferWithSigner, relayers map[uint64]Relayer) IbcTransfer {
	response := NewIbcTransfer(transfer.IbcTransfer)

	if transfer.SignerId != nil && len(relayers) > 0 {
		if relayer, ok := relayers[*transfer.SignerId]; ok {
			response.Relayer = &relayer
		}
	}

	return response
}

type IbcChainStats struct {
	Chain    string `example:"123456" format:"string" json:"chain"    swaggertype:"string"`
	Sent     string `example:"123445" format:"string" json:"sent"     swaggertype:"string"`
	Received string `example:"123445" format:"string" json:"received" swaggertype:"string"`
	Flow     string `example:"123445" format:"string" json:"flow"     swaggertype:"string"`
}

func NewIbcChainStats(stats storage.ChainStats) IbcChainStats {
	return IbcChainStats{
		Chain:    stats.Chain,
		Received: stats.Received.String(),
		Sent:     stats.Sent.String(),
		Flow:     stats.Flow.String(),
	}
}

type BusiestChannel struct {
	ChannelId      string `example:"channel-1" format:"string"  json:"channel_id"      swaggertype:"string"`
	TransfersCount int64  `example:"100"       format:"integer" json:"transfers_count" swaggertype:"integer"`
	ChainId        string `example:"osmosis-1" format:"string"  json:"chain_id"        swaggertype:"string"`
}

type IbcSummaryStats struct {
	LargestTransfer IbcTransfer    `json:"largest_transfer,omitempty"`
	BusiestChannel  BusiestChannel `json:"busiest_channel,omitempty"`
}

func NewIbcSummaryStats(transfer storage.IbcTransfer, channel storage.BusiestChannel) IbcSummaryStats {
	return IbcSummaryStats{
		LargestTransfer: NewIbcTransfer(transfer),
		BusiestChannel: BusiestChannel{
			ChannelId:      channel.ChannelId,
			TransfersCount: channel.TransfersCount,
			ChainId:        channel.ChainId,
		},
	}
}

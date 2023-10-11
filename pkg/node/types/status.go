// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package types

import (
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
	"time"
)

type Status struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
	ID              string          `json:"id"`
	ListenAddr      string          `json:"listen_addr"`
	Network         string          `json:"network"`
	Version         string          `json:"version"`
	Channels        string          `json:"channels"`
	Moniker         string          `json:"moniker"`
	Other           Other           `json:"other"`
}

type Other struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type ProtocolVersion struct {
	P2P   string `json:"p2p"`
	Block uint64 `json:"block,string"`
	App   uint64 `json:"app,string"`
}

type SyncInfo struct {
	LatestBlockHash     []byte         `json:"latest_block_hash"`
	LatestAppHash       []byte         `json:"latest_app_hash"`
	LatestBlockHeight   pkgTypes.Level `json:"latest_block_height,string"`
	LatestBlockTime     time.Time      `json:"latest_block_time"`
	EarliestBlockHash   []byte         `json:"earliest_block_hash"`
	EarliestAppHash     []byte         `json:"earliest_app_hash"`
	EarliestBlockHeight pkgTypes.Level `json:"earliest_block_height,string"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`
	CatchingUp          bool           `json:"catching_up"`
}

type ValidatorInfo struct {
	Address     []byte `json:"address"`
	PubKey      PubKey `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

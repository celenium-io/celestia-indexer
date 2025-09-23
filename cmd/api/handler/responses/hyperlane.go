// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type HyperlaneMailbox struct {
	Id               uint64         `example:"321"                                                              format:"int64"     json:"id"                      swaggertype:"integer"`
	InternalId       uint64         `example:"321"                                                              format:"int64"     json:"hyperlane_id"            swaggertype:"integer"`
	Height           pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"                  swaggertype:"integer"`
	Time             time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                    swaggertype:"string"`
	TxHash           string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"       swaggertype:"string"`
	Mailbox          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"mailbox"                 swaggertype:"string"`
	DefaultIsm       string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"default_ism,omitempty"   swaggertype:"string"`
	DefaultHook      string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"default_hook,omitempty"  swaggertype:"string"`
	RequiredHook     string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"required_hook,omitempty" swaggertype:"string"`
	Domain           uint64         `example:"100"                                                              format:"int64"     json:"domain,omitempty"        swaggertype:"integer"`
	SentMessages     uint64         `example:"100"                                                              format:"int64"     json:"sent_messages"           swaggertype:"integer"`
	ReceivedMessages uint64         `example:"100"                                                              format:"int64"     json:"received_messages"       swaggertype:"integer"`

	Owner *ShortAddress `json:"owner,omitempty"`
}

func NewHyperlaneMailbox(mailbox storage.HLMailbox) HyperlaneMailbox {
	result := HyperlaneMailbox{
		Id:               mailbox.Id,
		Height:           mailbox.Height,
		Time:             mailbox.Time,
		Mailbox:          hex.EncodeToString(mailbox.Mailbox),
		InternalId:       mailbox.InternalId,
		Domain:           mailbox.Domain,
		SentMessages:     mailbox.SentMessages,
		ReceivedMessages: mailbox.ReceivedMessages,
		Owner:            NewShortAddress(mailbox.Owner),
	}

	if len(mailbox.DefaultHook) > 0 {
		result.DefaultHook = hex.EncodeToString(mailbox.DefaultHook)
	}
	if len(mailbox.RequiredHook) > 0 {
		result.RequiredHook = hex.EncodeToString(mailbox.RequiredHook)
	}
	if len(mailbox.DefaultIsm) > 0 {
		result.DefaultIsm = hex.EncodeToString(mailbox.DefaultIsm)
	}

	if mailbox.Tx != nil {
		result.TxHash = hex.EncodeToString(mailbox.Tx.Hash)
	}

	return result
}

type HyperlaneToken struct {
	Id               uint64         `example:"321"                                                              format:"int64"     json:"id"                 swaggertype:"integer"`
	Height           pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"             swaggertype:"integer"`
	Time             time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"               swaggertype:"string"`
	Mailbox          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"mailbox"            swaggertype:"string"`
	TxHash           string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"  swaggertype:"string"`
	Type             string         `example:"collateral"                                                       format:"string"    json:"type"               swaggertype:"string"`
	Denom            string         `example:"utia"                                                             format:"string"    json:"denom"              swaggertype:"string"`
	TokenId          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"token_id"           swaggertype:"string"`
	SentTransfers    uint64         `example:"100"                                                              format:"int64"     json:"sent_transfers"     swaggertype:"integer"`
	ReceiveTransfers uint64         `example:"100"                                                              format:"int64"     json:"received_transfers" swaggertype:"integer"`
	Sent             string         `example:"123445"                                                           format:"string"    json:"sent"               swaggertype:"string"`
	Received         string         `example:"123445"                                                           format:"string"    json:"received"           swaggertype:"string"`

	Owner *ShortAddress `json:"owner,omitempty"`
}

func NewHyperlaneToken(token storage.HLToken) HyperlaneToken {
	result := HyperlaneToken{
		Id:               token.Id,
		Height:           token.Height,
		Time:             token.Time,
		Type:             token.Type.String(),
		Denom:            token.Denom,
		TokenId:          hex.EncodeToString(token.TokenId),
		SentTransfers:    token.SentTransfers,
		ReceiveTransfers: token.ReceiveTransfers,
		Sent:             token.Sent.String(),
		Received:         token.Received.String(),
		Owner:            NewShortAddress(token.Owner),
	}

	if token.Mailbox != nil {
		result.Mailbox = hex.EncodeToString(token.Mailbox.Mailbox)
	}
	if token.Tx != nil {
		result.TxHash = hex.EncodeToString(token.Tx.Hash)
	}

	return result
}

type HyperlaneTransfer struct {
	Id       uint64         `example:"321"                                                              format:"int64"     json:"id"                 swaggertype:"integer"`
	Height   pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"             swaggertype:"integer"`
	Time     time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"               swaggertype:"string"`
	TxHash   string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"  swaggertype:"string"`
	Mailbox  string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"mailbox"            swaggertype:"string"`
	TokenId  string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"token_id"           swaggertype:"string"`
	Type     string         `example:"collateral"                                                       format:"string"    json:"type"               swaggertype:"string"`
	Version  byte           `example:"1"                                                                format:"int64"     json:"version"            swaggertype:"integer"`
	Nonce    uint32         `example:"10"                                                               format:"int64"     json:"nonce"              swaggertype:"integer"`
	Body     []byte         `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="                         format:"string"    json:"body,omitempty"     swaggertype:"string"`
	Metadata []byte         `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="                         format:"string"    json:"metadata,omitempty" swaggertype:"string"`
	Amount   string         `example:"123445"                                                           format:"string"    json:"received"           swaggertype:"string"`
	Denom    string         `example:"utia"                                                             format:"string"    json:"denom"              swaggertype:"string"`

	Address      *ShortAddress         `json:"address,omitempty"`
	Relayer      *ShortAddress         `json:"relayer,omitempty"`
	Counterparty HyperlaneCounterparty `json:"counterparty"`
}

func NewHyperlaneTransfer(transfer storage.HLTransfer, store hyperlane.IChainStore) HyperlaneTransfer {
	counterparty := HyperlaneCounterparty{
		Hash:   transfer.CounterpartyAddress,
		Domain: transfer.Counterparty,
	}

	if store != nil {
		counterparty.ChainMetadata = NewChainMetadata(transfer.Counterparty, store)
	}

	result := HyperlaneTransfer{
		Id:           transfer.Id,
		Height:       transfer.Height,
		Time:         transfer.Time,
		Type:         transfer.Type.String(),
		Version:      transfer.Version,
		Nonce:        transfer.Nonce,
		Body:         transfer.Body,
		Metadata:     transfer.Metadata,
		Denom:        transfer.Denom,
		Amount:       transfer.Amount.String(),
		Address:      NewShortAddress(transfer.Address),
		Relayer:      NewShortAddress(transfer.Relayer),
		Counterparty: counterparty,
	}

	if transfer.Token != nil {
		result.TokenId = hex.EncodeToString(transfer.Token.TokenId)
	}
	if transfer.Mailbox != nil {
		result.Mailbox = hex.EncodeToString(transfer.Mailbox.Mailbox)
	}
	if transfer.Tx != nil {
		result.TxHash = hex.EncodeToString(transfer.Tx.Hash)
	}

	return result
}

type HyperlaneCounterparty struct {
	Hash          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" json:"hash"    swaggertype:"string"`
	Domain        uint64         `example:"100"                                                              format:"int64" json:"domain,omitempty" swaggertype:"integer"`
	ChainMetadata *ChainMetadata `json:"chain_metadata,omitempty"`
}

type ChainMetadata struct {
	Name           string          `example:"name"                   json:"name" swaggertype:"string"`
	BlockExplorers []BlockExplorer `json:"block_explorers,omitempty"`
	NativeToken    NativeToken     `json:"native_token,omitempty"`
}

type BlockExplorer struct {
	ApiUrl string `example:"https://api.scan.url.io" format:"string" json:"api_url" swaggertype:"string"`
	Family string `example:"etherscan"               format:"string" json:"family"  swaggertype:"string"`
	Name   string `example:"Block explorer"          format:"string" json:"name"    swaggertype:"string"`
	Url    string `example:"https://scan.url.io"     format:"string" json:"url"     swaggertype:"string"`
}

type NativeToken struct {
	Decimals uint64 `example:"18"    format:"int64"  json:"decimals" swaggertype:"integer"`
	Name     string `example:"Ether" format:"string" json:"name"     swaggertype:"string"`
	Symbol   string `example:"ETH"   format:"string" json:"symbol"   swaggertype:"string"`
}

func NewChainMetadata(domainId uint64, store hyperlane.IChainStore) *ChainMetadata {
	if metadata, ok := store.Get(domainId); ok {
		explorers := make([]BlockExplorer, len(metadata.BlockExplorers))
		for i := range explorers {
			explorers[i] = BlockExplorer(metadata.BlockExplorers[i])
		}
		return &ChainMetadata{
			Name:           metadata.DisplayName,
			BlockExplorers: explorers,
			NativeToken: NativeToken{
				Decimals: metadata.NativeToken.Decimals,
				Name:     metadata.NativeToken.Name,
				Symbol:   metadata.NativeToken.Symbol,
			},
		}
	}

	return nil
}

type DomainMetadata struct {
	Domain         uint64          `example:"1488"         json:"domain,omitempty" swaggertype:"integer"`
	Name           string          `example:"name"         json:"name,omitempty"   swaggertype:"string"`
	BlockExplorers []BlockExplorer `json:"block_explorers"`
	NativeToken    NativeToken     `json:"native_token"`
}

func NewDomainMetadata(domainId uint64, store hyperlane.IChainStore) *DomainMetadata {
	if metadata, ok := store.Get(domainId); ok {
		explorers := make([]BlockExplorer, len(metadata.BlockExplorers))
		for i := range explorers {
			explorers[i] = BlockExplorer(metadata.BlockExplorers[i])
		}
		return &DomainMetadata{
			Domain:         domainId,
			Name:           metadata.DisplayName,
			BlockExplorers: explorers,
			NativeToken: NativeToken{
				Decimals: metadata.NativeToken.Decimals,
				Name:     metadata.NativeToken.Name,
				Symbol:   metadata.NativeToken.Symbol,
			},
		}
	}

	return nil
}

type HlDomainStats struct {
	Domain         uint64         `example:"123456"                format:"integer" json:"domain_id"       swaggertype:"integer"`
	Amount         string         `example:"1234.5678"             format:"string"  json:"amount"          swaggertype:"string"`
	TransfersCount uint64         `example:"123445"                format:"integer" json:"transfers_count" swaggertype:"integer"`
	ChainMetadata  *ChainMetadata `json:"chain_metadata,omitempty"`
}

func NewHlDomainStats(stats storage.DomainStats, store hyperlane.IChainStore) HlDomainStats {
	result := HlDomainStats{
		Domain:         stats.Domain,
		Amount:         stats.Amount.String(),
		TransfersCount: stats.TxCount,
	}
	if store != nil {
		result.ChainMetadata = NewChainMetadata(stats.Domain, store)
	}

	return result
}

type HyperlaneIgp struct {
	Id     uint64         `example:"321"                                                              format:"int64"     json:"id"                 swaggertype:"integer"`
	Height pkgTypes.Level `example:"1488"                                                              format:"int64"     json:"height"             swaggertype:"integer"`
	Time   time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"               swaggertype:"string"`
	Denom  string         `example:"utia"                                                             format:"string"    json:"denom"              swaggertype:"string"`
	IgpId  string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"igp_id"           swaggertype:"string"`

	Owner  *ShortAddress       `json:"owner,omitempty"`
	Config *HyperlaneIgpConfig `json:"config,omitempty"`
}

func NewHyperlaneIgp(igp storage.HLIGP) HyperlaneIgp {
	result := HyperlaneIgp{
		Id:     igp.Id,
		Height: igp.Height,
		Time:   igp.Time,
		Denom:  igp.Denom,
		IgpId:  hex.EncodeToString(igp.IgpId),
		Owner:  NewShortAddress(igp.Owner),
	}

	if igp.Config != nil {
		result.Config = NewHyperlaneIgpConfig(igp.Config)
	}

	return result
}

type HyperlaneIgpConfig struct {
	Id                uint64         `example:"321"                                                              format:"int64"     json:"id"                 swaggertype:"integer"`
	Height            pkgTypes.Level `example:"1488"                                                              format:"int64"     json:"height"             swaggertype:"integer"`
	Time              time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"               swaggertype:"string"`
	IgpId             string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"igp_id"           swaggertype:"string"`
	GasOverhead       string         `example:"100000"                                                             format:"int64"     json:"gas_overhead"                 swaggertype:"string"`
	GasPrice          string         `example:"1"                                                             format:"int64"     json:"gas_price"                 swaggertype:"string"`
	RemoteDomain      uint64         `example:"100"                                                              format:"int64"     json:"remote_domain"        swaggertype:"integer"`
	TokenExchangeRate string         `example:"12345678"                                                              format:"int64"     json:"token_exchange_rate"        swaggertype:"string"`
}

func NewHyperlaneIgpConfig(igp *storage.HLIGPConfig) *HyperlaneIgpConfig {
	result := &HyperlaneIgpConfig{
		Id:                igp.Id,
		Height:            igp.Height,
		Time:              igp.Time,
		IgpId:             hex.EncodeToString(igp.IgpId),
		GasOverhead:       igp.GasOverhead.String(),
		GasPrice:          igp.GasPrice.String(),
		RemoteDomain:      igp.RemoteDomain,
		TokenExchangeRate: igp.TokenExchangeRate,
	}

	return result
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/hex"
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
)

type ZkISM struct {
	Id                  uint64         `example:"321"                                                              format:"int64"     json:"id"                              swaggertype:"integer"`
	ExternalId          string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"external_id"                     swaggertype:"string"`
	Height              pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"                          swaggertype:"integer"`
	Time                time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                            swaggertype:"string"`
	TxHash              string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"               swaggertype:"string"`
	State               string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"state,omitempty"                 swaggertype:"string"`
	StateRoot           string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"state_root,omitempty"            swaggertype:"string"`
	MerkleTreeAddress   string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"merkle_tree_address,omitempty"   swaggertype:"string"`
	StateTransitionVKey string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"state_transition_vkey,omitempty" swaggertype:"string"`
	StateMembershipVKey string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"state_membership_vkey,omitempty" swaggertype:"string"`
	Groth16VKey         string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"groth16_vkey,omitempty"          swaggertype:"string"`

	Creator *ShortAddress `json:"creator,omitempty"`
}

func NewZkISM(ism storage.ZkISM) ZkISM {
	result := ZkISM{
		Id:      ism.Id,
		Height:  ism.Height,
		Time:    ism.Time,
		Creator: NewShortAddress(ism.Creator),
	}

	if len(ism.ExternalId) > 0 {
		result.ExternalId = hex.EncodeToString(ism.ExternalId)
	}
	if len(ism.State) > 0 {
		result.State = hex.EncodeToString(ism.State)
	}
	if len(ism.StateRoot) > 0 {
		result.StateRoot = hex.EncodeToString(ism.StateRoot)
	}
	if len(ism.MerkleTreeAddress) > 0 {
		result.MerkleTreeAddress = hex.EncodeToString(ism.MerkleTreeAddress)
	}
	if len(ism.StateTransitionVKey) > 0 {
		result.StateTransitionVKey = hex.EncodeToString(ism.StateTransitionVKey)
	}
	if len(ism.StateMembershipVKey) > 0 {
		result.StateMembershipVKey = hex.EncodeToString(ism.StateMembershipVKey)
	}
	if len(ism.Groth16VKey) > 0 {
		result.Groth16VKey = hex.EncodeToString(ism.Groth16VKey)
	}
	if ism.Tx != nil {
		result.TxHash = hex.EncodeToString(ism.Tx.Hash)
	}

	return result
}

type ZkISMUpdate struct {
	Id           uint64         `example:"321"                                                              format:"int64"     json:"id"                       swaggertype:"integer"`
	Height       pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"                   swaggertype:"integer"`
	Time         time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                     swaggertype:"string"`
	TxHash       string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"        swaggertype:"string"`
	NewStateRoot string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"new_state_root,omitempty" swaggertype:"string"`
	NewState     string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"new_state,omitempty"      swaggertype:"string"`

	Signer *ShortAddress `json:"signer,omitempty"`
}

func NewZkISMUpdate(update storage.ZkISMUpdate) ZkISMUpdate {
	result := ZkISMUpdate{
		Id:     update.Id,
		Height: update.Height,
		Time:   update.Time,
		Signer: NewShortAddress(update.Signer),
	}

	if len(update.NewState) > 0 {
		result.NewState = hex.EncodeToString(update.NewState)
	}
	if len(update.NewStateRoot) > 0 {
		result.NewStateRoot = hex.EncodeToString(update.NewStateRoot)
	}
	if update.Tx != nil {
		result.TxHash = hex.EncodeToString(update.Tx.Hash)
	}

	return result
}

type ZkISMMessage struct {
	Id        uint64         `example:"321"                                                              format:"int64"     json:"id"                   swaggertype:"integer"`
	Height    pkgTypes.Level `example:"100"                                                              format:"int64"     json:"height"               swaggertype:"integer"`
	Time      time.Time      `example:"2023-07-04T03:10:57+00:00"                                        format:"date-time" json:"time"                 swaggertype:"string"`
	TxHash    string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"tx_hash,omitempty"    swaggertype:"string"`
	StateRoot string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"state_root,omitempty" swaggertype:"string"`
	MessageId string         `example:"652452A670018D629CC116E510BA88C1CABE061336661B1F3D206D248BD558AF" format:"binary"    json:"message_id,omitempty" swaggertype:"string"`

	Signer *ShortAddress `json:"signer,omitempty"`
}

func NewZkISMMessage(msg storage.ZkISMMessage) ZkISMMessage {
	result := ZkISMMessage{
		Id:     msg.Id,
		Height: msg.Height,
		Time:   msg.Time,
		Signer: NewShortAddress(msg.Signer),
	}

	if len(msg.StateRoot) > 0 {
		result.StateRoot = hex.EncodeToString(msg.StateRoot)
	}
	if len(msg.MessageId) > 0 {
		result.MessageId = hex.EncodeToString(msg.MessageId)
	}
	if msg.Tx != nil {
		result.TxHash = hex.EncodeToString(msg.Tx.Hash)
	}

	return result
}

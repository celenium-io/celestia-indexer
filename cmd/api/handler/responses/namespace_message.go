// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"time"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/pkg/errors"
)

type NamespaceMessage struct {
	Id       uint64    `example:"321"                       format:"int64"     json:"id"       swaggertype:"integer"`
	Height   int64     `example:"100"                       format:"int64"     json:"height"   swaggertype:"integer"`
	Time     time.Time `example:"2023-07-04T03:10:57+00:00" format:"date-time" json:"time"     swaggertype:"string"`
	Position int64     `example:"2"                         format:"int64"     json:"position" swaggertype:"integer"`

	Type string `enums:"MsgWithdrawValidatorCommission,MsgWithdrawDelegatorReward,MsgEditValidator,MsgBeginRedelegate,MsgCreateValidator,MsgDelegate,MsgUndelegate,MsgUnjail,MsgSend,MsgCreateVestingAccount,MsgCreatePeriodicVestingAccount,MsgPayForBlobs,MsgGrantAllowance" example:"MsgCreatePeriodicVestingAccount" format:"string" json:"type" swaggertype:"string"`

	Data      map[string]any `json:"data"`
	Tx        Tx             `json:"tx"`
	Namespace Namespace      `json:"namespace"`
}

func NewNamespaceMessage(msg storage.NamespaceMessage) (NamespaceMessage, error) {
	if msg.Message == nil {
		return NamespaceMessage{}, errors.New("nil message in namespace message constructor")
	}
	if msg.Tx == nil {
		return NamespaceMessage{}, errors.New("nil tx in namespace message constructor")
	}
	if msg.Namespace == nil {
		return NamespaceMessage{}, errors.New("nil namespace in namespace message constructor")
	}

	return NamespaceMessage{
		Id:        msg.Message.Id,
		Height:    int64(msg.Message.Height),
		Time:      msg.Message.Time,
		Position:  msg.Message.Position,
		Type:      string(msg.Message.Type),
		Data:      msg.Message.Data,
		Tx:        NewTx(*msg.Tx),
		Namespace: NewNamespace(*msg.Namespace),
	}, nil
}

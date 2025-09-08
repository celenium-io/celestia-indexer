// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package ibc_relayer

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
)

//go:generate mockgen -source=$GOFILE -destination=mock.go -package=ibc_relayer -typed
type IRelayerStore interface {
	List() (metadata map[uint64]responses.Relayer)
	All() (relayers []responses.Relayer)
}

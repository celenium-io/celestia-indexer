// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"

	"github.com/celestiaorg/go-square/namespace"
	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/rsmt2d"
)

type ODS struct {
	Width uint `example:"2" json:"width" swaggertype:"integer"`

	Items []ODSItem `json:"items"`
}

type ODSItem struct {
	From      []uint        `json:"from"`
	To        []uint        `json:"to"`
	Namespace string        `json:"namespace"`
	Type      NamespaceKind `json:"type"`
}

func NewODS(eds *rsmt2d.ExtendedDataSquare) (ODS, error) {
	ods := ODS{
		Width: eds.Width() / 2,
		Items: make([]ODSItem, 0),
	}

	var current ODSItem
	for i := uint(0); i < ods.Width; i++ {
		for j := uint(0); j < ods.Width; j++ {
			cell := eds.GetCell(i, j)
			share, err := shares.NewShare(cell)
			if err != nil {
				return ods, err
			}
			namespace, err := share.Namespace()
			if err != nil {
				return ods, err
			}
			base64Namespace := base64.StdEncoding.EncodeToString(namespace.Bytes())
			if base64Namespace != current.Namespace {
				if current.Namespace != "" {
					ods.Items = append(ods.Items, current)
				}
				current = ODSItem{
					From:      []uint{i, j},
					Namespace: base64Namespace,
					Type:      getNamespaceType(namespace),
				}
			}
			current.To = []uint{i, j}
		}
	}
	ods.Items = append(ods.Items, current)

	return ods, nil
}

type NamespaceKind string

const (
	PayForBlobNamespace      NamespaceKind = "pay_for_blob"
	TailPaddingNamespace     NamespaceKind = "tail_padding"
	TxNamespace              NamespaceKind = "tx"
	ParitySharesNamespace    NamespaceKind = "parity_shares"
	PrimaryReservedNamespace NamespaceKind = "primary_reserved_padding"
	DefaultNamespace         NamespaceKind = "namespace"
)

func getNamespaceType(ns namespace.Namespace) NamespaceKind {
	switch {
	case ns.IsPayForBlob():
		return PayForBlobNamespace
	case ns.IsTailPadding():
		return TailPaddingNamespace
	case ns.IsTx():
		return TxNamespace
	case ns.IsParityShares():
		return ParitySharesNamespace
	case ns.IsPrimaryReservedPadding():
		return PrimaryReservedNamespace
	default:
		return DefaultNamespace
	}
}

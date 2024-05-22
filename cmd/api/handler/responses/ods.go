// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"

	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/rsmt2d"
)

type ODS struct {
	Width uint `example:"2" json:"width" swaggertype:"integer"`

	Items []ODSItem `json:"items"`
}

type ODSItem struct {
	From      []uint `json:"from"`
	To        []uint `json:"to"`
	Namespace string `json:"namespace"`
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
				}
			}
			current.To = []uint{i, j}
		}
	}
	ods.Items = append(ods.Items, current)

	return ods, nil
}

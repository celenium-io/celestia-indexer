// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"encoding/base64"
	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	"github.com/celestiaorg/go-square/namespace"
	"github.com/celestiaorg/go-square/shares"
	incl "github.com/celestiaorg/go-square/v2/inclusion"
	"github.com/celestiaorg/go-square/v2/share"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/merkle"

	_ "github.com/celestiaorg/go-square/v2/share"
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

type sequence struct {
	ns            share.Namespace
	shareVersion  uint8
	startShareIdx int
	endShareIdx   int
	data          []byte
	sequenceLen   uint32
	signer        []byte
}

func (ods *ODS) FindODSByNamespace(namespace string) (*ODSItem, error) {
	for _, item := range ods.Items {
		if item.Namespace == namespace {
			return &item, nil
		}
	}
	return nil, errors.New("item with specified namespace not found")
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

func GetNamespaceShares(eds *rsmt2d.ExtendedDataSquare, from, to []uint) ([]share.Share, error) {
	if len(from) != len(to) {
		return nil, errors.New("length of 'from' and 'to' must match")
	}

	var resultShares []share.Share
	startRow, startCol := from[0], from[1]
	endRow, endCol := to[0], to[1]

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		return nil, errors.New("invalid from and to params")
	}
	currentRow, currentCol := startRow, startCol

	for {
		cell := eds.GetCell(currentRow, currentCol)
		cellShare, err := share.NewShare(cell)
		if err != nil {
			return nil, err
		}
		resultShares = append(resultShares, *cellShare)
		if currentRow == endRow && currentCol == endCol {
			break
		}
		currentCol++
		if currentCol == eds.Width()/2 {
			currentCol = 0
			currentRow++
		}
	}

	return resultShares, nil
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

func GetBlobShareIdxs(
	shares []share.Share,
	nsStartFromIdx []uint,
	edsWidth uint,
	b64commitment string,
) (blobStartIdx, blobEndIdx int, err error) {
	if len(shares) == 0 {
		return 0, 0, errors.New("invalid shares length")

	}
	sequences := make([]sequence, 0)
	startRow, startCol := nsStartFromIdx[0], nsStartFromIdx[1]
	nsStartIdx := int(startRow*edsWidth + startCol)

	for shareIdx, s := range shares {
		if !(s.Version() <= 1) {
			return 0, 0, errors.New("unsupported share version")
		}

		if s.IsPadding() {
			continue
		}

		if s.IsSequenceStart() {
			sequences = append(sequences, sequence{
				ns:            s.Namespace(),
				shareVersion:  s.Version(),
				startShareIdx: shareIdx,
				endShareIdx:   shareIdx,
				data:          s.RawData(),
				sequenceLen:   s.SequenceLen(),
				signer:        share.GetSigner(s),
			})
		} else {
			if len(sequences) == 0 {
				return 0, 0, errors.New("continuation share without a sequence start share")
			}
			prev := &sequences[len(sequences)-1]
			prev.data = append(prev.data, s.RawData()...)
			prev.endShareIdx = shareIdx
		}
	}
	for _, seq := range sequences {
		seq.data = seq.data[:seq.sequenceLen]
		blob, err := share.NewBlob(seq.ns, seq.data, seq.shareVersion, seq.signer)
		if err != nil {
			return 0, 0, err
		}
		commitment, err := incl.CreateCommitment(blob, merkle.HashFromByteSlices, appconsts.SubtreeRootThreshold(0))
		if err != nil {
			return 0, 0, err
		}
		if base64.StdEncoding.EncodeToString(commitment) == b64commitment {
			return nsStartIdx + seq.startShareIdx, nsStartIdx + seq.endShareIdx + 1, err
		}
	}
	return 0, 0, err
}

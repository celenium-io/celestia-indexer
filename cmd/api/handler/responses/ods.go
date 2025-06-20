// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package responses

import (
	"bytes"
	"encoding/base64"

	"github.com/celestiaorg/celestia-app/v4/pkg/appconsts"
	"github.com/celestiaorg/go-square/namespace"
	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/go-square/v2/inclusion"
	"github.com/celestiaorg/go-square/v2/share"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/pkg/errors"

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

type sequence struct {
	ns            share.Namespace
	shareVersion  uint8
	startShareIdx int
	endShareIdx   int
	data          []byte
	sequenceLen   uint32
	signer        []byte
}

func GetBlobShareIndexes(
	shares []share.Share,
	base64namespace string,
	base64commitment string,
) (blobStartIndex, blobEndIndex int, err error) {
	if len(shares) == 0 {
		return 0, 0, errors.New("invalid shares length")

	}
	sequences := make([]sequence, 0)
	namespaceBytes, err := base64.StdEncoding.DecodeString(base64namespace)
	if err != nil {
		return 0, 0, errors.Wrap(err, "decoding base64 namespace")
	}

	for shareIndex, s := range shares {
		if s.Version() > 1 {
			return 0, 0, errors.New("unsupported share version")
		}

		if s.IsPadding() {
			continue
		}

		if !bytes.Equal(s.Namespace().Bytes(), namespaceBytes) {
			continue
		}

		if s.IsSequenceStart() {
			sequences = append(sequences, sequence{
				ns:            s.Namespace(),
				shareVersion:  s.Version(),
				startShareIdx: shareIndex,
				endShareIdx:   shareIndex,
				data:          s.RawData(),
				sequenceLen:   s.SequenceLen(),
				signer:        share.GetSigner(s),
			})
		} else {
			if len(sequences) == 0 {
				return 0, 0, errors.New("continuation share without a s start share")
			}
			prev := &sequences[len(sequences)-1]
			prev.data = append(prev.data, s.RawData()...)
			prev.endShareIdx = shareIndex
		}
	}
	for _, s := range sequences {
		s.data = s.data[:s.sequenceLen]
		blob, err := share.NewBlob(s.ns, s.data, s.shareVersion, s.signer)
		if err != nil {
			return 0, 0, errors.Wrap(err, "creating blob")
		}
		commitment, err := inclusion.CreateCommitment(blob, merkle.HashFromByteSlices, appconsts.SubtreeRootThreshold)
		if err != nil {
			return 0, 0, errors.Wrap(err, "creating commitment")
		}
		if base64.StdEncoding.EncodeToString(commitment) == base64commitment {
			return s.startShareIdx, s.endShareIdx + 1, err
		}
	}
	return 0, 0, err
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package hyperlane

import "context"

type ChainMetadata struct {
	DomainId       uint64          `yaml:"domainId"`
	DisplayName    string          `yaml:"displayName"`
	BlockExplorers []BlockExplorer `yaml:"blockExplorers"`
	NativeToken    NativeToken     `yaml:"nativeToken"`
}

type BlockExplorer struct {
	ApiUrl string `yaml:"apiUrl"`
	Family string `yaml:"family"`
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
}

type NativeToken struct {
	Decimals uint64 `yaml:"decimals"`
	Name     string `yaml:"name"`
	Symbol   string `yaml:"symbol"`
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IApi interface {
	ChainMetadata(ctx context.Context) (map[uint64]ChainMetadata, error)
}

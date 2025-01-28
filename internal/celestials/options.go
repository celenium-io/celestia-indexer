// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package celestials

type ChangeOptions struct {
	Limit        int64  `json:"limit"`
	FromChangeId int64  `json:"from_change_id"`
	OnlyHead     bool   `json:"only_head"`
	Images       bool   `json:"with_images"`
	ChainId      string `json:"chain_id"`
}

type ChangeOption func(opts *ChangeOptions)

func WithLimit(limit int64) ChangeOption {
	return func(opts *ChangeOptions) {
		if limit > 0 {
			opts.Limit = limit
		}
	}
}

func WithOnlyHead() ChangeOption {
	return func(opts *ChangeOptions) {
		opts.OnlyHead = true
	}
}

func WithImages() ChangeOption {
	return func(opts *ChangeOptions) {
		opts.Images = true
	}
}

func WithFromChangeId(changeId int64) ChangeOption {
	return func(opts *ChangeOptions) {
		if changeId > 0 {
			opts.FromChangeId = changeId
		}
	}
}

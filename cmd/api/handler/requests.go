// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"strings"

	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

const (
	asc  = "asc"
	desc = "desc"
)

type limitOffsetPagination struct {
	Limit  int    `json:"limit"  param:"limit"  query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" param:"offset" query:"offset" validate:"omitempty,min=0"`
	Sort   string `json:"sort"   param:"sort"   query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *limitOffsetPagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

func pgSort(sort string) storage.SortOrder {
	switch sort {
	case asc:
		return storage.SortOrderAsc
	case desc:
		return storage.SortOrderDesc
	default:
		return storage.SortOrderAsc
	}
}

type txListRequest struct {
	Limit           int         `query:"limit"             validate:"omitempty,min=1,max=100"`
	Offset          int         `query:"offset"            validate:"omitempty,min=0"`
	Sort            string      `query:"sort"              validate:"omitempty,oneof=asc desc"`
	Height          *uint64     `query:"height"            validate:"omitempty,min=0"`
	Status          StringArray `query:"status"            validate:"omitempty,dive,status"`
	MsgType         StringArray `query:"msg_type"          validate:"omitempty,dive,msg_type"`
	ExcludedMsgType StringArray `query:"excluded_msg_type" validate:"omitempty,dive,msg_type"`
	Messages        bool        `query:"messages"          validate:"omitempty"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (p *txListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

type StringArray []string

func (s *StringArray) UnmarshalParam(param string) error {
	*s = StringArray(strings.Split(param, ","))
	return nil
}

type StatusArray StringArray
type MsgTypeArray StringArray

func bindAndValidate[T any](c echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return req, err
	}
	if err := c.Validate(req); err != nil {
		return req, err
	}
	return req, nil
}

type addressTxRequest struct {
	Hash    string      `param:"hash"     validate:"required,address"`
	Limit   int         `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  int         `query:"offset"   validate:"omitempty,min=0"`
	Sort    string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	Height  *uint64     `query:"height"   validate:"omitempty,min=0"`
	Status  StringArray `query:"status"   validate:"omitempty,dive,status"`
	MsgType StringArray `query:"msg_type" validate:"omitempty,dive,msg_type"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (p *addressTxRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

type listMessageByBlockRequest struct {
	Height          pkgTypes.Level `param:"height"            validate:"required,min=1"`
	Limit           int            `query:"limit"             validate:"omitempty,min=1,max=100"`
	Offset          int            `query:"offset"            validate:"omitempty,min=0"`
	MsgType         StringArray    `query:"msg_type"          validate:"omitempty,dive,msg_type"`
	ExcludedMsgType StringArray    `query:"excluded_msg_type" validate:"omitempty,dive,msg_type"`
}

func (p *listMessageByBlockRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

type namespaceList struct {
	Limit  int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=time pfb_count size"`
}

func (p *namespaceList) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
}

type getById struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

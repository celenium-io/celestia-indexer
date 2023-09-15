package handler

import (
	"strings"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

const (
	asc  = "asc"
	desc = "desc"
)

type limitOffsetPagination struct {
	Limit  uint64 `json:"limit"  param:"limit"  query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `json:"offset" param:"offset" query:"offset" validate:"omitempty,min=0"`
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
	Limit   uint64      `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  uint64      `query:"offset"   validate:"omitempty,min=0"`
	Sort    string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	Height  uint64      `query:"height"   validate:"omitempty,min=1"`
	Status  StringArray `query:"status"   validate:"omitempty,dive,status"`
	MsgType StringArray `query:"msg_type" validate:"omitempty,dive,msg_type"`

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
	Limit   uint64      `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  uint64      `query:"offset"   validate:"omitempty,min=0"`
	Sort    string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	Height  uint64      `query:"height"   validate:"omitempty,min=1"`
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

type namespacesByHeightRequest struct {
	Limit  uint64 `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset" validate:"omitempty,min=0"`
	Height uint64 `param:"height" validate:"required,min=1"`
}

func (p *namespacesByHeightRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

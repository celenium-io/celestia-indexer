// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type SearchHandler struct {
	search     storage.ISearch
	address    storage.IAddress
	block      storage.IBlock
	tx         storage.ITx
	namespace  storage.INamespace
	validator  storage.IValidator
	rollup     storage.IRollup
	celestials celestials.ICelestial
}

func NewSearchHandler(
	search storage.ISearch,
	address storage.IAddress,
	block storage.IBlock,
	tx storage.ITx,
	namespace storage.INamespace,
	validator storage.IValidator,
	rollup storage.IRollup,
	celestials celestials.ICelestial,
) SearchHandler {
	return SearchHandler{
		search:     search,
		address:    address,
		block:      block,
		tx:         tx,
		namespace:  namespace,
		validator:  validator,
		rollup:     rollup,
		celestials: celestials,
	}
}

type searchRequest struct {
	Search string `query:"query" validate:"required"`
}

var (
	hashRegexp      = regexp.MustCompile("^(0x)?[a-fA-F0-9]{64}$")
	namespaceRegexp = regexp.MustCompile("^[a-fA-f0-9]{58}$")
)

// Search godoc
//
//	@Summary				Search by hash
//	@Description.markdown	search
//	@Tags					search
//	@ID						search
//	@Param					query	query	string	true	"Search string"
//	@Produce				json
//	@Success				200	{array}	responses.SearchItem
//	@Success				204
//	@Failure				400	{object}	Error
//	@Failure				500	{object}	Error
//	@Router					/search [get]
func (handler SearchHandler) Search(c echo.Context) error {
	req, err := bindAndValidate[searchRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	data := make([]responses.SearchItem, 0)

	if height, err := strconv.ParseInt(req.Search, 10, 64); err == nil {
		block, err := handler.block.ByHeight(c.Request().Context(), types.Level(height))
		if err == nil {
			data = append(data, responses.SearchItem{
				Type:   "block",
				Result: responses.NewBlock(block, false),
			})
		}
	}

	var response []responses.SearchItem

	switch {
	case isAddress(req.Search):
		response, err = handler.searchAddress(c.Request().Context(), req.Search)
	case isValoperAddress(req.Search):
		response, err = handler.searchValoperAddress(c.Request().Context(), req.Search)
	case hashRegexp.MatchString(req.Search):
		response, err = handler.searchHash(c.Request().Context(), req.Search)
	case namespaceRegexp.MatchString(req.Search):
		response, err = handler.searchNamespaceById(c.Request().Context(), req.Search)
	case isNamespace(req.Search):
		response, err = handler.searchNamespaceByBase64(c.Request().Context(), req.Search)
	default:
		response, err = handler.searchText(c.Request().Context(), req.Search)
	}
	if err != nil {
		if !handler.address.IsNoRows(err) {
			return handleError(c, err, handler.address)
		}
	}

	data = append(data, response...)
	return returnArray(c, data)
}

func (handler SearchHandler) searchAddress(ctx context.Context, search string) ([]responses.SearchItem, error) {
	_, hash, err := types.Address(search).Decode()
	if err != nil {
		return nil, err
	}

	address, err := handler.address.ByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	result := responses.SearchItem{
		Type:   "address",
		Result: responses.NewAddress(address),
	}
	return []responses.SearchItem{result}, nil
}

func (handler SearchHandler) searchValoperAddress(ctx context.Context, search string) ([]responses.SearchItem, error) {
	_, hash, err := types.Address(search).Decode()
	if err != nil {
		return nil, err
	}

	address, err := handler.address.ByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	validator, err := handler.validator.ByAddress(ctx, search)
	if err != nil {
		return nil, err
	}

	return []responses.SearchItem{
		{
			Type:   "validator",
			Result: responses.NewValidator(validator),
		}, {
			Type:   "address",
			Result: responses.NewAddress(address),
		},
	}, nil
}

func (handler SearchHandler) searchHash(ctx context.Context, search string) ([]responses.SearchItem, error) {
	search = strings.TrimPrefix(search, "0x")
	data, err := hex.DecodeString(search)
	if err != nil {
		return nil, err
	}
	if len(data) != 32 {
		return nil, errors.Wrapf(errInvalidHashLength, "got %d", len(data))
	}
	result, err := handler.search.Search(ctx, data)
	if err != nil {
		return nil, err
	}

	response := make([]responses.SearchItem, len(result))
	for i := range result {
		response[i].Type = result[i].Type
		switch response[i].Type {
		case "tx":
			tx, err := handler.tx.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			response[i].Result = responses.NewTx(*tx)
		case "block":
			block, err := handler.block.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			response[i].Result = responses.NewBlock(*block, false)
		}
	}

	return response, nil
}

func (handler SearchHandler) searchNamespaceById(ctx context.Context, search string) ([]responses.SearchItem, error) {
	data, err := hex.DecodeString(search)
	if err != nil {
		return nil, err
	}

	return handler.getNamespace(ctx, data)
}

func (handler SearchHandler) searchNamespaceByBase64(ctx context.Context, search string) ([]responses.SearchItem, error) {
	data, err := base64.StdEncoding.DecodeString(search)
	if err != nil {
		return nil, err
	}

	return handler.getNamespace(ctx, data)
}

func (handler SearchHandler) getNamespace(ctx context.Context, data []byte) ([]responses.SearchItem, error) {
	version := data[0]
	namespaceId := data[1:]
	ns, err := handler.namespace.ByNamespaceIdAndVersion(ctx, namespaceId, version)
	if err != nil {
		return nil, err
	}
	result := responses.SearchItem{
		Type:   "namespace",
		Result: responses.NewNamespace(ns),
	}
	return []responses.SearchItem{result}, nil
}

func (handler SearchHandler) searchText(ctx context.Context, text string) ([]responses.SearchItem, error) {
	result, err := handler.search.SearchText(ctx, text)
	if err != nil {
		return nil, err
	}

	response := make([]responses.SearchItem, len(result))
	for i := range result {
		response[i].Type = result[i].Type
		switch response[i].Type {
		case "validator":
			validator, err := handler.validator.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			response[i].Result = responses.NewValidator(*validator)
		case "rollup":
			rollup, err := handler.rollup.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			response[i].Result = responses.NewRollup(rollup)
		case "namespace":
			namespace, err := handler.namespace.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			response[i].Result = responses.NewNamespace(*namespace)
		case "celestial":
			address, err := handler.address.GetByID(ctx, result[i].Id)
			if err != nil {
				return nil, err
			}
			addr := responses.NewAddress(*address)

			celestial, err := handler.celestials.ById(ctx, result[i].Value)
			if err != nil {
				return nil, err
			}
			addr.AddCelestails(&celestial)

			response[i].Result = addr
			response[i].Type = "address"
		default:
			return nil, errors.Errorf("unknown search text type: %s", response[i].Type)
		}
	}

	return response, nil
}

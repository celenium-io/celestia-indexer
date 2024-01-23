// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"regexp"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type SearchHandler struct {
	address   storage.IAddress
	block     storage.IBlock
	namespace storage.INamespace
	tx        storage.ITx
}

func NewSearchHandler(
	address storage.IAddress,
	block storage.IBlock,
	namespace storage.INamespace,
	tx storage.ITx,
) SearchHandler {
	return SearchHandler{
		address:   address,
		block:     block,
		namespace: namespace,
		tx:        tx,
	}
}

type searchRequest struct {
	Search string `query:"query" validate:"required"`
}

var (
	hashRegexp      = regexp.MustCompile("[a-fA-f0-9]{64}")
	namespaceRegexp = regexp.MustCompile("[a-fA-f0-9]{58}")
)

// Search godoc
//
//	@Summary				Search by hash
//	@Description.markdown	search
//	@Tags					search
//	@ID						search
//	@Param					query	query	string	true	"Search string"
//	@Produce				json
//	@Success				200	{object}	responses.SearchResponse[responses.Searchable]
//	@Success				204
//	@Failure				400	{object}	Error
//	@Failure				500	{object}	Error
//	@Router					/v1/search [get]
func (handler SearchHandler) Search(c echo.Context) error {
	req, err := bindAndValidate[searchRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	switch {
	case isAddress(req.Search):
		if err := handler.searchAddress(c, req.Search); err != nil {
			return internalServerError(c, err)
		}
	case hashRegexp.MatchString(req.Search):
		if err := handler.searchHash(c, req.Search); err != nil {
			return internalServerError(c, err)
		}
	case namespaceRegexp.MatchString(req.Search):
		if err := handler.searchNamespaceById(c, req.Search); err != nil {
			return internalServerError(c, err)
		}
	case isNamespace(req.Search):
		if err := handler.searchNamespaceByBase64(c, req.Search); err != nil {
			return internalServerError(c, err)
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler SearchHandler) searchAddress(c echo.Context, search string) error {
	_, hash, err := types.Address(search).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	return c.JSON(http.StatusOK, responses.NewSearchResponse(responses.NewAddress(address)))
}

func (handler SearchHandler) searchHash(c echo.Context, search string) error {
	data, err := hex.DecodeString(search)
	if err != nil {
		return badRequestError(c, err)
	}
	if len(data) != 32 {
		return badRequestError(c, errors.Wrapf(errInvalidHashLength, "got %d", len(data)))
	}
	tx, err := handler.tx.ByHash(c.Request().Context(), data)
	if err == nil {
		return c.JSON(http.StatusOK, responses.NewSearchResponse(responses.NewTx(tx)))
	}

	if !handler.tx.IsNoRows(err) {
		return internalServerError(c, err)
	}

	block, err := handler.block.ByHash(c.Request().Context(), data)
	if err == nil {
		return c.JSON(http.StatusOK, responses.NewSearchResponse(responses.NewBlock(block, false)))
	}
	if !handler.tx.IsNoRows(err) {
		return internalServerError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler SearchHandler) searchNamespaceById(c echo.Context, search string) error {
	data, err := hex.DecodeString(search)
	if err != nil {
		return badRequestError(c, err)
	}

	return handler.getNamespace(c, data)
}

func (handler SearchHandler) searchNamespaceByBase64(c echo.Context, search string) error {
	data, err := base64.StdEncoding.DecodeString(search)
	if err != nil {
		return badRequestError(c, err)
	}

	return handler.getNamespace(c, data)
}

func (handler SearchHandler) getNamespace(c echo.Context, data []byte) error {
	version := data[0]
	namespaceId := data[1:]
	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, version)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}
	response := responses.NewNamespace(ns)
	return c.JSON(http.StatusOK, responses.NewSearchResponse(response))
}

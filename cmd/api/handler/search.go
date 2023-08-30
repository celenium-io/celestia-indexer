package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"regexp"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
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
	base64Regexp    = regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")
)

// Get godoc
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
	case base64Regexp.MatchString(req.Search):
		if err := handler.searchNamespaceByBase64(c, req.Search); err != nil {
			return internalServerError(c, err)
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler SearchHandler) searchAddress(c echo.Context, search string) error {
	data, err := responses.DecodeAddress(search)
	if err != nil {
		return badRequestError(c, err)
	}
	address, err := handler.address.ByHash(c.Request().Context(), data)
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}
	response, err := responses.NewAddress(address)
	if err != nil {
		return internalServerError(c, err)
	}
	return c.JSON(http.StatusOK, responses.NewSearchResponse(response))
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
		return c.JSON(http.StatusOK, responses.NewSearchResponse(responses.NewBlock(block)))
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
	if len(data) != 29 {
		return badRequestError(c, errors.Wrapf(errInvalidNamespaceLength, "got %d", len(data)))
	}

	return handler.getNamespace(c, data)
}

func (handler SearchHandler) searchNamespaceByBase64(c echo.Context, search string) error {
	data, err := base64.URLEncoding.DecodeString(search)
	if err != nil {
		return badRequestError(c, err)
	}
	if len(data) != 29 {
		return badRequestError(c, errors.Wrapf(errInvalidNamespaceLength, "got %d", len(data)))
	}

	return handler.getNamespace(c, data)
}

func (handler SearchHandler) getNamespace(c echo.Context, data []byte) error {
	version := data[0]
	namespaceId := data[1:]
	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}
	response := responses.NewNamespace(ns)
	return c.JSON(http.StatusOK, responses.NewSearchResponse(response))
}

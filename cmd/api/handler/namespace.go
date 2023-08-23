package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type NamespaceHandler struct {
	namespace storage.INamespace
}

func NewNamespaceHandler(namespace storage.INamespace) *NamespaceHandler {
	return &NamespaceHandler{
		namespace: namespace,
	}
}

type getNamespaceRequest struct {
	Id string `param:"id" validate:"required,hexadecimal,len=56"`
}

// Get godoc
// @Summary Get namespace info
// @Description Returns array of namespace versions
// @Tags namespace
// @ID get-namespace
// @Param id path string true "Namespace id in hexadecimal" minlength(56) maxlength(56)
// @Produce  json
// @Success 200 {array} Namespace
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/namespace/{id} [get]
func (handler *NamespaceHandler) Get(c echo.Context) error {
	req := new(getNamespaceRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	namespace, err := handler.namespace.ByNamespaceId(c.Request().Context(), namespaceId)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	response := make([]responses.Namespace, len(namespace))
	for i := range namespace {
		response[i] = responses.NewNamespace(namespace[i])
	}

	return returnArray(c, response)
}

type getNamespaceByHashRequest struct {
	Hash string `param:"hash" validate:"required,base64"`
}

// GetByHash godoc
// @Summary Get namespace info by base64
// @Description Returns namespace by base64 encoded identity
// @Tags namespace
// @ID get-namespace-base64
// @Param hash path string true "Base64-encoded namespace id and version"
// @Produce  json
// @Success 200 {object} Namespace
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/namespace_by_hash/{hash} [get]
func (handler *NamespaceHandler) GetByHash(c echo.Context) error {
	req := new(getNamespaceByHashRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.URLEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}
	if len(hash) != 29 {
		return badRequestError(c, errors.Wrapf(errInvalidNamespaceLength, "got %d", len(hash)))
	}
	version := hash[0]
	namespaceId := hash[1:]

	namespace, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responses.NewNamespace(namespace))
}

type getNamespaceWithVersionRequest struct {
	Id      string `param:"id"      validate:"required,hexadecimal,len=56"`
	Version byte   `param:"version" validate:"required"`
}

// GetWithVersion godoc
// @Summary Get namespace info by id and version
// @Description Returns namespace by version byte and namespace id
// @Tags namespace
// @ID get-namespace-by-version-and-id
// @Param id      path string  true "Namespace id in hexadecimal" minlength(56) maxlength(56)
// @Param version path integer true "Version of namespace"
// @Produce  json
// @Success 200 {object} Namespace
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/namespace/{id}/{version} [get]
func (handler *NamespaceHandler) GetWithVersion(c echo.Context) error {
	req := new(getNamespaceWithVersionRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	namespace, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.NewNamespace(namespace))
}

// List godoc
// @Summary List namespace info
// @Description List namespace info
// @Tags namespace
// @ID list-namespace
// @Param limit  query integer false "Count of requested entities" mininum(1) maximum(100)
// @Param offset query integer false "Offset" mininum(1)
// @Param sort   query string  false "Sort order" Enums(asc, desc)
// @Produce json
// @Success 200 {array} Block
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/namespace [get]
func (handler *NamespaceHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespace, err := handler.namespace.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}
	response := make([]responses.Namespace, len(namespace))
	for i := range namespace {
		response[i] = responses.NewNamespace(*namespace[i])
	}
	return returnArray(c, response)
}

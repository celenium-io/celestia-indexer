package handler

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/pkg/node"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type NamespaceHandler struct {
	namespace   storage.INamespace
	blob        node.CelestiaNodeApi
	state       storage.IState
	indexerName string
}

func NewNamespaceHandler(
	namespace storage.INamespace,
	state storage.IState,
	indexerName string,
	blob node.CelestiaNodeApi,
) *NamespaceHandler {
	return &NamespaceHandler{
		namespace:   namespace,
		blob:        blob,
		state:       state,
		indexerName: indexerName,
	}
}

type getNamespaceRequest struct {
	Id string `param:"id" validate:"required,hexadecimal,len=56"`
}

// Get godoc
//
//	@Summary		Get namespace info
//	@Description	Returns array of namespace versions
//	@Tags			namespace
//	@ID				get-namespace
//	@Param			id	path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Produce		json
//	@Success		200	{array}		responses.Namespace
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace/{id} [get]
func (handler *NamespaceHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getNamespaceRequest](c)
	if err != nil {
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
//
//	@Summary		Get namespace info by base64
//	@Description	Returns namespace by base64 encoded identity
//	@Tags			namespace
//	@ID				get-namespace-base64
//	@Param			hash	path	string	true	"Base64-encoded namespace id and version"
//	@Produce		json
//	@Success		200	{object}	responses.Namespace
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace_by_hash/{hash} [get]
func (handler *NamespaceHandler) GetByHash(c echo.Context) error {
	req, err := bindAndValidate[getNamespaceByHashRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.StdEncoding.DecodeString(req.Hash)
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
	Version byte   `param:"version"`
}

// GetWithVersion godoc
//
//	@Summary		Get namespace info by id and version
//	@Description	Returns namespace by version byte and namespace id
//	@Tags			namespace
//	@ID				get-namespace-by-version-and-id
//	@Param			id		path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Param			version	path	integer	true	"Version of namespace"
//	@Produce		json
//	@Success		200	{object}	responses.Namespace
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace/{id}/{version} [get]
func (handler *NamespaceHandler) GetWithVersion(c echo.Context) error {
	req, err := bindAndValidate[getNamespaceWithVersionRequest](c)
	if err != nil {
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
//
//	@Summary		List namespace info
//	@Description	List namespace info
//	@Tags			namespace
//	@ID				list-namespace
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Namespace
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace [get]
func (handler *NamespaceHandler) List(c echo.Context) error {
	req, err := bindAndValidate[limitOffsetPagination](c)
	if err != nil {
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

type getBlobsRequest struct {
	Hash   string `param:"hash"   validate:"required,base64"`
	Height uint64 `param:"height" validation:"required,min=1"`
}

// GetBlobs godoc
//
//	@Summary		Get namespace blobs on height
//	@Description	Returns blobs
//	@Tags			namespace
//	@ID				get-namespace-blobs
//	@Param			hash	path	string	true	"Base64-encoded namespace id and version"
//	@Param			height	path	integer	true	"Block heigth"	minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Blob
//	@Failure		400	{object}	Error
//	@Router			/v1/namespace_by_hash/{hash}/{height} [get]
func (handler *NamespaceHandler) GetBlobs(c echo.Context) error {
	req, err := bindAndValidate[getBlobsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	blobs, err := handler.blob.Blobs(c.Request().Context(), req.Height, req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	return c.JSON(http.StatusOK, blobs)
}

type getBlobRequest struct {
	Hash       string `param:"hash"       validate:"required,base64"`
	Height     uint64 `param:"height"     validation:"required,min=1"`
	Commitment string `param:"commitment" validate:"required,base64"`
}

// GetBlob godoc
//
//	@Summary		Get namespace blob by commitment on height
//	@Description	Returns blob
//	@Tags			namespace
//	@ID				get-namespace-blob
//	@Param			hash		path	string	true	"Base64-encoded namespace id and version"
//	@Param			height		path	integer	true	"Block heigth"	minimum(1)
//	@Param			commitment	path	string	true	"Blob commitment"
//	@Produce		json
//	@Success		200	{object}	responses.Blob
//	@Failure		400	{object}	Error
//	@Router			/v1/namespace_by_hash/{hash}/{height}/{commitment} [get]
func (handler *NamespaceHandler) GetBlob(c echo.Context) error {
	req, err := bindAndValidate[getBlobRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	blob, err := handler.blob.Blob(c.Request().Context(), req.Height, req.Hash, req.Commitment)
	if err != nil {
		return badRequestError(c, err)
	}

	return c.JSON(http.StatusOK, blob)
}

type getNamespaceMessages struct {
	Id      string `param:"id"      validate:"required,hexadecimal,len=56"`
	Version byte   `param:"version"`
	Limit   uint64 `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  uint64 `query:"offset"  validate:"omitempty,min=0"`
}

// GetMessages godoc
//
//	@Summary		Get namespace messages by id and version
//	@Description	Returns namespace messages by version byte and namespace id
//	@Tags			namespace
//	@ID				get-namespace-messages
//	@Param			id		path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Param			version	path	integer	true	"Version of namespace"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}	responses.NamespaceMessage
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace/{id}/{version}/messages [get]
func (handler *NamespaceHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[getNamespaceMessages](c)
	if err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	messages, err := handler.namespace.Messages(c.Request().Context(), ns.Id, int(req.Limit), int(req.Offset))
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	response := make([]responses.NamespaceMessage, len(messages))
	for i := range response {
		msg, err := responses.NewNamespaceMessage(messages[i])
		if err := handleError(c, err, handler.namespace); err != nil {
			return err
		}
		response[i] = msg
	}

	return returnArray(c, response)
}

// GetActive godoc
//
//	@Summary		Get last used namespace
//	@Description	Get last used namespace
//	@Tags			namespace
//	@ID				get-namespace-active
//	@Produce		json
//	@Success		200	{array}		responses.ActiveNamespace
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace/active [get]
func (handler *NamespaceHandler) GetActive(c echo.Context) error {
	active, err := handler.namespace.Active(c.Request().Context(), 5)
	if err := handleError(c, err, handler.namespace); err != nil {
		return err
	}

	response := make([]responses.ActiveNamespace, len(active))
	for i := range response {
		response[i] = responses.NewActiveNamespace(active[i])
	}
	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of namespaces in network
//	@Description	Get count of namespaces in network
//	@Tags			namespace
//	@ID				get-namespace-count
//	@Produce		json
//	@Success		200	{integer}   uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/namespace/count [get]
func (handler *NamespaceHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err := handleError(c, err, handler.state); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, state.TotalNamespaces)
}

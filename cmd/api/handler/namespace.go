// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v3/pkg/da"
	"github.com/celestiaorg/go-square/v2/share"

	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/go-square/v2"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/node"
	"github.com/celestiaorg/celestia-app/v3/pkg/proof"
	"github.com/labstack/echo/v4"
)

type NamespaceHandler struct {
	namespace   storage.INamespace
	blobLogs    storage.IBlobLog
	rollups     storage.IRollup
	address     storage.IAddress
	blob        node.DalApi
	state       storage.IState
	node        node.Api
	indexerName string
}

func NewNamespaceHandler(
	namespace storage.INamespace,
	blobLogs storage.IBlobLog,
	rollups storage.IRollup,
	address storage.IAddress,
	state storage.IState,
	indexerName string,
	blob node.DalApi,
	node node.Api,
) *NamespaceHandler {
	return &NamespaceHandler{
		namespace:   namespace,
		blobLogs:    blobLogs,
		rollups:     rollups,
		address:     address,
		blob:        blob,
		state:       state,
		indexerName: indexerName,
		node:        node,
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
//	@Router			/namespace/{id} [get]
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
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	response := make([]responses.Namespace, len(namespace))
	for i := range namespace {
		response[i] = responses.NewNamespace(namespace[i])
	}

	return returnArray(c, response)
}

type getNamespaceByHashRequest struct {
	Hash string `param:"hash" validate:"required,base64,namespace"`
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
//	@Router			/namespace_by_hash/{hash} [get]
func (handler *NamespaceHandler) GetByHash(c echo.Context) error {
	req, err := bindAndValidate[getNamespaceByHashRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := base64.StdEncoding.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	version := hash[0]
	namespaceId := hash[1:]

	namespace, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, version)
	if err != nil {
		return handleError(c, err, handler.namespace)
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
//	@Router			/namespace/{id}/{version} [get]
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
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	return c.JSON(http.StatusOK, responses.NewNamespace(namespace))
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

// List godoc
//
//	@Summary		List namespace info
//	@Description	List namespace info
//	@Tags			namespace
//	@ID				list-namespace
//	@Param			limit	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"										mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, pfb_count, size)
//	@Produce		json
//	@Success		200	{array}		responses.Namespace
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/namespace [get]
func (handler *NamespaceHandler) List(c echo.Context) error {
	req, err := bindAndValidate[namespaceList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespace, err := handler.namespace.ListWithSort(c.Request().Context(), req.SortBy, pgSort(req.Sort), req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}
	response := make([]responses.Namespace, len(namespace))
	for i := range namespace {
		response[i] = responses.NewNamespace(namespace[i])
	}
	return returnArray(c, response)
}

type getBlobsRequest struct {
	Hash   string      `param:"hash"   validate:"required,base64"`
	Height types.Level `param:"height" validation:"required,min=1"`
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
//	@Router			/namespace_by_hash/{hash}/{height} [get]
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

type listByNamespace struct {
	Id      string `param:"id"      validate:"required,hexadecimal,len=56"`
	Version byte   `param:"version"`
	Limit   int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset"  validate:"omitempty,min=0"`
}

func (req *listByNamespace) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
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
//	@Router			/namespace/{id}/{version}/messages [get]
func (handler *NamespaceHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[listByNamespace](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	messages, err := handler.namespace.Messages(c.Request().Context(), ns.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	response := make([]responses.NamespaceMessage, len(messages))
	for i := range response {
		msg, err := responses.NewNamespaceMessage(messages[i])
		if err != nil {
			return handleError(c, err, handler.namespace)
		}
		response[i] = msg
	}

	return returnArray(c, response)
}

type listBlobsRequest struct {
	Limit      int         `query:"limit"      validate:"omitempty,min=1,max=100"`
	Offset     int         `query:"offset"     validate:"omitempty,min=0"`
	Sort       string      `query:"sort"       validate:"omitempty,oneof=asc desc"`
	SortBy     string      `query:"sort_by"    validate:"omitempty,oneof=time size"`
	Commitment string      `query:"commitment" validate:"omitempty,base64url"`
	Signers    StringArray `query:"signers"    validate:"omitempty,dive,address"`
	Namespaces StringArray `query:"namespaces" validate:"omitempty,dive,namespace"`
	Cursor     uint64      `query:"cursor"     validate:"omitempty,min=0"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (req *listBlobsRequest) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 0
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

func (req listBlobsRequest) toDbRequest(ctx context.Context, ns storage.INamespace, addrs storage.IAddress) (storage.ListBlobLogFilters, error) {
	fltrs := storage.ListBlobLogFilters{
		Limit:      req.Limit,
		Offset:     req.Offset,
		Sort:       pgSort(req.Sort),
		SortBy:     req.SortBy,
		Commitment: req.Commitment,
		Cursor:     req.Cursor,
		Namespaces: make([]uint64, len(req.Namespaces)),
	}
	if req.From > 0 {
		fltrs.From = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.To = time.Unix(req.To, 0).UTC()
	}

	var err error

	if len(req.Namespaces) > 0 {
		for i := range req.Namespaces {
			hash, err := base64.StdEncoding.DecodeString(req.Namespaces[i])
			if err != nil {
				return fltrs, err
			}

			n, err := ns.ByNamespaceIdAndVersion(ctx, hash[1:], hash[0])
			if err != nil {
				return fltrs, err
			}
			fltrs.Namespaces[i] = n.Id
		}
	}

	if len(req.Signers) > 0 {
		hash := make([][]byte, len(req.Signers))
		for i := range req.Signers {
			if _, hash[i], err = types.Address(req.Signers[i]).Decode(); err != nil {
				return fltrs, err
			}
		}

		if fltrs.Signers, err = addrs.IdByHash(ctx, hash...); err != nil {
			return fltrs, err
		}
	}

	return fltrs, nil
}

// Blobs godoc
//
//	@Summary		List all blobs with filters
//	@Description	Returns blobs
//	@Tags			namespace
//	@ID				get-blobs
//	@Param			limit		query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"										mininum(1)
//	@Param			sort		query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by		query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Param			commitment	query	string	false	"Commitment value in URLbase64 format"
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Param			signers		query	string	false	"Comma-separated celestia addresses"
//	@Param			namespaces	query	string	false	"Comma-separated celestia namespaces"
//	@Param			cursor		query	integer	false	"Last entity id which is used for cursor pagination"	mininum(1)
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		responses.LightBlobLog
//	@Failure		400	{object}	Error
//	@Router			/blob [get]
func (handler *NamespaceHandler) Blobs(c echo.Context) error {
	req, err := bindAndValidate[listBlobsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs, err := req.toDbRequest(c.Request().Context(), handler.namespace, handler.address)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}

	blob, err := handler.blobLogs.ListBlobs(c.Request().Context(), fltrs)
	if err != nil {
		return badRequestError(c, err)
	}

	response := make([]responses.LightBlobLog, len(blob))
	for i := range blob {
		response[i] = responses.NewLightBlobLog(blob[i])
	}

	return returnArray(c, response)
}

type postBlobRequest struct {
	Hash       string      `example:"AAAAAAAAAAAAAAAAAAAAAAAAAAAAs2bWWU6FOB0="     json:"hash"       validate:"required,namespace"`
	Height     types.Level `example:"123456"                                       json:"height"     validate:"required,min=1"`
	Commitment string      `example:"vbGakK59+Non81TE3ULg5Ve5ufT9SFm/bCyY+WLR3gg=" json:"commitment" validate:"required,base64"`
}

// Blob godoc
//
//	@Summary					Get namespace blob by commitment on height
//	@Description				Returns blob.
//	@Tags						namespace
//	@ID							get-blob
//	@Param						request	body postBlobRequest	true "Request body containing height, commitment and namespace hash"
//	@Accept						json
//	@Produce					json
//	@Success					200	{object}	responses.Blob
//	@Failure					400	{object}	Error
//	@Router						/blob [post]
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						apikey
//	@description				To authorize your requests you have to select the required tariff on our site. Then you receive api key to authorize. Api key should be passed via request header `apikey`.
func (handler *NamespaceHandler) Blob(c echo.Context) error {
	req, err := bindAndValidate[postBlobRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	blob, err := handler.blob.Blob(c.Request().Context(), req.Height, req.Hash, req.Commitment)
	if err != nil {
		return badRequestError(c, err)
	}

	response, err := responses.NewBlob(blob)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}

	return c.JSON(http.StatusOK, response)
}

// BlobMetadata godoc
//
//	@Summary		Get blob metadata by commitment on height
//	@Description	Returns blob metadata
//	@Tags			namespace
//	@ID				get-blob-metadata
//	@Param			request	body postBlobRequest	true "Request body containing height, commitment and namespace hash"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.BlobLog
//	@Failure		400	{object}	Error
//	@Router			/blob/metadata [post]
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						apikey
//	@description				To authorize your requests you have to select the required tariff on our site. Then you receive api key to authorize. Api key should be passed via request header `apikey`.
func (handler *NamespaceHandler) BlobMetadata(c echo.Context) error {
	req, err := bindAndValidate[postBlobRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	namespaceId, err := base64.StdEncoding.DecodeString(req.Hash)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId[1:], namespaceId[0])
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	blobMetadata, err := handler.blobLogs.Blob(c.Request().Context(), req.Height, ns.Id, req.Commitment)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	return c.JSON(http.StatusOK, responses.NewBlobLog(blobMetadata))
}

type getBlobLogsForNamespace struct {
	Id         string      `param:"id"         validate:"required,hexadecimal,len=56"`
	Version    byte        `param:"version"`
	Limit      int         `query:"limit"      validate:"omitempty,min=1,max=100"`
	Offset     int         `query:"offset"     validate:"omitempty,min=0"`
	Sort       string      `query:"sort"       validate:"omitempty,oneof=asc desc"`
	SortBy     string      `query:"sort_by"    validate:"omitempty,oneof=time size"`
	Commitment string      `query:"commitment" validate:"omitempty,base64url"`
	Joins      *bool       `query:"joins"      validate:"omitempty"`
	Signers    StringArray `query:"signers"    validate:"omitempty,dive,address"`
	Cursor     uint64      `query:"cursor"     validate:"omitempty,min=0"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (req getBlobLogsForNamespace) getCommitment() (string, error) {
	if req.Commitment == "" {
		return "", nil
	}
	data, err := base64.URLEncoding.DecodeString(req.Commitment)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (req *getBlobLogsForNamespace) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
	if req.Joins == nil {
		req.Joins = testsuite.Ptr(true)
	}
}

// GetBlobLogs godoc
//
//	@Summary		Get blob changes for namespace
//	@Description	Returns blob changes for namespace
//	@Tags			namespace
//	@ID				get-blob-logs
//	@Param			id			path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Param			version		path	integer	true	"Version of namespace"
//	@Param			limit		query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"										mininum(1)
//	@Param			sort		query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by		query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Param			commitment	query	string	false	"Commitment value in URLbase64 format"
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Param			joins		query	boolean	false	"Flag indicating whether entities of rollup, transaction and signer should be attached or not. Default: true"
//	@Param			signers		query	string	false	"Comma-separated celestia addresses"
//	@Param			cursor		query	integer	false	"Last entity id which is used for cursor pagination"	mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/namespace/{id}/{version}/blobs [get]
func (handler *NamespaceHandler) GetBlobLogs(c echo.Context) error {
	req, err := bindAndValidate[getBlobLogsForNamespace](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	cm, err := req.getCommitment()
	if err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	var ids []uint64
	if len(req.Signers) > 0 {
		hash := make([][]byte, len(req.Signers))
		for i := range req.Signers {
			_, h, err := types.Address(req.Signers[i]).Decode()
			if err != nil {
				return badRequestError(c, err)
			}
			hash[i] = h
		}
		ids, err = handler.address.IdByHash(c.Request().Context(), hash...)
		if err != nil {
			return handleError(c, err, handler.namespace)
		}
	}

	fltrs := storage.BlobLogFilters{
		Limit:      req.Limit,
		Offset:     req.Offset,
		Sort:       pgSort(req.Sort),
		SortBy:     req.SortBy,
		Commitment: cm,
		Joins:      *req.Joins,
		Signers:    ids,
		Cursor:     req.Cursor,
	}

	if req.From > 0 {
		fltrs.From = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.To = time.Unix(req.To, 0).UTC()
	}

	logs, err := handler.blobLogs.ByNamespace(
		c.Request().Context(),
		ns.Id,
		fltrs,
	)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	response := make([]responses.BlobLog, len(logs))
	for i := range response {
		response[i] = responses.NewBlobLog(logs[i])
	}

	return returnArray(c, response)
}

// Rollups godoc
//
//	@Summary		List rollups using the namespace
//	@Description	List rollups using the namespace
//	@Tags			namespace
//	@ID				get-namespace-rollups
//	@Param			id		path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Param			version	path	integer	true	"Version of namespace"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Rollup
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/namespace/{id}/{version}/rollups [get]
func (handler *NamespaceHandler) Rollups(c echo.Context) error {
	req, err := bindAndValidate[listByNamespace](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	ns, err := handler.namespace.ByNamespaceIdAndVersion(c.Request().Context(), namespaceId, req.Version)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	rollups, err := handler.rollups.RollupsByNamespace(
		c.Request().Context(),
		ns.Id,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	response := make([]responses.Rollup, len(rollups))
	for i := range response {
		response[i] = responses.NewRollup(&rollups[i])
	}

	return returnArray(c, response)
}

// BlobProofs godoc
//
//	@Summary		Get blob inclusion proofs
//	@Description	Returns blob inclusion proofs
//	@Tags			namespace
//	@ID				get-blob-proof
//	@Param			request	body postBlobRequest	true "Request body containing height, commitment and namespace hash"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.BlobLog
//	@Failure		400	{object}	Error
//	@Router			/blob/proofs [get]
func (handler *NamespaceHandler) BlobProofs(c echo.Context) error {
	req, err := bindAndValidate[postBlobRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	block, err := handler.node.Block(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	dataSquare, err := square.Construct(
		block.Block.Data.Txs.ToSliceOfBytes(),
		appconsts.SquareSizeUpperBound(0),
		appconsts.SubtreeRootThreshold(0),
	)

	if err != nil {
		return internalServerError(c, err)
	}

	startBlobIndex, endBlobIndex, err := responses.GetBlobShareIndexes(dataSquare, req.Hash, req.Commitment)
	if err != nil {
		return internalServerError(c, err)
	}
	blobSharesRange := share.Range{
		Start: startBlobIndex,
		End:   endBlobIndex,
	}

	eds, err := da.ExtendShares(share.ToBytes(dataSquare))
	if err != nil {
		return internalServerError(c, err)
	}

	namespaceBytes, err := base64.StdEncoding.DecodeString(req.Hash)
	if err != nil {
		return internalServerError(c, err)
	}

	namespace, err := share.NewNamespaceFromBytes(namespaceBytes)
	if err != nil {
		return internalServerError(c, err)
	}

	proofs, err := proof.NewShareInclusionProofFromEDS(eds, namespace, blobSharesRange)
	if err != nil {
		return handleError(c, err, handler.namespace)
	}

	return c.JSON(http.StatusOK, responses.NewProofs(proofs.ShareProofs))
}

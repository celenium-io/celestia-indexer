// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

type TxHandler struct {
	tx          storage.ITx
	events      storage.IEvent
	messages    storage.IMessage
	namespaces  storage.INamespace
	blobLogs    storage.IBlobLog
	state       storage.IState
	indexerName string
}

func NewTxHandler(
	tx storage.ITx,
	events storage.IEvent,
	messages storage.IMessage,
	namespaces storage.INamespace,
	blobLogs storage.IBlobLog,
	state storage.IState,
	indexerName string,
) *TxHandler {
	return &TxHandler{
		tx:          tx,
		events:      events,
		messages:    messages,
		namespaces:  namespaces,
		blobLogs:    blobLogs,
		state:       state,
		indexerName: indexerName,
	}
}

type getTxRequest struct {
	Hash string `param:"hash" validate:"required,hexadecimal,len=64"`
}

// Get godoc
//
//	@Summary		Get transaction by hash
//	@Description	Get transaction by hash
//	@Tags			transactions
//	@ID				get-transaction
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{object}	responses.Tx
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash} [get]
func (handler *TxHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	return c.JSON(http.StatusOK, responses.NewTx(tx))
}

// List godoc
//
//	@Summary		List transactions info
//	@Description	List transactions info
//	@Tags			transactions
//	@ID				list-transactions
//	@Param			limit				query	integer			false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset				query	integer			false	"Offset"						mininum(1)
//	@Param			sort				query	string			false	"Sort order"					Enums(asc, desc)
//	@Param			status				query	types.Status	false	"Comma-separated status list"
//	@Param			msg_type			query	types.MsgType	false	"Comma-separated message types list"
//	@Param			excluded_msg_type	query	types.MsgType	false	"Comma-separated message types list which should be excluded"
//	@Param			from				query	integer			false	"Time from in unix timestamp"	mininum(1)
//	@Param			to					query	integer			false	"Time to in unix timestamp"		mininum(1)
//	@Param			height				query	integer			false	"Block number"					mininum(1)
//	@Param			messages			query	boolean			false	"If true join messages"			mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx [get]
func (handler *TxHandler) List(c echo.Context) error {
	req, err := bindAndValidate[txListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.TxFilter{
		Limit:                int(req.Limit),
		Offset:               int(req.Offset),
		Sort:                 pgSort(req.Sort),
		Status:               req.Status,
		Height:               req.Height,
		MessageTypes:         types.NewMsgTypeBitMask(),
		ExcludedMessageTypes: types.NewMsgTypeBitMask(),
		WithMessages:         req.Messages,
	}
	if req.From > 0 {
		fltrs.TimeFrom = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.TimeTo = time.Unix(req.To, 0).UTC()
	}
	for i := range req.MsgType {
		fltrs.MessageTypes.SetByMsgType(storageTypes.MsgType(req.MsgType[i]))
	}
	for i := range req.ExcludedMsgType {
		fltrs.ExcludedMessageTypes.SetByMsgType(storageTypes.MsgType(req.ExcludedMsgType[i]))
	}

	txs, err := handler.tx.Filter(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

type getTxRequestWithPagination struct {
	Hash   string `param:"hash"   validate:"required,hexadecimal,len=64"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (p *getTxRequestWithPagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// GetEvents godoc
//
//	@Summary		Get transaction events
//	@Description	Get transaction events
//	@Tags			transactions
//	@ID				get-transaction-events
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Event
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/events [get]
func (handler *TxHandler) GetEvents(c echo.Context) error {
	req, err := bindAndValidate[getTxRequestWithPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	fltrs := storage.EventFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		Time:   tx.Time.UTC(),
	}

	events, err := handler.events.ByTxId(c.Request().Context(), tx.Id, fltrs)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Event, len(events))
	for i := range events {
		response[i] = responses.NewEvent(events[i])
	}
	return returnArray(c, response)
}

// GetMessages godoc
//
//	@Summary		Get transaction messages
//	@Description	Get transaction messages
//	@Tags			transactions
//	@ID				get-transaction-messages
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Message
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/messages [get]
func (handler *TxHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[getTxRequestWithPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	messages, err := handler.messages.ByTxId(c.Request().Context(), tx.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Message, len(messages))
	for i := range messages {
		response[i] = responses.NewMessage(messages[i])
	}
	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of transactions in network
//	@Description	Get count of transactions in network
//	@Tags			transactions
//	@ID				get-transactions-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/count [get]
func (handler *TxHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	return c.JSON(http.StatusOK, state.TotalTx)
}

// Genesis godoc
//
//	@Summary		List genesis transactions info
//	@Description	List genesis transactions info
//	@Tags			transactions
//	@ID				list-genesis-transactions
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/genesis [get]
func (handler *TxHandler) Genesis(c echo.Context) error {
	req, err := bindAndValidate[limitOffsetPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	txs, err := handler.tx.Genesis(c.Request().Context(), int(req.Limit), int(req.Offset), pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

// Namespaces godoc
//
//	@Summary		List namespaces affected by transaction
//	@Description	List namespaces affected by transaction
//	@Tags			transactions
//	@ID				list-namespaces-transactions
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"		mininum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"							mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.NamespaceMessage
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/namespace [get]
func (handler *TxHandler) Namespaces(c echo.Context) error {
	req, err := bindAndValidate[listForTx](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	messages, err := handler.namespaces.MessagesByTxId(c.Request().Context(), tx.Id, int(req.Limit), int(req.Offset))
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	response := make([]responses.NamespaceMessage, len(messages))
	for i := range messages {
		response[i], err = responses.NewNamespaceMessage(messages[i])
		if err != nil {
			return handleError(c, err, handler.tx)
		}
	}
	return returnArray(c, response)
}

// NamespacesCount godoc
//
//	@Summary		Count of namespaces affected by transaction
//	@Description	Count of namespaces affected by transaction
//	@Tags			transactions
//	@ID				list-namespaces-count-transactions
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/namespace/count [get]
func (handler *TxHandler) NamespacesCount(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	count, err := handler.namespaces.CountMessagesByTxId(c.Request().Context(), tx.Id)
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	return c.JSON(http.StatusOK, count)
}

type getBlobsForTx struct {
	Hash   string `param:"hash"    validate:"required,hexadecimal,len=64"`
	Limit  uint64 `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=time size"`
}

func (req *getBlobsForTx) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// Blobs godoc
//
//	@Summary		List blobs which was pushed by transaction
//	@Description	List blobs which was pushed by transaction
//	@Tags			transactions
//	@ID				list-transaction-blobs
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"				minlength(64)	maxlength(64)
//	@Param			limit	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"										mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/blobs [get]
func (handler *TxHandler) Blobs(c echo.Context) error {
	req, err := bindAndValidate[getBlobsForTx](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	blobs, err := handler.blobLogs.ByTxId(
		c.Request().Context(),
		tx.Id,
		storage.BlobLogFilters{
			Limit:  int(req.Limit),
			Offset: int(req.Offset),
			Sort:   pgSort(req.Sort),
			SortBy: req.SortBy,
		},
	)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}

	response := make([]responses.BlobLog, len(blobs))
	for i := range blobs {
		response[i] = responses.NewBlobLog(blobs[i])
	}
	return returnArray(c, response)
}

// BlobsCount godoc
//
//	@Summary		Count of blobs which was pushed by transaction
//	@Description	Count of blobs which was pushed by transaction
//	@Tags			transactions
//	@ID				transaction-blobs-count
//	@Param			hash	path	string	true	"Transaction hash in hexadecimal"				minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/tx/{hash}/blobs/count [get]
func (handler *TxHandler) BlobsCount(c echo.Context) error {
	req, err := bindAndValidate[getTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	tx, err := handler.tx.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	count, err := handler.blobLogs.CountByTxId(c.Request().Context(), tx.Id)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}
	return c.JSON(http.StatusOK, count)
}

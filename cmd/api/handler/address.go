// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	address     storage.IAddress
	txs         storage.ITx
	blobLogs    storage.IBlobLog
	messages    storage.IMessage
	state       storage.IState
	indexerName string
}

func NewAddressHandler(
	address storage.IAddress,
	txs storage.ITx,
	blobLogs storage.IBlobLog,
	messages storage.IMessage,
	state storage.IState,
	indexerName string,
) *AddressHandler {
	return &AddressHandler{
		address:     address,
		txs:         txs,
		blobLogs:    blobLogs,
		messages:    messages,
		state:       state,
		indexerName: indexerName,
	}
}

type getAddressRequest struct {
	Hash string `param:"hash" validate:"required,address"`
}

// Get godoc
//
//	@Summary		Get address info
//	@Description	Get address info
//	@Tags			address
//	@ID				get-address
//	@Param			hash	path	string	true	"Hash"	minlength(48)	maxlength(48)
//	@Produce		json
//	@Success		200	{object}	responses.Address
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash} [get]
func (handler *AddressHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getAddressRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewAddress(address))
}

// List godoc
//
//	@Summary		List address info
//	@Description	List address info
//	@Tags			address
//	@ID				list-address
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Address
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address [get]
func (handler *AddressHandler) List(c echo.Context) error {
	req, err := bindAndValidate[addressListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.AddressListFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
		Sort:   pgSort(req.Sort),
	}

	address, err := handler.address.ListWithBalance(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Address, len(address))
	for i := range address {
		response[i] = responses.NewAddress(address[i])
	}

	return returnArray(c, response)
}

// Transactions godoc
//
//	@Summary		Get address transactions
//	@Description	Get address transactions
//	@Tags			address
//	@ID				address-transactions
//	@Param			hash		path	string					true	"Hash"							minlength(48)	maxlength(48)
//	@Param			limit		query	integer					false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset		query	integer					false	"Offset"						minimum(1)
//	@Param			sort		query	string					false	"Sort order"					Enums(asc, desc)
//	@Param			status		query	storageTypes.Status		false	"Comma-separated status list"
//	@Param			msg_type	query	storageTypes.MsgType	false	"Comma-separated message types list"
//	@Param			from		query	integer					false	"Time from in unix timestamp"	minimum(1)
//	@Param			to			query	integer					false	"Time to in unix timestamp"		minimum(1)
//	@Param			height		query	integer					false	"Block number"					minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Tx
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/txs [get]
func (handler *AddressHandler) Transactions(c echo.Context) error {
	req, err := bindAndValidate[addressTxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	fltrs := storage.TxFilter{
		Limit:        int(req.Limit),
		Offset:       int(req.Offset),
		Sort:         pgSort(req.Sort),
		Status:       req.Status,
		Height:       req.Height,
		MessageTypes: storageTypes.NewMsgTypeBitMask(),
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

	txs, err := handler.txs.ByAddress(c.Request().Context(), address.Id, fltrs)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

type getAddressMessages struct {
	Hash    string      `param:"hash"     validate:"required,address"`
	Limit   uint64      `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  uint64      `query:"offset"   validate:"omitempty,min=0"`
	Sort    string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	MsgType StringArray `query:"msg_type" validate:"omitempty,dive,msg_type"`
}

func (p *getAddressMessages) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
	if p.MsgType == nil {
		p.MsgType = make(StringArray, 0)
	}
}

func (p *getAddressMessages) ToFilters() storage.AddressMsgsFilter {
	return storage.AddressMsgsFilter{
		Limit:        int(p.Limit),
		Offset:       int(p.Offset),
		Sort:         pgSort(p.Sort),
		MessageTypes: p.MsgType,
	}
}

// Messages godoc
//
//	@Summary		Get address messages
//	@Description	Get address messages
//	@Tags			address
//	@ID				address-messages
//	@Param			hash	path	string	true	"Hash"							minlength(48)	maxlength(48)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.MessageForAddress
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/messages [get]
func (handler *AddressHandler) Messages(c echo.Context) error {
	req, err := bindAndValidate[getAddressMessages](c)
	if err != nil {
		return badRequestError(c, err)
	}

	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	filters := req.ToFilters()
	msgs, err := handler.messages.ByAddress(c.Request().Context(), address.Id, filters)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.MessageForAddress, len(msgs))
	for i := range msgs {
		response[i] = responses.NewMessageForAddress(msgs[i])
	}

	return returnArray(c, response)
}

type getBlobLogsForAddress struct {
	Hash   string `param:"hash"    validate:"required,address"`
	Limit  uint64 `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=time size"`
}

func (req *getBlobLogsForAddress) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// Blobs godoc
//
//	@Summary		Get blobs pushed by address
//	@Description	Get blobs pushed by address
//	@Tags			address
//	@ID				address-blobs
//	@Param			hash	path	string	true	"Hash"											minlength(48)	maxlength(48)
//	@Param			limit	query	integer	false	"Count of requested entities"					minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"										minimum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/address/{hash}/blobs [get]
func (handler *AddressHandler) Blobs(c echo.Context) error {
	req, err := bindAndValidate[getBlobLogsForAddress](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	logs, err := handler.blobLogs.BySigner(
		c.Request().Context(),
		address.Id,
		storage.BlobLogFilters{
			Limit:  int(req.Limit),
			Offset: int(req.Offset),
			Sort:   pgSort(req.Sort),
			SortBy: req.SortBy,
		},
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.BlobLog, len(logs))
	for i := range response {
		response[i] = responses.NewBlobLog(logs[i])
	}

	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of addresses in network
//	@Description	Get count of addresses in network
//	@Tags			address
//	@ID				get-address-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/address/count [get]
func (handler *AddressHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	return c.JSON(http.StatusOK, state.TotalAccounts)
}

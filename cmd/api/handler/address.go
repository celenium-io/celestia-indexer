// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	celestials "github.com/celenium-io/celestial-module/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type AddressHandler struct {
	address       storage.IAddress
	txs           storage.ITx
	blobLogs      storage.IBlobLog
	messages      storage.IMessage
	delegations   storage.IDelegation
	undelegations storage.IUndelegation
	redelegations storage.IRedelegation
	vestings      storage.IVestingAccount
	grants        storage.IGrant
	celestial     celestials.ICelestial
	votes         storage.IVote
	state         storage.IState
	indexerName   string
}

func NewAddressHandler(
	address storage.IAddress,
	txs storage.ITx,
	blobLogs storage.IBlobLog,
	messages storage.IMessage,
	delegations storage.IDelegation,
	undelegations storage.IUndelegation,
	redelegations storage.IRedelegation,
	vestings storage.IVestingAccount,
	grants storage.IGrant,
	celestial celestials.ICelestial,
	votes storage.IVote,
	state storage.IState,
	indexerName string,
) *AddressHandler {
	return &AddressHandler{
		address:       address,
		txs:           txs,
		blobLogs:      blobLogs,
		messages:      messages,
		delegations:   delegations,
		undelegations: undelegations,
		redelegations: redelegations,
		vestings:      vestings,
		grants:        grants,
		celestial:     celestial,
		votes:         votes,
		state:         state,
		indexerName:   indexerName,
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
//	@Param			hash	path	string	true	"Hash"	minlength(47)	maxlength(128)
//	@Produce		json
//	@Success		200	{object}	responses.Address
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash} [get]
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

type addressListRequest struct {
	Limit  int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=id spendable delegated unbonding first_height last_height"`
}

func (p *addressListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
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
//	@Param			sort_by	query	string	false	"Sort field"					Enums(id, delegated, spendable, unbonding, first_height, last_height)
//	@Produce		json
//	@Success		200	{array}		responses.Address
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address [get]
func (handler *AddressHandler) List(c echo.Context) error {
	req, err := bindAndValidate[addressListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.AddressListFilter{
		Limit:     req.Limit,
		Offset:    req.Offset,
		Sort:      pgSort(req.Sort),
		SortField: req.SortBy,
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
//	@Param			hash		path	string					true	"Hash"							minlength(47)	maxlength(128)
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
//	@Router			/address/{hash}/txs [get]
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

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	fltrs := storage.TxFilter{
		Limit:        req.Limit,
		Offset:       req.Offset,
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

	txs, err := handler.txs.ByAddress(c.Request().Context(), addressId, fltrs)
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
	Limit   int         `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  int         `query:"offset"   validate:"omitempty,min=0"`
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
		Limit:        p.Limit,
		Offset:       p.Offset,
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
//	@Param			hash		path	string					true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit		query	integer					false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset		query	integer					false	"Offset"						minimum(1)
//	@Param			sort		query	string					false	"Sort order"					Enums(asc, desc)
//	@Param			msg_type	query	storageTypes.MsgType	false	"Comma-separated message types list"
//	@Produce		json
//	@Success		200	{array}		responses.MessageForAddress
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/messages [get]
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

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	filters := req.ToFilters()
	msgs, err := handler.messages.ByAddress(c.Request().Context(), addressId, filters)
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
	Limit  int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=time size"`
	Joins  *bool  `query:"joins"   validate:"omitempty"`
}

func (req *getBlobLogsForAddress) SetDefault() {
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

// Blobs godoc
//
//	@Summary		Get blobs pushed by address
//	@Description	Get blobs pushed by address
//	@Tags			address
//	@ID				address-blobs
//	@Param			hash	path	string	true	"Hash"											minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"					minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"										minimum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Param			joins	query	boolean	false	"Flag indicating whether entities of transaction and namespace should be attached or not. Default: true"
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/blobs [get]
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

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	logs, err := handler.blobLogs.BySigner(
		c.Request().Context(),
		addressId,
		storage.BlobLogFilters{
			Limit:  req.Limit,
			Offset: req.Offset,
			Sort:   pgSort(req.Sort),
			SortBy: req.SortBy,
			Joins:  *req.Joins,
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
//	@Router			/address/count [get]
func (handler *AddressHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	return c.JSON(http.StatusOK, state.TotalAccounts)
}

type getAddressDelegations struct {
	Hash     string `param:"hash"      validate:"required,address"`
	Limit    int    `query:"limit"     validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset"    validate:"omitempty,min=0"`
	ShowZero bool   `query:"show_zero" validate:"omitempty"`
}

func (req *getAddressDelegations) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
}

// Delegations godoc
//
//	@Summary		Get delegations made by address
//	@Description	Get delegations made by address
//	@Tags			address
//	@ID				address-delegations
//	@Param			hash		path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit		query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset		query	integer	false	"Offset"						minimum(1)
//	@Param			show_zero	query	boolean	false	"Show zero delegations"
//	@Produce		json
//	@Success		200	{array}		responses.Delegation
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/delegations [get]
func (handler *AddressHandler) Delegations(c echo.Context) error {
	req, err := bindAndValidate[getAddressDelegations](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	delegations, err := handler.delegations.ByAddress(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
		req.ShowZero,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Delegation, len(delegations))
	for i := range response {
		response[i] = responses.NewDelegation(delegations[i])
	}

	return returnArray(c, response)
}

type getAddressPageable struct {
	Hash   string `param:"hash"   validate:"required,address"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (req *getAddressPageable) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
}

// Undelegations godoc
//
//	@Summary		Get undelegations made by address
//	@Description	Get undelegations made by address
//	@Tags			address
//	@ID				address-undelegations
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Undelegation
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/undelegations [get]
func (handler *AddressHandler) Undelegations(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	undelegations, err := handler.undelegations.ByAddress(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Undelegation, len(undelegations))
	for i := range response {
		response[i] = responses.NewUndelegation(undelegations[i])
	}

	return returnArray(c, response)
}

// Redelegations godoc
//
//	@Summary		Get redelegations made by address
//	@Description	Get redelegations made by address
//	@Tags			address
//	@ID				address-redelegations
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Redelegation
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/redelegations [get]
func (handler *AddressHandler) Redelegations(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	redelegations, err := handler.redelegations.ByAddress(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Redelegation, len(redelegations))
	for i := range response {
		response[i] = responses.NewRedelegation(redelegations[i])
	}

	return returnArray(c, response)
}

type getAddressVestings struct {
	Hash      string `param:"hash"       validate:"required,address"`
	Limit     int    `query:"limit"      validate:"omitempty,min=1,max=100"`
	Offset    int    `query:"offset"     validate:"omitempty,min=0"`
	ShowEnded bool   `query:"show_ended" validate:"omitempty"`
}

func (req *getAddressVestings) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
}

// Vestings godoc
//
//	@Summary		Get vesting for address
//	@Description	Get vesting for address
//	@Tags			address
//	@ID				address-vesting
//	@Param			hash		path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit		query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset		query	integer	false	"Offset"						minimum(1)
//	@Param			show_ended	query	boolean	false	"Show finished vestings delegations"
//	@Produce		json
//	@Success		200	{array}		responses.Vesting
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/vestings [get]
func (handler *AddressHandler) Vestings(c echo.Context) error {
	req, err := bindAndValidate[getAddressVestings](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	vestings, err := handler.vestings.ByAddress(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
		req.ShowEnded,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Vesting, len(vestings))
	for i := range response {
		response[i] = responses.NewVesting(vestings[i])
	}

	return returnArray(c, response)
}

// Grants godoc
//
//	@Summary		Get grants made by address
//	@Description	Get grants made by address
//	@Tags			address
//	@ID				address-grants
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Grant
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/grants [get]
func (handler *AddressHandler) Grants(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	grants, err := handler.grants.ByGranter(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Grant, len(grants))
	for i := range response {
		response[i] = responses.NewGrant(grants[i])
	}
	return returnArray(c, response)
}

// Grantee godoc
//
//	@Summary		Get grants where address is grantee
//	@Description	Get grants where address is grantee
//	@Tags			address
//	@ID				address-grantee
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Grant
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/granters [get]
func (handler *AddressHandler) Grantee(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}
	grants, err := handler.grants.ByGrantee(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Grant, len(grants))
	for i := range response {
		response[i] = responses.NewGrant(grants[i])
	}
	return returnArray(c, response)
}

type addressStatsRequest struct {
	Hash       string `example:"celestia1glfkehhpvl55amdew2fnm6wxt7egy560mxdrj7" param:"hash"      swaggertype:"string"  validate:"required,address"`
	Timeframe  string `example:"hour"                                            param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"tps"                                             param:"name"      swaggertype:"string"  validate:"required,oneof=gas_used gas_wanted fee tx_count"`
	From       int64  `example:"1692892095"                                      query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         int64  `example:"1692892095"                                      query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Stats godoc
//
//	@Summary		Get address stats
//	@Description	Get address stats
//	@Tags			address
//	@ID				address-stats
//	@Param			hash		path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			name		path	string	true	"Series name"					Enums(gas_used, gas_wanted, fee, tx_count)
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.HistogramItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/stats/{name}/{timeframe} [get]
func (handler *AddressHandler) Stats(c echo.Context) error {
	req, err := bindAndValidate[addressStatsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	series, err := handler.address.Series(
		c.Request().Context(),
		addressId,
		storage.Timeframe(req.Timeframe),
		req.SeriesName,
		storage.NewSeriesRequest(req.From, req.To),
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.HistogramItem, len(series))
	for i := range series {
		response[i] = responses.NewHistogramItem(series[i])
	}
	return returnArray(c, response)
}

func (handler *AddressHandler) getIdByHash(ctx context.Context, hash []byte, address string) (uint64, error) {
	addressId, err := handler.address.IdByHash(ctx, hash)
	if err != nil {
		return 0, err
	}

	switch len(addressId) {
	case 0:
		return 0, errors.Wrap(errUnknownAddress, address)
	case 1:
		return addressId[0], nil
	default:
		return handler.address.IdByAddress(ctx, address, addressId...)
	}
}

// Celestials godoc
//
//	@Summary		Get list of celestial id for address
//	@Description	Get list of celestial id for address
//	@Tags			address
//	@ID				address-celestials
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Celestial
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/celestials [get]
func (handler *AddressHandler) Celestials(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}

	_, hash, err := types.Address(req.Hash).Decode()
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.getIdByHash(c.Request().Context(), hash, req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	celestials, err := handler.celestial.ByAddressId(
		c.Request().Context(),
		addressId,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]*responses.Celestial, len(celestials))
	for i := range celestials {
		response[i] = responses.NewCelestial(&celestials[i])
	}
	return returnArray(c, response)
}

// Votes godoc
//
//	@Summary		Get list of votes for address
//	@Description	Get list of votes for address
//	@Tags			address
//	@ID				address-votes
//	@Param			hash	path	string	true	"Hash"							minlength(47)	maxlength(128)
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)		maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Vote
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/address/{hash}/votes [get]
func (handler *AddressHandler) Votes(c echo.Context) error {
	req, err := bindAndValidate[getAddressPageable](c)
	if err != nil {
		return badRequestError(c, err)
	}

	addressId, err := handler.address.IdByAddress(c.Request().Context(), req.Hash)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	votes, err := handler.votes.ByVoterId(
		c.Request().Context(),
		addressId,
		storage.VoteFilters{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.Vote, len(votes))
	for i := range votes {
		response[i] = responses.NewVote(votes[i])
	}
	return returnArray(c, response)
}

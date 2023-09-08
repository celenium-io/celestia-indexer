package handler

import (
	"net/http"
	"time"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	_ "github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	address storage.IAddress
	txs     storage.ITx
}

func NewAddressHandler(address storage.IAddress, txs storage.ITx) *AddressHandler {
	return &AddressHandler{
		address: address,
		txs:     txs,
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

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err := handleError(c, err, handler.address); err != nil {
		return err
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
	req, err := bindAndValidate[limitOffsetPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	address, err := handler.address.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}

	response := make([]responses.Address, len(address))
	for i := range address {
		response[i] = responses.NewAddress(*address[i])
	}

	return returnArray(c, response)
}

// Transactions godoc
//
//	@Summary		Get address transactions
//	@Description	Get address transactions
//	@Tags			address
//	@ID				address-transactions
//	@Param			limit		query	integer	false	"Count of requested entities"			mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"								mininum(1)
//	@Param			sort		query	string	false	"Sort order"							Enums(asc, desc)
//	@Param			status		query	types.Status	false	"Comma-separated status list"
//	@Param			msg_type	query	types.MsgType	false	"Comma-separated message types list"
//	@Param			from		query	integer	false	"Time from in unix timestamp"			mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"				mininum(1)
//	@Param			height		query	integer	false	"Block number"							mininum(1)
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

	address, err := handler.address.ByHash(c.Request().Context(), req.Hash)
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}

	fltrs := storage.TxFilter{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
		Sort:   pgSort(req.Sort),
		Status: req.Status,
		Height: req.Height,
	}
	if req.From > 0 {
		fltrs.TimeFrom = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		fltrs.TimeTo = time.Unix(req.To, 0).UTC()
	}

	txs, err := handler.txs.ByAddress(c.Request().Context(), address.Id, fltrs)
	if err := handleError(c, err, handler.txs); err != nil {
		return err
	}
	response := make([]responses.Tx, len(txs))
	for i := range txs {
		response[i] = responses.NewTx(txs[i])
	}
	return returnArray(c, response)
}

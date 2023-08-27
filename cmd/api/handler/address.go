package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	address storage.IAddress
}

func NewAddressHandler(address storage.IAddress) *AddressHandler {
	return &AddressHandler{
		address: address,
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

	hash, err := responses.DecodeAddress(req.Hash)
	if err != nil {
		return badRequestError(c, err)
	}

	address, err := handler.address.ByHash(c.Request().Context(), hash)
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}
	response, err := responses.NewAddress(address)
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
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
		response[i], err = responses.NewAddress(*address[i])
		if err := handleError(c, err, handler.address); err != nil {
			return err
		}
	}

	return returnArray(c, response)
}

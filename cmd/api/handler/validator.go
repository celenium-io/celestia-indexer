// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type ValidatorHandler struct {
	validators storage.IValidator
}

func NewValidatorHandler(
	validators storage.IValidator,
) *ValidatorHandler {
	return &ValidatorHandler{
		validators: validators,
	}
}

type validatorRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

// Get godoc
//
//	@Summary		Get validator info
//	@Description	Get validator info
//	@Tags			validator
//	@ID				get-validator
//	@Param			id	path	integer	true	"Internal validator id"
//	@Produce		json
//	@Success		200	{object}	responses.Validator
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators/{id} [get]
func (handler *ValidatorHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[validatorRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	validator, err := handler.validators.GetByID(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	return c.JSON(http.StatusOK, responses.NewValidator(*validator))
}

// List godoc
//
//	@Summary		List validators
//	@Description	List validators
//	@Tags			validator
//	@ID				list-validator
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Produce		json
//	@Produce		json
//	@Success		200	{array}	responses.Validator
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/validators [get]
func (handler *ValidatorHandler) List(c echo.Context) error {
	req, err := bindAndValidate[limitOffsetPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	validators, err := handler.validators.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	response := make([]responses.Validator, len(validators))
	for i := range validators {
		response[i] = *responses.NewValidator(*validators[i])
	}

	return returnArray(c, response)
}

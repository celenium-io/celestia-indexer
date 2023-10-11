// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type ConstantHandler struct {
	constants     storage.IConstant
	denomMetadata storage.IDenomMetadata
	address       storage.IAddress
}

func NewConstantHandler(constants storage.IConstant, denomMetadata storage.IDenomMetadata, address storage.IAddress) *ConstantHandler {
	return &ConstantHandler{
		constants:     constants,
		denomMetadata: denomMetadata,
	}
}

// Get godoc
//
//	@Summary		Get network constants
//	@Description	Get network constants
//	@Tags			general
//	@ID				get-constants
//	@Produce		json
//	@Success		200	{object}	responses.Constants
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/v1/constants [get]
func (handler *ConstantHandler) Get(c echo.Context) error {
	consts, err := handler.constants.All(c.Request().Context())
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}
	dm, err := handler.denomMetadata.All(c.Request().Context())
	if err := handleError(c, err, handler.address); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responses.NewConstants(consts, dm))
}

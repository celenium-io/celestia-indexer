// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type ConstantHandler struct {
	constants     storage.IConstant
	denomMetadata storage.IDenomMetadata
	address       storage.IAddress
}

func NewConstantHandler(
	constants storage.IConstant,
	denomMetadata storage.IDenomMetadata,
	address storage.IAddress,
) *ConstantHandler {
	return &ConstantHandler{
		constants:     constants,
		denomMetadata: denomMetadata,
		address:       address,
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
//	@Router			/constants [get]
func (handler *ConstantHandler) Get(c echo.Context) error {
	consts, err := handler.constants.All(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.address)
	}
	dm, err := handler.denomMetadata.All(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.address)
	}
	return c.JSON(http.StatusOK, responses.NewConstants(consts, dm))
}

// Enums godoc
//
//	@Summary		Get celenium enumerators
//	@Description	Get celenium enumerators
//	@Tags			general
//	@ID				get-enums
//	@Produce		json
//	@Success		200	{object}	responses.Enums
//	@Router			/enums [get]
func (handler *ConstantHandler) Enums(c echo.Context) error {
	return c.JSON(http.StatusOK, responses.NewEnums())
}

// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
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
	rollup        storage.IRollup
}

func NewConstantHandler(
	constants storage.IConstant,
	denomMetadata storage.IDenomMetadata,
	rollup storage.IRollup,
) *ConstantHandler {
	return &ConstantHandler{
		constants:     constants,
		denomMetadata: denomMetadata,
		rollup:        rollup,
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
		return handleError(c, err, handler.rollup)
	}
	dm, err := handler.denomMetadata.All(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.rollup)
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
	tags, err := handler.rollup.Tags(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.rollup)
	}
	return c.JSON(http.StatusOK, responses.NewEnums(tags))
}

// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"strconv"

	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/labstack/echo/v4"
)

// GasHandler -
type GasHandler struct{}

func NewGasHandler() GasHandler {
	return GasHandler{}
}

type estimatePfbGas struct {
	Sizes StringArray `query:"sizes" validate:"required"`
}

// EstimateForPfb godoc
//
//	@Summary		Get estimated gas for pay for blob
//	@Description	Get estimated gas for pay for blob message with certain values of blob sizes
//	@Tags			gas
//	@ID				gas-estimate-for-pfb
//	@Param			sizes	query	string	true 	"Comma-separated array of blob sizes"
//	@Produce		json
//	@Success		200	{object}	uint64
//	@Failure		400	{object}	Error
//	@Router			/v1/gas/estimate_for_pfb [get]
func (handler GasHandler) EstimateForPfb(c echo.Context) error {
	req, err := bindAndValidate[estimatePfbGas](c)
	if err != nil {
		return badRequestError(c, err)
	}
	sizes := make([]uint32, len(req.Sizes))
	for i := range req.Sizes {
		size, err := strconv.ParseUint(req.Sizes[i], 10, 32)
		if err != nil {
			return badRequestError(c, err)
		}
		sizes[i] = uint32(size)
	}

	return c.JSON(http.StatusOK, blobtypes.DefaultEstimateGas(sizes))
}

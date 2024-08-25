// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"strconv"

	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	"github.com/labstack/echo/v4"
)

// GasHandler -
type GasHandler struct {
	state      storage.IState
	tx         storage.ITx
	constant   storage.IConstant
	blockStats storage.IBlockStats
	tracker    *gas.Tracker
}

func NewGasHandler(
	state storage.IState,
	tx storage.ITx,
	constant storage.IConstant,
	blockStats storage.IBlockStats,
	tracker *gas.Tracker,
) GasHandler {
	return GasHandler{
		state:      state,
		tx:         tx,
		blockStats: blockStats,
		constant:   constant,
		tracker:    tracker,
	}
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
//	@Param			sizes	query	string	true	"Comma-separated array of blob sizes"
//	@Produce		json
//	@Success		200	{object}	uint64
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/gas/estimate_for_pfb [get]
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
		//nolint:gosec
		sizes[i] = uint32(size)
	}

	gasPerBlobByteConst, err := handler.constant.Get(c.Request().Context(), types.ModuleNameBlob, "gas_per_blob_byte")
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	gasPerBlobByte := gasPerBlobByteConst.MustUint32()

	txSizeCostConst, err := handler.constant.Get(c.Request().Context(), types.ModuleNameAuth, "tx_size_cost_per_byte")
	if err != nil {
		return handleError(c, err, handler.tx)
	}
	txSizeCost := txSizeCostConst.MustUint64()

	return c.JSON(http.StatusOK, blobtypes.EstimateGas(sizes, gasPerBlobByte, txSizeCost))
}

// EstimatePrice godoc
//
//	@Summary		Get estimated gas price
//	@Description	Get estimated gas price based on historical data
//	@Tags			gas
//	@ID				gas-price
//	@Produce		json
//	@Success		200	{object}	responses.GasPrice
//	@Router			/gas/price [get]
func (handler GasHandler) EstimatePrice(c echo.Context) error {
	data := handler.tracker.State()
	return c.JSON(http.StatusOK, responses.GasPrice{
		Slow:   data.Slow,
		Median: data.Median,
		Fast:   data.Fast,
	})
}

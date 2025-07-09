// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"strconv"

	"github.com/celenium-io/celestia-indexer/cmd/api/gas"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	blobtypes "github.com/celestiaorg/celestia-app/v4/x/blob/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// GasHandler -
type GasHandler struct {
	state      storage.IState
	tx         storage.ITx
	constant   storage.IConstant
	blockStats storage.IBlockStats
	tracker    gas.ITracker
}

func NewGasHandler(
	state storage.IState,
	tx storage.ITx,
	constant storage.IConstant,
	blockStats storage.IBlockStats,
	tracker gas.ITracker,
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

type estimatePricePriorityRequest struct {
	Priority string `param:"priority" validate:"required,oneof=slow median fast"`
}

// EstimatePricePriority godoc
//
//	@Summary		Get estimated gas price with priority filter
//	@Description	Get estimated gas price with priority filter based on historical data
//	@Tags			gas
//	@ID				gas-price-priority
//	@Param			priority	path	string	true	"Priority"	Enums(slow, median, fast)
//	@Produce		json
//	@Success		200	{string} string
//	@Router			/gas/price/{priority} [get]
func (handler GasHandler) EstimatePricePriority(c echo.Context) error {
	req, err := bindAndValidate[estimatePricePriorityRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	data := handler.tracker.State()
	switch req.Priority {
	case "slow":
		return c.JSON(http.StatusOK, data.Slow)
	case "median":
		return c.JSON(http.StatusOK, data.Median)
	case "fast":
		return c.JSON(http.StatusOK, data.Fast)
	default:
		return badRequestError(c, errors.Errorf("invalid priority: %s", req.Priority))
	}
}

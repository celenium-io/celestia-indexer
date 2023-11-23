// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

// GasHandler -
type GasHandler struct {
	state      storage.IState
	tx         storage.ITx
	blockStats storage.IBlockStats
}

func NewGasHandler(
	state storage.IState,
	tx storage.ITx,
	blockStats storage.IBlockStats,
) GasHandler {
	return GasHandler{
		state:      state,
		tx:         tx,
		blockStats: blockStats,
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

const (
	estimationGasPriceBlocksCount = 5
)

// EstimatePrice godoc
//
//	@Summary		Get estimated gas price
//	@Description	Get estimated gas price based on historical data
//	@Tags			gas
//	@ID				gas-price
//	@Produce		json
//	@Success		200	{object}	responses.GasPrice
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/v1/gas/price [get]
func (handler GasHandler) EstimatePrice(c echo.Context) error {
	ctx := c.Request().Context()
	states, err := handler.state.List(ctx, 1, 0, sdk.SortOrderAsc)
	if err != nil {
		return handleError(c, err, handler.state)
	}
	if len(states) == 0 {
		return c.JSON(http.StatusNoContent, []any{})
	}
	state := states[0]
	lastBlockHeight := state.LastHeight - estimationGasPriceBlocksCount

	var (
		wg     sync.WaitGroup
		result = make(chan gasPrice, estimationGasPriceBlocksCount)
		errs   = make(chan error)
	)
	for height := state.LastHeight; height > lastBlockHeight; height-- {
		wg.Add(1)
		go handler.computeGasPriceEstimationForBlock(ctx, height, result, errs, &wg)
	}
	wg.Wait()

	gas := newGasPrice()

	for {
		select {
		case <-ctx.Done():
			return c.JSON(http.StatusOK, Error{
				Message: ctx.Err().Error(),
			})
		case err := <-errs:
			return internalServerError(c, err)
		case gp := <-result:
			gas.percentiles[0] = gas.percentiles[0].Add(gp.percentiles[0])
			gas.percentiles[1] = gas.percentiles[1].Add(gp.percentiles[1])
			gas.percentiles[2] = gas.percentiles[2].Add(gp.percentiles[2])

			gas.blocks = append(gas.blocks, responses.GasBlock{
				Height:       uint64(gp.stats.Height),
				GasWanted:    uint64(gp.stats.GasLimit),
				GasUsed:      uint64(gp.stats.GasUsed),
				Fee:          gp.stats.Fee.String(),
				GasPrice:     gp.stats.Fee.Div(decimal.NewFromInt(gp.stats.GasLimit)).String(),
				GasUsedRatio: decimal.NewFromInt(gp.stats.GasUsed).Div(decimal.NewFromInt(gp.stats.GasLimit)).String(),
				TxCount:      uint64(gp.stats.TxCount),
				Percentiles: []string{
					gp.percentiles[0].String(),
					gp.percentiles[1].String(),
					gp.percentiles[2].String(),
				},
			})

			if len(result) == 0 {
				gas.percentiles[0] = gas.percentiles[0].Div(decimal.NewFromInt(estimationGasPriceBlocksCount))
				gas.percentiles[1] = gas.percentiles[1].Div(decimal.NewFromInt(estimationGasPriceBlocksCount))
				gas.percentiles[2] = gas.percentiles[2].Div(decimal.NewFromInt(estimationGasPriceBlocksCount))

				return c.JSON(http.StatusOK, gas.toResponse())
			}
		}
	}
}

type gasPrice struct {
	percentiles []decimal.Decimal
	stats       storage.BlockStats

	blocks []responses.GasBlock
}

func newGasPrice() gasPrice {
	return gasPrice{
		percentiles: []decimal.Decimal{
			decimal.New(0, 1),
			decimal.New(0, 1),
			decimal.New(0, 1),
		},
		blocks: make([]responses.GasBlock, 0),
	}
}

func (gp gasPrice) toResponse() responses.GasPrice {
	return responses.GasPrice{
		Slow:           gp.percentiles[0].String(),
		Median:         gp.percentiles[1].String(),
		Fast:           gp.percentiles[2].String(),
		ComputedBlocks: gp.blocks,
	}
}

var (
	minGasPrice = decimal.NewFromFloat(appconsts.DefaultMinGasPrice)
)

func (handler GasHandler) computeGasPriceEstimationForBlock(ctx context.Context, height types.Level, result chan<- gasPrice, errs chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	block, err := handler.blockStats.ByHeight(ctx, height)
	if err != nil {
		errs <- err
		return
	}
	if block.TxCount == 0 {
		result <- newGasPrice()
		return
	}

	txs, err := handler.tx.Gas(ctx, height)
	if err != nil {
		errs <- err
		return
	}
	sort.Sort(storage.ByGasPrice(txs))

	var (
		gp      = newGasPrice()
		sumGas  = txs[0].GasWanted
		txIndex = 0
	)

	gp.stats = block

	for i, p := range []float64{.10, .50, .99} {
		threshold := uint64(float64(block.GasLimit) * p)
		for sumGas < int64(threshold) && txIndex < len(txs) {
			txIndex++
			sumGas += txs[txIndex].GasWanted
		}
		if txs[txIndex].GasPrice.LessThan(minGasPrice) {
			gp.percentiles[txIndex] = minGasPrice.Copy()
		} else {
			gp.percentiles[i] = txs[txIndex].GasPrice
		}
	}

	result <- gp
}

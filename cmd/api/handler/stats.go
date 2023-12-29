// SPDX-FileCopyrightText: 2023 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	repo   storage.IStats
	nsRepo storage.INamespace
	price  storage.IPrice
	state  storage.IState
}

func NewStatsHandler(repo storage.IStats, nsRepo storage.INamespace, price storage.IPrice, state storage.IState) StatsHandler {
	return StatsHandler{
		repo:   repo,
		nsRepo: nsRepo,
		price:  price,
		state:  state,
	}
}

type summaryRequest struct {
	Table    string `example:"block"      param:"table"    swaggertype:"string"  validate:"required,oneof=block block_stats tx event message validator"`
	Function string `example:"count"      param:"function" swaggertype:"string"  validate:"required,oneof=avg sum min max count"`
	Column   string `example:"fee"        query:"column"   swaggertype:"string"  validate:"omitempty"`
	From     uint64 `example:"1692892095" query:"from"     swaggertype:"integer" validate:"omitempty,min=1"`
	To       uint64 `example:"1692892095" query:"to"       swaggertype:"integer" validate:"omitempty,min=1"`
}

// Summary godoc
//
//	@Summary				Get value by table and function
//	@Description.markdown	summary
//	@Tags					stats
//	@ID						stats-summary
//	@Param					table		path	string	true	"Table name"	Enums(block, block_stats, tx, event, message, validator)
//	@Param					function	path	string	true	"Function name"	Enums(min, max, avg, sum, count)
//	@Param					column		query	string	false	"Column name which will be used for computation. Optional for count."
//	@Param					from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param					to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce				json
//	@Success				200	{object}	string
//	@Failure				400	{object}	Error
//	@Failure				500	{object}	Error
//	@Router					/v1/stats/summary/{table}/{function} [get]
func (sh StatsHandler) Summary(c echo.Context) error {
	req, err := bindAndValidate[summaryRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	var (
		summary      string
		countRequest = storage.CountRequest{
			Table: req.Table,
			From:  req.From,
			To:    req.To,
		}
	)

	if req.Function == "count" {
		summary, err = sh.repo.Count(c.Request().Context(), countRequest)
	} else {
		summary, err = sh.repo.Summary(c.Request().Context(), storage.SummaryRequest{
			CountRequest: countRequest,
			Function:     req.Function,
			Column:       req.Column,
		})
	}
	if err != nil {
		if errors.Is(err, storage.ErrValidation) {
			return badRequestError(c, err)
		}
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, summary)
}

type histogramRequest struct {
	Table     string `example:"block"      param:"table"     swaggertype:"string"  validate:"required,oneof=block block_stats tx event message"`
	Function  string `example:"count"      param:"function"  swaggertype:"string"  validate:"required,oneof=avg sum min max count"`
	Timeframe string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day week month year"`
	Column    string `example:"fee"        query:"column"    swaggertype:"string"  validate:"omitempty"`
	From      uint64 `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To        uint64 `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Histogram godoc
//
//	@Summary				Get histogram
//	@Description.markdown	histogram
//	@Tags					stats
//	@ID						stats-histogram
//	@Param					table		path	string	true	"Table name"	Enums(block, block_stats, tx, event, message)
//	@Param					function	path	string	true	"Function name"	Enums(min, max, avg, sum, count)
//	@Param					timeframe	path	string	true	"Timeframe"		Enums(hour, day, week, month, year)
//	@Param					column		query	string	false	"Column name which will be used for computation. Optional for count"
//	@Param					from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param					to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce				json
//	@Success				200	{array}		responses.HistogramItem
//	@Failure				400	{object}	Error
//	@Failure				500	{object}	Error
//	@Router					/v1/stats/histogram/{table}/{function}/{timeframe} [get]
func (sh StatsHandler) Histogram(c echo.Context) error {
	req, err := bindAndValidate[histogramRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	var (
		histogram    []storage.HistogramItem
		countRequest = storage.CountRequest{
			Table: req.Table,
			From:  req.From,
			To:    req.To,
		}
	)

	if req.Function == "count" {
		histogram, err = sh.repo.HistogramCount(c.Request().Context(), storage.HistogramCountRequest{
			CountRequest: countRequest,
			Timeframe:    storage.Timeframe(req.Timeframe),
		})
	} else {
		histogram, err = sh.repo.Histogram(c.Request().Context(), storage.HistogramRequest{
			SummaryRequest: storage.SummaryRequest{
				CountRequest: countRequest,
				Function:     req.Function,
				Column:       req.Column,
			},
			Timeframe: storage.Timeframe(req.Timeframe),
		})
	}
	if err != nil {
		return internalServerError(c, err)
	}

	response := make([]responses.HistogramItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewHistogramItem(histogram[i])
	}

	return returnArray(c, response)
}

// TPS godoc
//
//	@Summary		Get tps
//	@Description	Returns transaction per seconds statistics
//	@Tags			stats
//	@ID				stats-tps
//	@Produce		json
//	@Success		200	{object}	responses.TPS
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/tps [get]
func (sh StatsHandler) TPS(c echo.Context) error {
	tps, err := sh.repo.TPS(c.Request().Context())
	if err != nil {
		return internalServerError(c, err)
	}
	return c.JSON(http.StatusOK, responses.NewTPS(tps))
}

// TxCountHourly24h godoc
//
//	@Summary		Get tx count histogram for last 24 hours by hour
//	@Description	Get tx count histogram for last 24 hours by hour
//	@Tags			stats
//	@ID				stats-tx-count-24h
//	@Produce		json
//	@Success		200	{array}		responses.TxCountHistogramItem
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/tx_count_24h [get]
func (sh StatsHandler) TxCountHourly24h(c echo.Context) error {
	histogram, err := sh.repo.TxCountForLast24h(c.Request().Context())
	if err != nil {
		return internalServerError(c, err)
	}
	response := make([]responses.TxCountHistogramItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewTxCountHistogramItem(histogram[i])
	}
	return returnArray(c, response)
}

type namespaceUsageRequest struct {
	Top *int `example:"100" query:"top" validate:"omitempty,min=1,max=100"`
}

// NamespaceUsage godoc
//
//	@Summary		Get namespaces with sorting by size.
//	@Description	Get namespaces with sorting by size. Returns top 100 namespaces. Namespaces which is not included to top 100 grouped into 'others' item
//	@Tags			stats
//	@ID				stats-namespace-usage
//	@Param			top	query	integer	false	"Count of entities"	mininum(1)	maximum(100)
//	@Produce		json
//	@Success		200	{array}		responses.NamespaceUsage
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/namespace/usage [get]
func (sh StatsHandler) NamespaceUsage(c echo.Context) error {
	req, err := bindAndValidate[namespaceUsageRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	if req.Top == nil {
		top := 10
		req.Top = &top
	}

	namespaces, err := sh.nsRepo.ListWithSort(c.Request().Context(), "size", sdk.SortOrderDesc, *req.Top, 0)
	if err != nil {
		return internalServerError(c, err)
	}

	var top100Size int64
	response := make([]responses.NamespaceUsage, len(namespaces))
	for i := range namespaces {
		response[i] = responses.NewNamespaceUsage(namespaces[i])
		top100Size += response[i].Size
	}

	state, err := sh.state.List(c.Request().Context(), 1, 0, sdk.SortOrderAsc)
	if err != nil {
		return internalServerError(c, err)
	}
	if len(state) == 0 {
		return returnArray(c, response)
	}

	response = append(response, responses.NamespaceUsage{
		Name:    "others",
		Size:    state[0].TotalBlobsSize - top100Size,
		Version: nil,
	})

	return returnArray(c, response)
}

type seriesRequest struct {
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day week month year"`
	SeriesName string `example:"tps"        param:"name"      swaggertype:"string"  validate:"required,oneof=blobs_size tps bps fee supply_change block_time tx_count events_count gas_price gas_efficiency gas_used gas_limit"`
	From       uint64 `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         uint64 `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Series godoc
//
//	@Summary		Get histogram with precomputed stats
//	@Description	Get histogram with precomputed stats by series name and timeframe
//	@Tags			stats
//	@ID				stats-series
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, week, month, year)
//	@Param			name		path	string	true	"Series name"					Enums(blobs_size, tps, bps, fee, supply_change, block_time, tx_count, events_count, gas_price, gas_efficiency, gas_used, gas_limit)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.SeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/series/{name}/{timeframe} [get]
func (sh StatsHandler) Series(c echo.Context) error {
	req, err := bindAndValidate[seriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	histogram, err := sh.repo.Series(c.Request().Context(), storage.Timeframe(req.Timeframe), req.SeriesName, storage.SeriesRequest{
		From: req.From,
		To:   req.To,
	})
	if err != nil {
		return internalServerError(c, err)
	}

	response := make([]responses.SeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

type namespaceSeriesRequest struct {
	Id         string `example:"0011223344" param:"id"        swaggertype:"string"  validate:"required,hexadecimal,len=56"`
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day week month year"`
	SeriesName string `example:"size"       param:"name"      swaggertype:"string"  validate:"required,oneof=pfb_count size"`
	From       uint64 `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         uint64 `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// NamespaceSeries godoc
//
//	@Summary		Get histogram for namespace with precomputed stats
//	@Description	Get histogram for namespace with precomputed stats by series name and timeframe
//	@Tags			stats
//	@ID				stats-ns-series
//	@Param			id			path	string	true	"Namespace id in hexadecimal"	minlength(56)	maxlength(56)
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, week, month, year)
//	@Param			name		path	string	true	"Series name"					Enums(pfb_count, size)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.SeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/namespace/series/{id}/{name}/{timeframe} [get]
func (sh StatsHandler) NamespaceSeries(c echo.Context) error {
	req, err := bindAndValidate[namespaceSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	namespaceId, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	namespace, err := sh.nsRepo.ByNamespaceId(c.Request().Context(), namespaceId)
	if err != nil {
		return handleError(c, err, sh.nsRepo)
	}
	if len(namespace) == 0 {
		return c.JSON(http.StatusOK, []any{})
	}

	histogram, err := sh.repo.NamespaceSeries(
		c.Request().Context(),
		storage.Timeframe(req.Timeframe),
		req.SeriesName,
		namespace[0].Id,
		storage.SeriesRequest{
			From: req.From,
			To:   req.To,
		})
	if err != nil {
		return internalServerError(c, err)
	}

	response := make([]responses.SeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

type priceSeriesRequest struct {
	Timeframe string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=1m 1h 1d"`
	From      uint64 `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To        uint64 `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// PriceSeries godoc
//
//	@Summary		Get histogram with TIA price
//	@Description	Get histogram with TIA price
//	@Tags			stats
//	@ID				stats-price-series
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(1m, 1h, 1d)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Price
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/price/series/{timeframe} [get]
func (sh StatsHandler) PriceSeries(c echo.Context) error {
	req, err := bindAndValidate[priceSeriesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	histogram, err := sh.price.Get(c.Request().Context(), req.Timeframe, int64(req.From), int64(req.To), 100)
	if err != nil {
		return handleError(c, err, sh.nsRepo)
	}

	response := make([]responses.Price, len(histogram))
	for i := range histogram {
		response[i] = responses.NewPrice(histogram[i])
	}
	return returnArray(c, response)
}

// PriceCurrent godoc
//
//	@Summary		Get current TIA price
//	@Description	Get current TIA price
//	@Tags			stats
//	@ID				stats-price-current
//	@Produce		json
//	@Success		200	{object}	responses.Price
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/stats/price/current [get]
func (sh StatsHandler) PriceCurrent(c echo.Context) error {
	price, err := sh.price.Last(c.Request().Context())
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, responses.NewPrice(price))
}

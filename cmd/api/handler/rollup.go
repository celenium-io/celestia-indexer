// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	testsuite "github.com/celenium-io/celestia-indexer/internal/test_suite"
	"github.com/labstack/echo/v4"
)

type RollupHandler struct {
	rollups   storage.IRollup
	namespace storage.INamespace
	blobs     storage.IBlobLog
}

func NewRollupHandler(
	rollups storage.IRollup,
	namespace storage.INamespace,
	blobs storage.IBlobLog,
) RollupHandler {
	return RollupHandler{
		rollups:   rollups,
		namespace: namespace,
		blobs:     blobs,
	}
}

// Leaderboard godoc
//
//	@Summary		List rollups info
//	@Description	List rollups info
//	@Tags			rollup
//	@ID				list-rollup
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"		Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. Default: size"		Enums(time, blobs_count, size, fee)
//	@Produce		json
//	@Success		200	{array}		responses.RollupWithStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup [get]
func (handler RollupHandler) Leaderboard(c echo.Context) error {
	req, err := bindAndValidate[rollupList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	rollups, err := handler.rollups.Leaderboard(c.Request().Context(), req.SortBy, pgSort(req.Sort), req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	response := make([]responses.RollupWithStats, len(rollups))
	for i := range rollups {
		response[i] = responses.NewRollupWithStats(rollups[i])
	}
	return returnArray(c, response)
}

// LeaderboardDay godoc
//
//	@Summary		List rollups info with stats by previous 24 hours
//	@Description	List rollups info with stats by previous 24 hours
//	@Tags			rollup
//	@ID				list-rollup-24h
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"		Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. Default: mb_price"	Enums(avg_size, blobs_count, total_size, total_fee, throughput, namespace_count, pfb_count, mb_price)
//	@Produce		json
//	@Success		200	{array}		responses.RollupWithDayStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/day [get]
func (handler RollupHandler) LeaderboardDay(c echo.Context) error {
	req, err := bindAndValidate[rollupDayList](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	rollups, err := handler.rollups.LeaderboardDay(c.Request().Context(), req.SortBy, pgSort(req.Sort), req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	response := make([]responses.RollupWithDayStats, len(rollups))
	for i := range rollups {
		response[i] = responses.NewRollupWithDayStats(rollups[i])
	}
	return returnArray(c, response)
}

// Get godoc
//
//	@Summary		Get rollup info
//	@Description	Get rollup info
//	@Tags			rollup
//	@ID				get-rollup
//	@Param			id	path	integer	true	"Internal identity"	mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.Rollup
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id} [get]
func (handler RollupHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getById](c)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	return c.JSON(http.StatusOK, responses.NewRollupWithStats(rollup))
}

type getRollupPages struct {
	Id     uint64 `param:"id"     validate:"required,min=1"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (req *getRollupPages) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
}

// GetNamespaces godoc
//
//	@Summary		Get rollup namespaces info
//	@Description	Get rollup namespaces info
//	@Tags			rollup
//	@ID				get-rollup-namespaces
//	@Param			id		path	integer	true	"Internal identity"				mininum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Namespace
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id}/namespaces [get]
func (handler RollupHandler) GetNamespaces(c echo.Context) error {
	req, err := bindAndValidate[getRollupPages](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	namespaceIds, err := handler.rollups.Namespaces(c.Request().Context(), req.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	if len(namespaceIds) == 0 {
		return c.JSON(http.StatusOK, []any{})
	}

	namespaces, err := handler.namespace.GetByIds(c.Request().Context(), namespaceIds...)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	response := make([]responses.Namespace, len(namespaces))
	for i := range namespaces {
		response[i] = responses.NewNamespace(namespaces[i])
	}

	return returnArray(c, response)
}

type getRollupPagesWithSort struct {
	Id     uint64 `param:"id"      validate:"required,min=1"`
	Limit  int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset"  validate:"omitempty,min=0"`
	Sort   string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=time size"`
	Joins  *bool  `query:"joins"   validate:"omitempty"`
}

func (p *getRollupPagesWithSort) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
	}
	if p.Joins == nil {
		p.Joins = testsuite.Ptr(true)
	}
}

// GetBlobs godoc
//
//	@Summary		Get rollup blobs
//	@Description	Get rollup blobs
//	@Tags			rollup
//	@ID				get-rollup-blobs
//	@Param			id		path	integer	true	"Internal identity"								mininum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"										mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Param			joins	query	boolean	false	"Flag indicating whether entities of transaction and signer should be attached or not. Default: true"
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id}/blobs [get]
func (handler RollupHandler) GetBlobs(c echo.Context) error {
	req, err := bindAndValidate[getRollupPagesWithSort](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	providers, err := handler.rollups.Providers(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	if len(providers) == 0 {
		return c.JSON(http.StatusOK, []any{})
	}

	blobs, err := handler.blobs.ByProviders(c.Request().Context(), providers, storage.BlobLogFilters{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   pgSort(req.Sort),
		SortBy: req.SortBy,
		Joins:  *req.Joins,
	})
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	response := make([]responses.BlobLog, len(blobs))
	for i := range blobs {
		response[i] = responses.NewBlobLog(blobs[i])
	}
	return returnArray(c, response)
}

type rollupStatsRequest struct {
	Id         uint64 `example:"1"          param:"id"        swaggertype:"integer" validate:"required,min=1"`
	Timeframe  string `example:"hour"       param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day month"`
	SeriesName string `example:"tps"        param:"name"      swaggertype:"string"  validate:"required,oneof=blobs_count size size_per_blob fee"`
	From       int64  `example:"1692892095" query:"from"      swaggertype:"integer" validate:"omitempty,min=1"`
	To         int64  `example:"1692892095" query:"to"        swaggertype:"integer" validate:"omitempty,min=1"`
}

// Stats godoc
//
//	@Summary		Get rollup stats
//	@Description	Get rollup stats
//	@Tags			rollup
//	@ID				get-rollup-stats
//	@Param			id			path	integer	true	"Internal identity"				mininum(1)
//	@Param			name		path	string	true	"Series name"					Enums(blobs_count, size, size_per_blob, fee)
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.HistogramItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id}/stats/{name}/{timeframe} [get]
func (handler RollupHandler) Stats(c echo.Context) error {
	req, err := bindAndValidate[rollupStatsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	histogram, err := handler.rollups.Series(
		c.Request().Context(),
		req.Id,
		req.Timeframe,
		req.SeriesName,
		storage.NewSeriesRequest(req.From, req.To),
	)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.HistogramItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewHistogramItem(histogram[i])
	}
	return returnArray(c, response)
}

// AllSeries godoc
//
//	@Summary		Get series for all rollups
//	@Description	Get series for all rollups
//	@Tags			rollup
//	@ID				get-rollup-all-series
//	@Produce		json
//	@Success		200	{array}		responses.RollupAllSeriesItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/stats/series [get]
func (handler RollupHandler) AllSeries(c echo.Context) error {
	histogram, err := handler.rollups.AllSeries(
		c.Request().Context(),
	)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.RollupAllSeriesItem, len(histogram))
	for i := range histogram {
		response[i] = responses.NewRollupAllSeriesItem(histogram[i])
	}
	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get count of rollups in network
//	@Description	Get count of rollups in network
//	@Tags			rollup
//	@ID				get-rollups-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/rollup/count [get]
func (handler RollupHandler) Count(c echo.Context) error {
	count, err := handler.rollups.Count(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	return c.JSON(http.StatusOK, count)
}

type rollupBySlugRequest struct {
	Slug string `param:"slug" validate:"required"`
}

// BySlug godoc
//
//	@Summary		Get rollup by slug
//	@Description	Get rollup by slug
//	@Tags			rollup
//	@ID				get-rollup-by-slug
//	@Param			slug	path	string	true	"Slug"
//	@Produce		json
//	@Success		200	{object}	responses.Rollup
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/slug/{slug} [get]
func (handler RollupHandler) BySlug(c echo.Context) error {
	req, err := bindAndValidate[rollupBySlugRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.BySlug(c.Request().Context(), req.Slug)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	return c.JSON(http.StatusOK, responses.NewRollupWithStats(rollup))
}

type rollupDistributionRequest struct {
	Id         uint64 `example:"1"    param:"id"        swaggertype:"integer" validate:"required,min=1"`
	Timeframe  string `example:"hour" param:"timeframe" swaggertype:"string"  validate:"required,oneof=hour day"`
	SeriesName string `example:"tps"  param:"name"      swaggertype:"string"  validate:"required,oneof=blobs_count size size_per_blob fee_per_blob"`
}

// Distribution godoc
//
//	@Summary		Get rollup distribution
//	@Description	Get rollup distribution
//	@Tags			rollup
//	@ID				get-rollup-distribution
//	@Param			id			path	integer	true	"Internal identity"	mininum(1)
//	@Param			name		path	string	true	"Series name"		Enums(blobs_count, size, size_per_blob, fee_per_blob)
//	@Param			timeframe	path	string	true	"Timeframe"			Enums(hour, day)
//	@Produce		json
//	@Success		200	{array}		responses.DistributionItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id}/distribution/{name}/{timeframe} [get]
func (handler RollupHandler) Distribution(c echo.Context) error {
	req, err := bindAndValidate[rollupDistributionRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	distr, err := handler.rollups.Distribution(
		c.Request().Context(),
		req.Id,
		req.SeriesName,
		req.Timeframe,
	)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.DistributionItem, len(distr))
	for i := range distr {
		response[i] = responses.NewDistributionItem(distr[i], req.Timeframe)
	}
	return returnArray(c, response)
}

type exportBlobsRequest struct {
	Id   uint64 `example:"10"         param:"id"   swaggertype:"integer" validate:"required,min=1"`
	From int64  `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64  `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

// ExportBlobs godoc
//
//	@Summary		Export rollup blobs
//	@Description	Export rollup blobs
//	@Tags			rollup
//	@ID				rollup-export
//	@Param			id		path	integer	true	"Internal identity"				mininum(1)
//	@Param			from	query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to		query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Success		200
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/rollup/{id}/export [get]
func (handler RollupHandler) ExportBlobs(c echo.Context) error {
	req, err := bindAndValidate[exportBlobsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	providers, err := handler.rollups.Providers(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	if len(providers) == 0 {
		return c.JSON(http.StatusOK, []any{})
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
	c.Response().WriteHeader(http.StatusOK)

	var (
		from time.Time
		to   time.Time
	)
	if req.From > 0 {
		from = time.Unix(req.From, 0).UTC()
	}
	if req.To > 0 {
		to = time.Unix(req.To, 0).UTC()
	}

	err = handler.blobs.ExportByProviders(
		c.Request().Context(),
		providers,
		from,
		to,
		c.Response(),
	)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	return nil
}

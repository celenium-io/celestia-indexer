// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
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
//	@Param			sort_by	query	string	false	"Sort field. Default: size"		Enums(time, blobs_count, size)
//	@Produce		json
//	@Success		200	{array}		responses.RollupWithStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup [get]
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
//	@Router			/v1/rollup/{id} [get]
func (handler RollupHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getById](c)
	if err != nil {
		return badRequestError(c, err)
	}

	rollup, err := handler.rollups.GetByID(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	stats, err := handler.rollups.Stats(c.Request().Context(), rollup.Id)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	return c.JSON(http.StatusOK, responses.NewRollupWithStats(storage.RollupWithStats{
		Rollup:      *rollup,
		RollupStats: stats,
	}))
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
//	@Router			/v1/rollup/{id}/namespaces [get]
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
}

func (p *getRollupPagesWithSort) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = desc
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
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{id}/blobs [get]
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
	SeriesName string `example:"tps"        param:"name"      swaggertype:"string"  validate:"required,oneof=blobs_count size"`
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
//	@Param			name		path	string	true	"Series name"					Enums(blobs_count, size)
//	@Param			timeframe	path	string	true	"Timeframe"						Enums(hour, day, month)
//	@Param			from		query	integer	false	"Time from in unix timestamp"	mininum(1)
//	@Param			to			query	integer	false	"Time to in unix timestamp"		mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.HistogramItem
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/{id}/stats/{name}/{timeframe} [get]
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

// Count godoc
//
//	@Summary		Get count of rollups in network
//	@Description	Get count of rollups in network
//	@Tags			rollup
//	@ID				get-rollups-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/v1/rollup/count [get]
func (handler RollupHandler) Count(c echo.Context) error {
	count, err := handler.rollups.Count(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.rollups)
	}
	return c.JSON(http.StatusOK, count)
}

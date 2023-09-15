package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	repo storage.IStats
}

func NewStatsHandler(repo storage.IStats) StatsHandler {
	return StatsHandler{
		repo: repo,
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
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, summary)
}

type histogramRequest struct {
	Table     string `example:"block"      param:"table"     swaggertype:"string"  validate:"required,oneof=block tx event message"`
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
//	@Param					table		path	string	true	"Table name"	Enums(block, tx, event, message)
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

	return c.JSON(http.StatusOK, response)
}

package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type BlockHandler struct {
	block      storage.IBlock
	blockStats storage.IBlockStats
	events     storage.IEvent
}

func NewBlockHandler(block storage.IBlock, blockStats storage.IBlockStats, events storage.IEvent) *BlockHandler {
	return &BlockHandler{
		block:      block,
		blockStats: blockStats,
		events:     events,
	}
}

type getBlockRequest struct {
	Height uint64 `param:"height" validate:"required,min=1"`

	Stats bool `query:"stats" validate:"omitempty"`
}

// Get godoc
//
//	@Summary		Get block info
//	@Description	Get block info
//	@Tags			block
//	@ID				get-block
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Param			stats	query	boolean	false 	"Need join stats for block"
//	@Produce		json
//	@Success		200	{object}	responses.Block
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height} [get]
func (handler *BlockHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	block, err := handler.block.ByHeight(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.NewBlock(block, req.Stats))
}

type blockListRequest struct {
	Limit  uint64 `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `query:"offset" validate:"omitempty,min=0"`
	Sort   string `query:"sort"   validate:"omitempty,oneof=asc desc"`
	Stats  bool   `query:"stats"  validate:"omitempty"`
}

func (p *blockListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

// List godoc
//
//	@Summary		List blocks info
//	@Description	List blocks info
//	@Tags			block
//	@ID				list-block
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order"					Enums(asc, desc)
//	@Param			stats	query	boolean	false 	"Need join stats for block"
//	@Produce		json
//	@Success		200	{array}		responses.Block
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block [get]
func (handler *BlockHandler) List(c echo.Context) error {
	req, err := bindAndValidate[blockListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blocks, err := handler.block.ListWithStats(c.Request().Context(), req.Stats, req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(blocks[i], req.Stats)
	}

	return returnArray(c, response)
}

// GetEvents godoc
//
//	@Summary		Get events from begin and end of block
//	@Description	Get events from begin and end of block
//	@Tags			block
//	@ID				get-block-events
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Event
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/events [get]
func (handler *BlockHandler) GetEvents(c echo.Context) error {
	req, err := bindAndValidate[getBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	events, err := handler.events.ByBlock(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.events); err != nil {
		return err
	}

	response := make([]responses.Event, len(events))
	for i := range events {
		response[i] = responses.NewEvent(events[i])
	}

	return returnArray(c, response)
}

// GetStats godoc
//
//	@Summary		Get block stats by height
//	@Description	Get block stats by height
//	@Tags			block
//	@ID				get-block-stats
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@Success		200	{object}		responses.BlockStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/v1/block/{height}/stats [get]
func (handler *BlockHandler) GetStats(c echo.Context) error {
	req, err := bindAndValidate[getBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	stats, err := handler.blockStats.ByHeight(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.events); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responses.NewBlockStats(stats))
}

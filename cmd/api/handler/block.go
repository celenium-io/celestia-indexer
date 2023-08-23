package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type BlockHandler struct {
	block  storage.IBlock
	events storage.IEvent
}

func NewBlockHandler(block storage.IBlock, events storage.IEvent) *BlockHandler {
	return &BlockHandler{
		block:  block,
		events: events,
	}
}

type getBlockRequest struct {
	Height uint64 `param:"height" validate:"required,min=1"`
}

// Get godoc
// @Summary Get block info
// @Description Get block info
// @Tags block
// @ID get-block
// @Param height path integer true "Block height" minimum(1)
// @Produce  json
// @Success 200 {object} Block
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/block/{height} [get]
func (handler *BlockHandler) Get(c echo.Context) error {
	req := new(getBlockRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}

	block, err := handler.block.ByHeight(c.Request().Context(), req.Height)
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.NewBlock(block))
}

// List godoc
// @Summary List blocks info
// @Description List blocks info
// @Tags block
// @ID list-block
// @Param limit  query integer false "Count of requested entities" mininum(1) maximum(100)
// @Param offset query integer false "Offset" mininum(1)
// @Param sort   query string  false "Sort order" Enums(asc, desc)
// @Produce json
// @Success 200 {array} Block
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/block [get]
func (handler *BlockHandler) List(c echo.Context) error {
	req := new(limitOffsetPagination)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blocks, err := handler.block.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	if err := handleError(c, err, handler.block); err != nil {
		return err
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(*blocks[i])
	}

	return returnArray(c, response)
}

// GetEvents godoc
// @Summary Get events from begin and end of block
// @Description Get events from begin and end of block
// @Tags block
// @ID get-block-events
// @Param height path integer true "Block height" minimum(1)
// @Produce json
// @Success 200 {array} Event
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/block/{height}/events [get]
func (handler *BlockHandler) GetEvents(c echo.Context) error {
	req := new(getBlockRequest)
	if err := c.Bind(req); err != nil {
		return badRequestError(c, err)
	}
	if err := c.Validate(req); err != nil {
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

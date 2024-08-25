// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/pkg/node"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/da"
	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/go-square/square"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type BlockHandler struct {
	block       storage.IBlock
	blockStats  storage.IBlockStats
	events      storage.IEvent
	namespace   storage.INamespace
	blobLogs    storage.IBlobLog
	message     storage.IMessage
	state       storage.IState
	node        node.Api
	indexerName string
}

func NewBlockHandler(
	block storage.IBlock,
	blockStats storage.IBlockStats,
	events storage.IEvent,
	namespace storage.INamespace,
	message storage.IMessage,
	blobLogs storage.IBlobLog,
	state storage.IState,
	node node.Api,
	indexerName string,
) *BlockHandler {
	return &BlockHandler{
		block:       block,
		blockStats:  blockStats,
		events:      events,
		namespace:   namespace,
		blobLogs:    blobLogs,
		message:     message,
		state:       state,
		node:        node,
		indexerName: indexerName,
	}
}

type getBlockByHeightRequest struct {
	Height types.Level `param:"height" validate:"min=0"`
}

type getBlockRequest struct {
	Height types.Level `param:"height" validate:"min=0"`

	Stats bool `query:"stats" validate:"omitempty"`
}

// Get godoc
//
//	@Summary		Get block info
//	@Description	Get block info
//	@Tags			block
//	@ID				get-block
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Param			stats	query	boolean	false	"Need join stats for block"
//	@Produce		json
//	@Success		200	{object}	responses.Block
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height} [get]
func (handler *BlockHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	var block storage.Block
	if req.Stats {
		block, err = handler.block.ByHeightWithStats(c.Request().Context(), req.Height)
	} else {
		block, err = handler.block.ByHeight(c.Request().Context(), req.Height)
	}

	if err != nil {
		return handleError(c, err, handler.block)
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
//	@Param			stats	query	boolean	false	"Need join stats for block"
//	@Produce		json
//	@Success		200	{array}		responses.Block
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block [get]
func (handler *BlockHandler) List(c echo.Context) error {
	req, err := bindAndValidate[blockListRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	var blocks []*storage.Block
	if req.Stats {
		blocks, err = handler.block.ListWithStats(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	} else {
		blocks, err = handler.block.List(c.Request().Context(), req.Limit, req.Offset, pgSort(req.Sort))
	}

	if err != nil {
		return handleError(c, err, handler.block)
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(*blocks[i], req.Stats)
	}

	return returnArray(c, response)
}

type getBlockEvents struct {
	Height types.Level `param:"height" validate:"min=0"`
	Limit  int         `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int         `query:"offset" validate:"omitempty,min=0"`
}

func (p *getBlockEvents) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// GetEvents godoc
//
//	@Summary		Get events from begin and end of block
//	@Description	Get events from begin and end of block
//	@Tags			block
//	@ID				get-block-events
//	@Param			height	path	integer	true	"Block height"					minimum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Event
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/events [get]
func (handler *BlockHandler) GetEvents(c echo.Context) error {
	req, err := bindAndValidate[getBlockEvents](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blockTime, err := handler.block.Time(c.Request().Context(), req.Height)
	if err != nil {
		if handler.block.IsNoRows(err) {
			return returnArray(c, []any{})
		}
		return handleError(c, err, handler.block)
	}

	fltrs := storage.EventFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		Time:   blockTime.UTC(),
	}

	events, err := handler.events.ByBlock(c.Request().Context(), req.Height, fltrs)
	if err != nil {
		return handleError(c, err, handler.block)
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
//	@Success		200	{object}	responses.BlockStats
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/stats [get]
func (handler *BlockHandler) GetStats(c echo.Context) error {
	req, err := bindAndValidate[getBlockByHeightRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	stats, err := handler.blockStats.ByHeight(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	return c.JSON(http.StatusOK, responses.NewBlockStats(stats))
}

// Count godoc
//
//	@Summary		Get count of blocks in network
//	@Description	Get count of blocks in network
//	@Tags			block
//	@ID				get-block-count
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		500	{object}	Error
//	@Router			/block/count [get]
func (handler *BlockHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	return c.JSON(http.StatusOK, state.LastHeight+1) // + genesis block
}

// GetMessages godoc
//
//	@Summary		Get messages contained in the block
//	@Description	Get messages contained in the block
//	@Tags			block
//	@ID				get-block-messages
//	@Param			height				path	integer			true	"Block height"					minimum(1)
//	@Param			limit				query	integer			false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset				query	integer			false	"Offset"						mininum(1)
//	@Param			msg_type			query	types.MsgType	false	"Comma-separated message types list"
//	@Param			excluded_msg_type	query	types.MsgType	false	"Comma-separated message types which should be excluded from list"
//	@Produce		json
//	@Success		200	{array}		responses.Message
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/messages [get]
func (handler *BlockHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[listMessageByBlockRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.MessageListWithTxFilters{
		Limit:                req.Limit,
		Offset:               req.Offset,
		Height:               req.Height,
		MessageTypes:         req.MsgType,
		ExcludedMessageTypes: req.ExcludedMsgType,
	}

	messages, err := handler.message.ListWithTx(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.block)
	}
	response := make([]responses.Message, len(messages))
	for i := range response {
		msg := responses.NewMessageWithTx(messages[i])
		response[i] = msg
	}

	return c.JSON(http.StatusOK, response)
}

type getBlobsForBlock struct {
	Height types.Level `param:"height"  validate:"min=0"`
	Limit  int         `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset int         `query:"offset"  validate:"omitempty,min=0"`
	Sort   string      `query:"sort"    validate:"omitempty,oneof=asc desc"`
	SortBy string      `query:"sort_by" validate:"omitempty,oneof=time size"`
}

func (req *getBlobsForBlock) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// Blobs godoc
//
//	@Summary		List blobs which was pushed in the block
//	@Description	List blobs which was pushed in the block
//	@Tags			block
//	@ID				get-block-blobs
//	@Param			height	path	integer	true	"Block height"									minimum(1)
//	@Param			limit	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"										mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			sort_by	query	string	false	"Sort field. If it's empty internal id is used"	Enums(time, size)
//	@Produce		json
//	@Success		200	{array}		responses.BlobLog
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/blobs [get]
func (handler *BlockHandler) Blobs(c echo.Context) error {
	req, err := bindAndValidate[getBlobsForBlock](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blockTime, err := handler.block.Time(c.Request().Context(), req.Height)
	if err != nil {
		if handler.block.IsNoRows(err) {
			return returnArray(c, []any{})
		}
		return handleError(c, err, handler.block)
	}

	blobs, err := handler.blobLogs.ByHeight(
		c.Request().Context(),
		req.Height,
		storage.BlobLogFilters{
			Limit:  req.Limit,
			Offset: req.Offset,
			Sort:   pgSort(req.Sort),
			SortBy: req.SortBy,
			// using time filters to take certain partition
			From: blockTime,
			To:   blockTime.Add(time.Minute),
		},
	)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}

	response := make([]responses.BlobLog, len(blobs))
	for i := range blobs {
		response[i] = responses.NewBlobLog(blobs[i])
	}
	return returnArray(c, response)
}

// BlobsCount godoc
//
//	@Summary		Count of blobs which was pushed by transaction
//	@Description	Count of blobs which was pushed by transaction
//	@Tags			block
//	@ID				block-blobs-count
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@Success		200	{integer}	uint64
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/blobs/count [get]
func (handler *BlockHandler) BlobsCount(c echo.Context) error {
	req, err := bindAndValidate[getBlockByHeightRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	count, err := handler.blobLogs.CountByHeight(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.blobLogs)
	}

	return c.JSON(http.StatusOK, count)
}

// BlockODS godoc
//
//	@Summary		ODS for block
//	@Description	ODS for block
//	@Tags			block
//	@ID				block-ods
//	@Param			height	path	integer	true	"Block height"	minimum(1)
//	@Produce		json
//	@x-internal		true
//	@Success		200	{object}	responses.ODS
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/block/{height}/ods [get]
func (handler *BlockHandler) BlockODS(c echo.Context) error {
	req, err := bindAndValidate[getBlockByHeightRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	blockStats, err := handler.blockStats.ByHeight(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.block)
	}

	if blockStats.TxCount == 0 {
		return c.JSON(http.StatusOK, responses.ODS{
			Width: 1,
			Items: []responses.ODSItem{
				{
					From:      []uint{0, 0},
					To:        []uint{0, 0},
					Namespace: "//////////////////////////////////////4=",
					Type:      responses.TailPaddingNamespace,
				},
			},
		})
	}

	block, err := handler.node.Block(c.Request().Context(), req.Height)
	if err != nil {
		return handleError(c, err, handler.block)
	}

	dataSquare, err := square.Construct(
		block.Block.Data.Txs.ToSliceOfBytes(),
		appconsts.SquareSizeUpperBound(0),
		appconsts.SubtreeRootThreshold(0),
	)
	if err != nil {
		return internalServerError(c, err)
	}

	eds, err := da.ExtendShares(shares.ToBytes(dataSquare))
	if err != nil {
		return internalServerError(c, err)
	}

	ods, err := responses.NewODS(eds)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, ods)
}

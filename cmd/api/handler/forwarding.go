// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/hex"
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/cmd/api/hyperlane"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

// ForwardingsHandler -
type ForwardingsHandler struct {
	forwardings storage.IForwarding
	address     storage.IAddress
	txs         storage.ITx
	chainStore  hyperlane.IChainStore
}

func NewForwardingsHandler(
	forwardings storage.IForwarding,
	address storage.IAddress,
	txs storage.ITx,
	chainStore hyperlane.IChainStore,
) ForwardingsHandler {
	return ForwardingsHandler{
		forwardings: forwardings,
		address:     address,
		txs:         txs,
		chainStore:  chainStore,
	}
}

type listForwardingsRequest struct {
	Limit   int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset"  validate:"omitempty,min=0"`
	Sort    string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	TxHash  string `query:"tx_hash" validate:"omitempty,hexadecimal,len=64"`
	Address string `query:"address" validate:"omitempty,address"`
	Height  uint64 `query:"height"  validate:"omitempty,min=1"`

	From int64 `example:"1692892095" query:"from" swaggertype:"integer" validate:"omitempty,min=1"`
	To   int64 `example:"1692892095" query:"to"   swaggertype:"integer" validate:"omitempty,min=1"`
}

func (req *listForwardingsRequest) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

func (req *listForwardingsRequest) toFilters(ctx context.Context, address storage.IAddress, txs storage.ITx) (storage.ForwardingFilter, error) {
	filters := storage.ForwardingFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   pgSort(req.Sort),
	}
	if req.Height > 0 {
		filters.Height = &req.Height
	}
	if req.TxHash != "" {
		txHash, err := hex.DecodeString(req.TxHash)
		if err != nil {
			return filters, err
		}
		txId, txTime, err := txs.IdAndTimeByHash(ctx, txHash)
		if err != nil {
			return filters, err
		}
		filters.TxId = &txId
		filters.From = txTime
	}
	if req.Address != "" {
		addressId, err := address.IdByAddress(ctx, req.Address)
		if err != nil {
			return filters, err
		}
		filters.AddressId = &addressId
	}

	return filters, nil
}

// List godoc
//
//	@Summary		List forwarding events
//	@Description	Returns a paginated list of forwarding events. Forwarding events represent cross-domain token transfers
//	@Description	where tokens are forwarded from a Celestia address to a destination address on another domain.
//	@Description	Results can be filtered by transaction hash, address, or block height.
//	@Tags			forwarding
//	@ID				list-forwarding
//	@Param			limit	query	integer	false	"Count of requested entities"			mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset for pagination"					mininum(0)
//	@Param			sort	query	string	false	"Sort order. Default: desc"				Enums(asc, desc)
//	@Param			tx_hash	query	string	false	"Filter by transaction hash (hex)"		minlength(64)	maxlength(64)
//	@Param			address	query	string	false	"Filter by Celestia address"			minlength(47)	maxlength(47)
//	@Param			height	query	integer	false	"Filter by block height"				mininum(1)
//	@Param			from	query	integer	false	"Filter by start time (Unix timestamp)"	minimum(1)
//	@Param			to		query	integer	false	"Filter by end time (Unix timestamp)"	minimum(1)
//	@Success		200	{array}		responses.Forwarding
//	@Produce		json
//	@Success		200	{array}		responses.Forwarding
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/forwarding [get]
func (handler *ForwardingsHandler) List(c echo.Context) error {
	req, err := bindAndValidate[listForwardingsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	filters, err := req.toFilters(c.Request().Context(), handler.address, handler.txs)
	if err != nil {
		return badRequestError(c, err)
	}

	forwardings, err := handler.forwardings.Filter(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.forwardings)
	}
	response := make([]responses.Forwarding, len(forwardings))
	for i := range forwardings {
		response[i] = responses.NewForwarding(forwardings[i], handler.chainStore)
	}
	return returnArray(c, response)
}

// Get godoc
//
//	@Summary		Get forwarding event by ID
//	@Description	Returns a single forwarding event by its internal ID. The response includes details about
//	@Description	the cross-domain token transfer such as destination domain, destination address,
//	@Description	forwarding address, success/failed counts, and the list of individual transfers.
//	@Tags			forwarding
//	@ID				get-forwarding
//	@Param			id	path	integer	true	"Internal forwarding event ID"	mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.Forwarding
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/forwarding/{id} [get]
func (handler *ForwardingsHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getById](c)
	if err != nil {
		return badRequestError(c, err)
	}

	forwarding, prevTime, err := handler.forwardings.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.forwardings)
	}

	inputs, err := handler.forwardings.Inputs(c.Request().Context(), forwarding.AddressId, prevTime, forwarding.Time)
	if err != nil {
		return handleError(c, err, handler.forwardings)
	}

	response := responses.NewForwarding(forwarding, handler.chainStore)

	response.Inputs = make([]responses.ForwardingInput, len(inputs))
	for i, input := range inputs {
		response.Inputs[i] = responses.NewForwardingInputFromHyperlaneTransfer(input, handler.chainStore)
	}
	return c.JSON(http.StatusOK, response)
}

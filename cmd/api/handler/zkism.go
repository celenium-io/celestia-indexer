// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

type ZkISMHandler struct {
	zkism   storage.IZkISM
	address storage.IAddress
	txs     storage.ITx
}

func NewZkISMHandler(zkism storage.IZkISM, address storage.IAddress, txs storage.ITx) *ZkISMHandler {
	return &ZkISMHandler{
		zkism:   zkism,
		address: address,
		txs:     txs,
	}
}

type listZkISMRequest struct {
	Limit   int    `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset"  validate:"omitempty,min=0"`
	Sort    string `query:"sort"    validate:"omitempty,oneof=asc desc"`
	TxHash  string `query:"tx_hash" validate:"omitempty,hexadecimal,len=64"`
	Address string `query:"address" validate:"omitempty,address"`
}

func (req *listZkISMRequest) toFilter(ctx context.Context, address storage.IAddress, txs storage.ITx) (storage.ZkISMFilter, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}

	filter := storage.ZkISMFilter{
		Sort:   sdk.SortOrder(req.Sort),
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if req.TxHash != "" {
		hash, err := hex.DecodeString(req.TxHash)
		if err != nil {
			return filter, err
		}
		txId, _, err := txs.IdAndTimeByHash(ctx, hash)
		if err != nil {
			return filter, err
		}
		filter.TxId = &txId
	}
	if req.Address != "" {
		addressId, err := address.IdByAddress(ctx, req.Address)
		if err != nil {
			return filter, err
		}
		filter.CreatorId = &addressId
	}

	return filter, nil
}

// List godoc
//
//	@Summary		List ZK Interchain Security Modules
//	@Description	Returns a paginated list of ZK Interchain Security Modules (ZK ISMs). ZK ISMs use
//	@Description	Groth16 zero-knowledge proofs to trustlessly verify cross-chain messages from external
//	@Description	chains. Results can be filtered by transaction hash or creator address.
//	@Tags			hyperlane
//	@ID				list-zkism
//	@Param			limit	 query	integer	false	"Count of requested entities"				mininum(1)	maximum(100)
//	@Param			offset	 query	integer	false	"Offset for pagination"						mininum(0)
//	@Param			sort	 query	string	false	"Sort order. Default: desc"					Enums(asc, desc)
//	@Param			tx_hash	 query	string	false	"Filter by transaction hash (hex)"			minlength(64)	maxlength(64)
//	@Param			address	 query	string	false	"Filter by creator Celestia address"		minlength(47)	maxlength(47)
//	@Produce		json
//	@Success		200	{array}		responses.ZkISM
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/zkism [get]
func (h *ZkISMHandler) List(c echo.Context) error {
	req, err := bindAndValidate[listZkISMRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	filter, err := req.toFilter(c.Request().Context(), h.address, h.txs)
	if err != nil {
		return badRequestError(c, err)
	}

	items, err := h.zkism.List(c.Request().Context(), filter)
	if err != nil {
		return handleError(c, err, h.address)
	}

	response := make([]responses.ZkISM, len(items))
	for i := range items {
		response[i] = responses.NewZkISM(items[i])
	}
	return returnArray(c, response)
}

type getZkISMRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

// Get godoc
//
//	@Summary		Get ZK ISM by id
//	@Description	Returns a single ZK Interchain Security Module by its internal id. The response includes
//	@Description	the current trusted state, state root, merkle tree address, and verifier key commitments
//	@Description	used for ZK proof verification.
//	@Tags			hyperlane
//	@ID				get-zkism
//	@Param			id	path	integer	true	"Internal ZK ISM identity"	mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.ZkISM
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/zkism/{id} [get]
func (h *ZkISMHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getZkISMRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	ism, err := h.zkism.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, h.address)
	}

	return c.JSON(http.StatusOK, responses.NewZkISM(ism))
}

type listZkISMHistoryRequest struct {
	Id      uint64 `param:"id"           validate:"required,min=1"`
	Limit   int    `query:"limit"        validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset"       validate:"omitempty,min=0"`
	Sort    string `query:"sort"         validate:"omitempty,oneof=asc desc"`
	TxHash  string `query:"tx_hash"      validate:"omitempty,hexadecimal,len=64"`
	Address string `query:"address"      validate:"omitempty,address"`
	From    int64  `example:"1692892095" query:"from"                            swaggertype:"integer" validate:"omitempty,min=1"`
	To      int64  `example:"1692892095" query:"to"                              swaggertype:"integer" validate:"omitempty,min=1"`
}

func (req *listZkISMHistoryRequest) toFilter(ctx context.Context, address storage.IAddress, txs storage.ITx) (storage.ZkISMUpdatesFilter, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}

	filter := storage.ZkISMUpdatesFilter{
		Sort:   sdk.SortOrder(req.Sort),
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if req.From > 0 {
		filter.From = time.Unix(req.From, 0)
	}
	if req.To > 0 {
		filter.To = time.Unix(req.To, 0)
	}
	if req.TxHash != "" {
		hash, err := hex.DecodeString(req.TxHash)
		if err != nil {
			return filter, err
		}
		txId, _, err := txs.IdAndTimeByHash(ctx, hash)
		if err != nil {
			return filter, err
		}
		filter.TxId = &txId
	}
	if req.Address != "" {
		addressId, err := address.IdByAddress(ctx, req.Address)
		if err != nil {
			return filter, err
		}
		filter.SignerId = &addressId
	}

	return filter, nil
}

// GetUpdates godoc
//
//	@Summary		Get ZK ISM state update history
//	@Description	Returns the history of state updates for a given ZK ISM. Each update corresponds to a
//	@Description	MsgUpdateInterchainSecurityModule transaction that advanced the trusted state of the ISM
//	@Description	via a Groth16 state-transition ZK proof. Results can be filtered by signer address,
//	@Description	transaction hash, or time range.
//	@Tags			hyperlane
//	@ID				get-zkism-updates
//	@Param			id		 path	integer	true	"Internal ZK ISM identity"					mininum(1)
//	@Param			limit	 query	integer	false	"Count of requested entities"				mininum(1)	maximum(100)
//	@Param			offset	 query	integer	false	"Offset for pagination"						mininum(0)
//	@Param			sort	 query	string	false	"Sort order. Default: desc"					Enums(asc, desc)
//	@Param			tx_hash	 query	string	false	"Filter by transaction hash (hex)"			minlength(64)	maxlength(64)
//	@Param			address	 query	string	false	"Filter by signer Celestia address"			minlength(47)	maxlength(47)
//	@Param			from	 query	integer	false	"Filter by start time (Unix timestamp)"		minimum(1)
//	@Param			to		 query	integer	false	"Filter by end time (Unix timestamp)"		minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.ZkISMUpdate
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/zkism/{id}/updates [get]
func (h *ZkISMHandler) GetUpdates(c echo.Context) error {
	req, err := bindAndValidate[listZkISMHistoryRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	filter, err := req.toFilter(c.Request().Context(), h.address, h.txs)
	if err != nil {
		return badRequestError(c, err)
	}

	items, err := h.zkism.Updates(c.Request().Context(), req.Id, filter)
	if err != nil {
		return handleError(c, err, h.address)
	}

	response := make([]responses.ZkISMUpdate, len(items))
	for i := range items {
		response[i] = responses.NewZkISMUpdate(items[i])
	}
	return returnArray(c, response)
}

// GetMessages godoc
//
//	@Summary		Get ZK ISM authorized messages
//	@Description	Returns the list of Hyperlane messages that were authorized via a given ZK ISM through
//	@Description	MsgSubmitMessages transactions. Each entry contains the message id, state root used for
//	@Description	the membership proof, and the signer who submitted the authorization. Results can be
//	@Description	filtered by signer address, transaction hash, or time range.
//	@Tags			hyperlane
//	@ID				get-zkism-messages
//	@Param			id		 path	integer	true	"Internal ZK ISM identity"					mininum(1)
//	@Param			limit	 query	integer	false	"Count of requested entities"				mininum(1)	maximum(100)
//	@Param			offset	 query	integer	false	"Offset for pagination"						mininum(0)
//	@Param			sort	 query	string	false	"Sort order. Default: desc"					Enums(asc, desc)
//	@Param			tx_hash	 query	string	false	"Filter by transaction hash (hex)"			minlength(64)	maxlength(64)
//	@Param			address	 query	string	false	"Filter by signer Celestia address"			minlength(47)	maxlength(47)
//	@Param			from	 query	integer	false	"Filter by start time (Unix timestamp)"		minimum(1)
//	@Param			to		 query	integer	false	"Filter by end time (Unix timestamp)"		minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.ZkISMMessage
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/zkism/{id}/messages [get]
func (h *ZkISMHandler) GetMessages(c echo.Context) error {
	req, err := bindAndValidate[listZkISMHistoryRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	filter, err := req.toFilter(c.Request().Context(), h.address, h.txs)
	if err != nil {
		return badRequestError(c, err)
	}

	items, err := h.zkism.Messages(c.Request().Context(), req.Id, filter)
	if err != nil {
		return handleError(c, err, h.address)
	}

	response := make([]responses.ZkISMMessage, len(items))
	for i := range items {
		response[i] = responses.NewZkISMMessage(items[i])
	}
	return returnArray(c, response)
}

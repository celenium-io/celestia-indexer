// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/hex"
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type SignalHandler struct {
	signals   storage.ISignalVersion
	upgrades  storage.IUpgrade
	validator storage.IValidator
	tx        storage.ITx
	address   storage.IAddress
}

func NewSignalHandler(
	signals storage.ISignalVersion,
	upgrades storage.IUpgrade,
	validator storage.IValidator,
	tx storage.ITx,
	address storage.IAddress,
) *SignalHandler {
	return &SignalHandler{
		signals:   signals,
		upgrades:  upgrades,
		validator: validator,
		tx:        tx,
		address:   address,
	}
}

type signalsRequest struct {
	Limit       int    `example:"10"                                                               query:"limit"        swaggertype:"integer" validate:"omitempty,min=1,max=100"`
	Offset      int    `example:"0"                                                                query:"offset"       swaggertype:"integer" validate:"omitempty,min=0"`
	Version     uint64 `example:"7"                                                                query:"version"      swaggertype:"integer" validate:"omitempty"`
	ValidatorId uint64 `example:"1488"                                                             query:"validator_id" swaggertype:"string"  validate:"omitempty"`
	TxHash      string `example:"97589d917f13c7d1d5b01dcc0a3df84c8b4337e47dae492e03c274cc77ded173" query:"tx_hash"      swaggertype:"string"  validate:"omitempty"`
	From        int64  `example:"1692892095"                                                       query:"from"         swaggertype:"integer" validate:"omitempty,min=1"`
	To          int64  `example:"1692892095"                                                       query:"to"           swaggertype:"integer" validate:"omitempty,min=1"`
	Sort        string `example:"asc"                                                              query:"sort"         swaggertype:"string"  validate:"omitempty,oneof=asc desc"`
}

func (req *signalsRequest) ToFilters(
	ctx context.Context,
	txs storage.ITx) (storage.ListSignalsFilter, error) {
	var filters = storage.ListSignalsFilter{
		Limit:       req.Limit,
		Offset:      req.Offset,
		ValidatorId: req.ValidatorId,
		Version:     req.Version,
		TxHash:      req.TxHash,
		Sort:        desc,
	}

	if req.Limit > 0 {
		filters.Limit = req.Limit
	}
	if req.Sort != "" {
		filters.Sort = pgSort(req.Sort)
	}

	if req.TxHash != "" {
		hash, err := hex.DecodeString(req.TxHash)
		if err != nil {
			return filters, err
		}

		tx, err := txs.ByHash(ctx, hash)
		if err != nil {
			return filters, err
		}

		filters.TxId = &tx.Id
	}

	return filters, nil
}

// List godoc
//
//		@Summary		List signals
//		@Description	List signals
//		@Tags			signal
//		@ID				list-signal
//	    @Param			version	query	integer	false	"Version"
//	    @Param			validator_id	query	integer	false	"Validator internal id"
//	    @Param			tx_hash	query	string	false	"Transaction hash"
//		@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//		@Param			offset	query	integer	false	"Offset"						mininum(1)
//		@Param			from	query	integer	false	"Time from in unix timestamp"	mininum(1)
//		@Param			to		query	integer	false	"Time to in unix timestamp"		mininum(1)
//		@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//		@Produce		json
//		@Success		200	{array}		responses.SignalVersion
//		@Failure		400	{object}	Error
//		@Failure		500	{object}	Error
//		@Router			/signal [get]
func (handler *SignalHandler) List(c echo.Context) error {
	req, err := bindAndValidate[signalsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	filters, err := req.ToFilters(c.Request().Context(), handler.tx)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	signals, err := handler.signals.List(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	response := make([]responses.SignalVersion, len(signals))
	for i := range signals {
		response[i] = responses.NewSignalVersion(signals[i])
	}

	return returnArray(c, response)
}

type upgradesRequest struct {
	Limit  int    `example:"10"                                                               query:"limit"   swaggertype:"integer" validate:"omitempty,min=1,max=100"`
	Offset int    `example:"0"                                                                query:"offset"  swaggertype:"integer" validate:"omitempty,min=0"`
	Height uint64 `example:"12345678"                                                         query:"height"  swaggertype:"integer" validate:"omitempty"`
	TxHash string `example:"97589d917f13c7d1d5b01dcc0a3df84c8b4337e47dae492e03c274cc77ded173" query:"tx_hash" swaggertype:"string"  validate:"omitempty"`
	Signer string `example:"celestia1ps2778x42p833xesk7jh6vy0qu8485e3pz8g72"                  query:"signer"  swaggertype:"string"  validate:"omitempty"`
	Sort   string `example:"asc"                                                              query:"sort"    swaggertype:"string"  validate:"omitempty,oneof=asc desc"`
}

func (req *upgradesRequest) ToFilters(
	ctx context.Context,
	address storage.IAddress,
	txs storage.ITx) (storage.ListUpgradesFilter, error) {
	var filters = storage.ListUpgradesFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		Height: req.Height,
		TxHash: req.TxHash,
		Signer: req.Signer,
		Sort:   desc,
	}

	if req.Limit > 0 {
		filters.Limit = req.Limit
	}
	if req.Sort != "" {
		filters.Sort = pgSort(req.Sort)
	}
	if req.Signer != "" {
		addrId, err := address.IdByAddress(ctx, req.Signer)
		if err != nil {
			return filters, err
		}
		filters.SignerId = &addrId
	}

	if req.TxHash != "" {
		hash, err := hex.DecodeString(req.TxHash)
		if err != nil {
			return filters, err
		}

		tx, err := txs.ByHash(ctx, hash)
		if err != nil {
			return filters, err
		}

		filters.TxId = &tx.Id
	}

	return filters, nil
}

// Upgrades godoc
//
//	@Summary		List upgrades
//	@Description	List upgrades
//	@Tags			signal
//	@ID				list-upgrades
//	@Param			height	query	integer	false	"Number of block"
//	@Param			tx_hash	query	string	false	"Transaction hash"
//	@Param			signer	query	string	false	"Signer address"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			sort	query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Produce		json
//	@Success		200	{array}		responses.Upgrade
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/signal/upgrade [get]
func (handler *SignalHandler) Upgrades(c echo.Context) error {
	req, err := bindAndValidate[upgradesRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	filters, err := req.ToFilters(c.Request().Context(), handler.address, handler.tx)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	upgrades, err := handler.upgrades.List(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.tx)
	}

	response := make([]responses.Upgrade, len(upgrades))
	for i := range upgrades {
		response[i] = responses.NewUpgrade(upgrades[i])
	}

	return returnArray(c, response)
}

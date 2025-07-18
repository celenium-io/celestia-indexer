// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

type IbcHandler struct {
	clients   storage.IIbcClient
	conns     storage.IIbcConnection
	channels  storage.IIbcChannel
	transfers storage.IIbcTransfer
	address   storage.IAddress
}

func NewIbcHandler(
	clients storage.IIbcClient,
	conns storage.IIbcConnection,
	channels storage.IIbcChannel,
	transfers storage.IIbcTransfer,
	address storage.IAddress,
) *IbcHandler {
	return &IbcHandler{
		clients:   clients,
		conns:     conns,
		channels:  channels,
		transfers: transfers,
		address:   address,
	}
}

type getIbcClientRequest struct {
	Id string `param:"id" validate:"required"`
}

// Get godoc
//
//	@Summary		Get ibc client info
//	@Description	Get ibc client info
//	@Tags			ibc
//	@ID				get-ibc-client
//	@Param			id	path	string	true	"IBC client id"
//	@Produce		json
//	@Success		200	{object}	responses.IbcClient
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/client/{id} [get]
func (handler *IbcHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[getIbcClientRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	client, err := handler.clients.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewIbcClient(client))
}

type getIbcClientsRequest struct {
	Limit   int    `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset"   validate:"omitempty,min=0"`
	Sort    string `query:"sort"     validate:"omitempty,oneof=asc desc"`
	ChainId string `query:"chain_id" validate:"omitempty"`
	Creator string `query:"creator"  validate:"omitempty,address"`
}

func (req *getIbcClientsRequest) ToFilters(ctx context.Context, address storage.IAddress) (storage.ListIbcClientsFilters, error) {
	var filters = storage.ListIbcClientsFilters{
		Limit:   10,
		Sort:    desc,
		Offset:  req.Offset,
		ChainId: req.ChainId,
	}

	if req.Limit > 0 {
		filters.Limit = req.Limit
	}
	if req.Sort != "" {
		filters.Sort = pgSort(req.Sort)
	}
	if req.Creator != "" {
		addrId, err := address.IdByAddress(ctx, req.Creator)
		if err != nil {
			return filters, err
		}
		filters.CreatorId = &addrId
	}

	return filters, nil
}

// List godoc
//
//	@Summary		Get ibc clients info
//	@Description	Get ibc clients info
//	@Tags			ibc
//	@ID				get-ibc-clients
//	@Param			limit		query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"										mininum(1)
//	@Param			sort		query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			chain_id	query	string	false	"Chain id"
//	@Param			creator		query	string	false	"Creator address"						    	minlength(47)	maxlength(47)
//	@Produce		json
//	@Success		200	{array}	responses.IbcClient
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/client [get]
func (handler *IbcHandler) List(c echo.Context) error {
	req, err := bindAndValidate[getIbcClientsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	filters, err := req.ToFilters(c.Request().Context(), handler.address)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	clients, err := handler.clients.List(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.IbcClient, len(clients))
	for i := range clients {
		response[i] = responses.NewIbcClient(clients[i])
	}
	return returnArray(c, response)
}

type getIbcConnectionRequest struct {
	Id string `param:"id" validate:"required"`
}

// GetConnection godoc
//
//	@Summary		Get ibc connection info
//	@Description	Get ibc client info
//	@Tags			ibc
//	@ID				get-ibc-conn
//	@Param			id	path	string	true	"IBC connection id"
//	@Produce		json
//	@Success		200	{object}	responses.IbcConnection
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/connection/{id} [get]
func (handler *IbcHandler) GetConnection(c echo.Context) error {
	req, err := bindAndValidate[getIbcConnectionRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	conn, err := handler.conns.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewIbcConnection(conn))
}

type getIbcConnsRequest struct {
	Limit    int    `query:"limit"     validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset"    validate:"omitempty,min=0"`
	Sort     string `query:"sort"      validate:"omitempty,oneof=asc desc"`
	ClientId string `query:"client_id" validate:"omitempty"`
}

func (req *getIbcConnsRequest) SetDefault() {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// ListConnections godoc
//
//	@Summary		Get ibc connections info
//	@Description	Get ibc connections info
//	@Tags			ibc
//	@ID				get-ibc-conns
//	@Param			limit	    query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	    query	integer	false	"Offset"										mininum(1)
//	@Param			sort	    query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			client_id	query	string	false	"Client id"
//	@Produce		json
//	@Success		200	{array}	responses.IbcConnection
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/connection [get]
func (handler *IbcHandler) ListConnections(c echo.Context) error {
	req, err := bindAndValidate[getIbcConnsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	conns, err := handler.conns.List(c.Request().Context(), storage.ListConnectionFilters{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Sort:     pgSort(req.Sort),
		ClientId: req.ClientId,
	})
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.IbcConnection, len(conns))
	for i := range conns {
		response[i] = responses.NewIbcConnection(conns[i])
	}
	return returnArray(c, response)
}

type getIbcChannelRequest struct {
	Id string `param:"id" validate:"required"`
}

// GetChannel godoc
//
//	@Summary		Get ibc channel info
//	@Description	Get ibc channel info
//	@Tags			ibc
//	@ID				get-ibc-channel
//	@Param			id	path	string	true	"IBC channel id"
//	@Produce		json
//	@Success		200	{object}	responses.IbcChannel
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/channel/{id} [get]
func (handler *IbcHandler) GetChannel(c echo.Context) error {
	req, err := bindAndValidate[getIbcChannelRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	channel, err := handler.channels.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewIbcChannel(channel))
}

type getIbcChannelsRequest struct {
	Limit        int                    `query:"limit"         validate:"omitempty,min=1,max=100"`
	Offset       int                    `query:"offset"        validate:"omitempty,min=0"`
	Sort         string                 `query:"sort"          validate:"omitempty,oneof=asc desc"`
	ClientId     string                 `query:"client_id"     validate:"omitempty"`
	ConnectionId string                 `query:"connection_id" validate:"omitempty"`
	Status       types.IbcChannelStatus `query:"status"        validate:"omitempty,ibc_channel_status"`
}

func (req *getIbcChannelsRequest) SetDefault() {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// ListChannels godoc
//
//	@Summary		Get ibc channels info
//	@Description	Get ibc channels info
//	@Tags			ibc
//	@ID				get-ibc-channels
//	@Param			limit	        query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	        query	integer	false	"Offset"										mininum(1)
//	@Param			sort	        query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			client_id	    query	string	false	"Client id"
//	@Param			connection_id	query	string	false	"Connection id"
//	@Param			status	        query	string	false	"Channel status"					        	Enums(initialization, opened, closed)
//	@Produce		json
//	@Success		200	{array}	responses.IbcChannel
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/channel [get]
func (handler *IbcHandler) ListChannels(c echo.Context) error {
	req, err := bindAndValidate[getIbcChannelsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	channels, err := handler.channels.List(c.Request().Context(), storage.ListChannelFilters{
		Limit:        req.Limit,
		Offset:       req.Offset,
		Sort:         pgSort(req.Sort),
		ClientId:     req.ClientId,
		Status:       req.Status,
		ConnectionId: req.ConnectionId,
	})
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.IbcChannel, len(channels))
	for i := range channels {
		response[i] = responses.NewIbcChannel(channels[i])
	}
	return returnArray(c, response)
}

type getIbcTransfersRequest struct {
	Limit     int    `query:"limit"      validate:"omitempty,min=1,max=100"`
	Offset    int    `query:"offset"     validate:"omitempty,min=0"`
	Sort      string `query:"sort"       validate:"omitempty,oneof=asc desc"`
	ChannelId string `query:"channel_id" validate:"omitempty"`
	ChainId   string `query:"chain_id"   validate:"omitempty"`
	Receiver  string `query:"receiver"   validate:"omitempty,address"`
	Sender    string `query:"sender"     validate:"omitempty,address"`
	Address   string `query:"address"    validate:"omitempty,address"`
}

func (req *getIbcTransfersRequest) SetDefault() {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

// ListTransfers godoc
//
//	@Summary		Get ibc transfers info
//	@Description	Get ibc transfers info
//	@Tags			ibc
//	@ID				get-ibc-transfers
//	@Param			limit	    	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	    	query	integer	false	"Offset"										mininum(1)
//	@Param			sort	   		query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	@Param			channel_id		query	string	false	"Channel id"
//	@Param			chain_id		query	string	false	"Chain id"
//	@Param			receiver		query	string	false	"Receiver address"								minlength(47)	maxlength(47)
//	@Param			sender			query	string	false	"Sender address"								minlength(47)	maxlength(47)
//	@Param			address			query	string	false	"Address: receiver or sender"					minlength(47)	maxlength(47)
//	@Produce		json
//	@Success		200	{array}	responses.IbcTransfer
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/ibc/transfer [get]
func (handler *IbcHandler) ListTransfers(c echo.Context) error {
	req, err := bindAndValidate[getIbcTransfersRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	fltrs := storage.ListIbcTransferFilters{
		Limit:     req.Limit,
		Offset:    req.Offset,
		Sort:      pgSort(req.Sort),
		ChannelId: req.ChannelId,
	}
	if req.Address != "" {
		id, err := handler.address.IdByAddress(c.Request().Context(), req.Address)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		fltrs.AddressId = &id
	}
	if req.Sender != "" {
		id, err := handler.address.IdByAddress(c.Request().Context(), req.Sender)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		fltrs.SenderId = &id
	}
	if req.Receiver != "" {
		id, err := handler.address.IdByAddress(c.Request().Context(), req.Receiver)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		fltrs.ReceiverId = &id
	}
	if req.ChainId != "" {
		clients, err := handler.clients.ByChainId(c.Request().Context(), req.ChainId)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		if len(clients) == 0 {
			return returnArray(c, []any{})
		}
		conns, err := handler.conns.IdsByClients(c.Request().Context(), clients...)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		fltrs.ConnectionIds = conns
	}

	transfers, err := handler.transfers.List(c.Request().Context(), fltrs)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.IbcTransfer, len(transfers))
	for i := range transfers {
		response[i] = responses.NewIbcTransfer(transfers[i])
	}
	return returnArray(c, response)
}

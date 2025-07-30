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
	storageTypes "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HyperlaneHandler struct {
	mailbox    storage.IHLMailbox
	tokens     storage.IHLToken
	transfers  storage.IHLTransfer
	txs        storage.ITx
	address    storage.IAddress
	chainStore hyperlane.IChainStore
}

func NewHyperlaneHandler(
	mailbox storage.IHLMailbox,
	tokens storage.IHLToken,
	transfers storage.IHLTransfer,
	txs storage.ITx,
	address storage.IAddress,
	chainStore hyperlane.IChainStore,
) *HyperlaneHandler {
	return &HyperlaneHandler{
		mailbox:    mailbox,
		tokens:     tokens,
		transfers:  transfers,
		txs:        txs,
		address:    address,
		chainStore: chainStore,
	}
}

type getHyperlaneMailboxRequest struct {
	Id string `param:"id" validate:"required,hexadecimal"`
}

// GetMailbox godoc
//
//	@Summary		Get hyperlane mailbox info
//	@Description	Get hyperlane mailbox info
//	@Tags			hyperlane
//	@ID				get-hyperlane-mailbox
//	@Param			id	path	string	true	"Hyperlane mailbox id"
//	@Produce		json
//	@Success		200	{object}	responses.HyperlaneMailbox
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/mailbox/{id} [get]
func (handler *HyperlaneHandler) GetMailbox(c echo.Context) error {
	req, err := bindAndValidate[getHyperlaneMailboxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	mailbox, err := handler.mailbox.ByHash(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewHyperlaneMailbox(mailbox))
}

type listHyperlaneMailboxRequest struct {
	Limit  int `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int `query:"offset" validate:"omitempty,min=0"`
}

func (req *listHyperlaneMailboxRequest) SetDefault() {
	if req.Limit <= 0 {
		req.Limit = 10
	}
}

// ListMailboxes godoc
//
//	@Summary		List hyperlane mailboxes info
//	@Description	List hyperlane mailboxes info
//	@Tags			hyperlane
//	@ID				list-hyperlane-mailbox
//	@Param			limit	query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"										mininum(1)
//	@Produce		json
//	@Success		200	{array}	responses.HyperlaneMailbox
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/mailbox [get]
func (handler *HyperlaneHandler) ListMailboxes(c echo.Context) error {
	req, err := bindAndValidate[listHyperlaneMailboxRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	mailbox, err := handler.mailbox.List(c.Request().Context(), req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.HyperlaneMailbox, len(mailbox))
	for i := range mailbox {
		response[i] = responses.NewHyperlaneMailbox(mailbox[i])
	}
	return returnArray(c, response)
}

type getHyperlaneTokenRequest struct {
	Id string `param:"id" validate:"required,hexadecimal"`
}

// GetToken godoc
//
//	@Summary		Get hyperlane token info
//	@Description	Get hyperlane token info
//	@Tags			hyperlane
//	@ID				get-hyperlane-token
//	@Param			id	path	string	true	"Hyperlane token id"
//	@Produce		json
//	@Success		200	{object}	responses.HyperlaneToken
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/token/{id} [get]
func (handler *HyperlaneHandler) GetToken(c echo.Context) error {
	req, err := bindAndValidate[getHyperlaneTokenRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return badRequestError(c, err)
	}

	token, err := handler.tokens.ByHash(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewHyperlaneToken(token))
}

type listHyperlaneTokenRequest struct {
	Limit   int         `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  int         `query:"offset"  validate:"omitempty,min=0"`
	Sort    string      `query:"sort"    validate:"omitempty,oneof=asc desc"`
	Owner   string      `query:"owner"   validate:"omitempty,address"`
	Mailbox string      `query:"mailbox" validate:"omitempty,hexadecimal"`
	Type    StringArray `query:"type"    validate:"omitempty,dive,hl_token_type"`
}

func (req *listHyperlaneTokenRequest) ToFilters(ctx context.Context, address storage.IAddress, mailbox storage.IHLMailbox) (storage.ListHyperlaneTokens, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
	filters := storage.ListHyperlaneTokens{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   pgSort(req.Sort),
	}

	if req.Mailbox != "" {
		id, err := hex.DecodeString(req.Mailbox)
		if err != nil {
			return filters, errors.Wrapf(err, "decoding mailbox id: %s", req.Mailbox)
		}
		mbx, err := mailbox.ByHash(ctx, id)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving mailbox by id: %x", id)
		}
		filters.MailboxId = mbx.Id
	}
	if req.Owner != "" {
		_, hash, err := types.Address(req.Owner).Decode()
		if err != nil {
			return filters, errors.Wrapf(err, "decoding owner address: %s", req.Owner)
		}
		addr, err := address.ByHash(ctx, hash)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving address by hash: %x", hash)
		}
		filters.OwnerId = addr.Id
	}
	if len(req.Type) > 0 {
		filters.Type = make([]storageTypes.HLTokenType, len(req.Type))
		for i := range req.Type {
			filters.Type[i] = storageTypes.HLTokenType(req.Type[i])
		}
	}

	return filters, nil
}

// ListTokens godoc
//
//	@Summary		List hyperlane tokens info
//	@Description	List hyperlane tokens info
//	@Tags			hyperlane
//	@ID				list-hyperlane-tokens
//	@Param			limit	query	integer	false	"Count of requested entities"				mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"									mininum(1)
//	@Param			sort    query	string	false	"Sort order. Default: desc"					Enums(asc, desc)
//	@Param			owner	query	string	false	"Owner celestia address"					minlength(47)	maxlength(47)
//	@Param			mailbox	query	string	false	"Mailbox hexademical identity"
//	@Param			type    query	string	false	"Comma-separated string of tokens type"		Enums(synthetic, collateral)
//	@Produce		json
//	@Success		200	{array}	responses.HyperlaneToken
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/token [get]
func (handler *HyperlaneHandler) ListTokens(c echo.Context) error {
	req, err := bindAndValidate[listHyperlaneTokenRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	filters, err := req.ToFilters(c.Request().Context(), handler.address, handler.mailbox)
	if err != nil {
		return badRequestError(c, err)
	}

	tokens, err := handler.tokens.List(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.HyperlaneToken, len(tokens))
	for i := range tokens {
		response[i] = responses.NewHyperlaneToken(tokens[i])
	}
	return returnArray(c, response)
}

type listHyperlaneTransferRequest struct {
	Limit   int         `query:"limit"   validate:"omitempty,min=1,max=100"`
	Offset  int         `query:"offset"  validate:"omitempty,min=0"`
	Sort    string      `query:"sort"    validate:"omitempty,oneof=asc desc"`
	Address string      `query:"address" validate:"omitempty,address"`
	Relayer string      `query:"relayer" validate:"omitempty,address"`
	Mailbox string      `query:"mailbox" validate:"omitempty,hexadecimal"`
	Token   string      `query:"token"   validate:"omitempty,hexadecimal"`
	Type    StringArray `query:"type"    validate:"omitempty,dive,hl_transfer_type"`
	Domain  uint64      `query:"domain"  validate:"omitempty,min=1"`
	Hash    string      `query:"hash"    validate:"omitempty,hexadecimal,len=64"`
}

func (req *listHyperlaneTransferRequest) ToFilters(ctx context.Context, address storage.IAddress, mailbox storage.IHLMailbox, tokens storage.IHLToken, tx storage.ITx) (storage.ListHyperlaneTransferFilters, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
	filters := storage.ListHyperlaneTransferFilters{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   pgSort(req.Sort),
		Domain: req.Domain,
	}

	if req.Mailbox != "" {
		id, err := hex.DecodeString(req.Mailbox)
		if err != nil {
			return filters, errors.Wrapf(err, "decoding mailbox id: %s", req.Mailbox)
		}
		mbx, err := mailbox.ByHash(ctx, id)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving mailbox by id: %x", id)
		}
		filters.MailboxId = mbx.Id
	}
	if req.Hash != "" {
		hash, err := hex.DecodeString(req.Hash)
		if err != nil {
			return filters, errors.Wrapf(err, "decoding tx hash: %s", req.Hash)
		}
		transaction, err := tx.ByHash(ctx, hash)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving tx by hash: %x", hash)
		}
		filters.TxId = transaction.Id
	}
	if req.Token != "" {
		id, err := hex.DecodeString(req.Token)
		if err != nil {
			return filters, errors.Wrapf(err, "decoding token id: %s", req.Token)
		}
		token, err := tokens.ByHash(ctx, id)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving token by id: %x", id)
		}
		filters.TokenId = token.Id
	}
	if req.Address != "" {
		_, hash, err := types.Address(req.Address).Decode()
		if err != nil {
			return filters, errors.Wrapf(err, "decoding address: %s", req.Address)
		}
		addr, err := address.ByHash(ctx, hash)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving address by hash: %x", hash)
		}
		filters.AddressId = addr.Id
	}
	if req.Relayer != "" {
		_, hash, err := types.Address(req.Relayer).Decode()
		if err != nil {
			return filters, errors.Wrapf(err, "decoding relayer address: %s", req.Relayer)
		}
		addr, err := address.ByHash(ctx, hash)
		if err != nil {
			return filters, errors.Wrapf(err, "receiving address by hash: %x", hash)
		}
		filters.RelayerId = addr.Id
	}
	if len(req.Type) > 0 {
		filters.Type = make([]storageTypes.HLTransferType, len(req.Type))
		for i := range req.Type {
			filters.Type[i] = storageTypes.HLTransferType(req.Type[i])
		}
	}

	return filters, nil
}

// ListTransfers godoc
//
//	@Summary		List hyperlane transfers info
//	@Description	List hyperlane transfers info
//	@Tags			hyperlane
//	@ID				list-hyperlane-transfers
//	@Param			limit	query	integer	false	"Count of requested entities"				mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"									mininum(1)
//	@Param			sort    query	string	false	"Sort order. Default: desc"					Enums(asc, desc)
//	@Param			address	query	string	false	"Celestia address"				         	minlength(47)	maxlength(47)
//	@Param			relayer	query	string	false	"Celestia address of relayer"				minlength(47)	maxlength(47)
//	@Param			mailbox	query	string	false	"Mailbox hexademical identity"
//	@Param			token	query	string	false	"Token hexademical identity"
//	@Param			type    query	string	false	"Comma-separated string of transfer type"	Enums(send, receive)
//	@Param			domain	query	integer	false	"Domain of counterparty chain"				mininum(1)
//	@Param			hash	query	string	false	"Transaction hash in hexadecimal"	minlength(64)	maxlength(64)
//	@Produce		json
//	@Success		200	{array}	responses.HyperlaneTransfer
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/transfer [get]
func (handler *HyperlaneHandler) ListTransfers(c echo.Context) error {
	req, err := bindAndValidate[listHyperlaneTransferRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	filters, err := req.ToFilters(c.Request().Context(), handler.address, handler.mailbox, handler.tokens, handler.txs)
	if err != nil {
		if handler.txs.IsNoRows(err) {
			return returnArray(c, []any{})
		}
		return badRequestError(c, err)
	}

	transfers, err := handler.transfers.List(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	response := make([]responses.HyperlaneTransfer, len(transfers))
	for i := range transfers {
		response[i] = responses.NewHyperlaneTransfer(transfers[i], handler.chainStore)
	}
	return returnArray(c, response)
}

// GetTransfer godoc
//
//	@Summary		Get transfer by id
//	@Description	Get transfer by id
//	@Tags			hyperlane
//	@ID				get-hyperlane-transfer
//	@Param			id	path	integer	true	"Internal identity"	mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.HyperlaneTransfer
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/hyperlane/transfer/{id} [get]
func (handler *HyperlaneHandler) GetTransfer(c echo.Context) error {
	req, err := bindAndValidate[getById](c)
	if err != nil {
		return badRequestError(c, err)
	}

	transfer, err := handler.transfers.ById(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.address)
	}

	return c.JSON(http.StatusOK, responses.NewHyperlaneTransfer(transfer, handler.chainStore))
}

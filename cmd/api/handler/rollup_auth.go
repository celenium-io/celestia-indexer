// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	enums "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type RollupAuthHandler struct {
	address   storage.IAddress
	namespace storage.INamespace
	rollups   storage.IRollup
	tx        sdk.Transactable
}

func NewRollupAuthHandler(
	rollups storage.IRollup,
	address storage.IAddress,
	namespace storage.INamespace,
	tx sdk.Transactable,
) RollupAuthHandler {
	return RollupAuthHandler{
		rollups:   rollups,
		address:   address,
		namespace: namespace,
		tx:        tx,
	}
}

type createRollupRequest struct {
	Name        string           `json:"name"        validate:"required,min=1"`
	Description string           `json:"description" validate:"required,min=1"`
	Website     string           `json:"website"     validate:"omitempty,url"`
	GitHub      string           `json:"github"      validate:"omitempty,url"`
	Twitter     string           `json:"twitter"     validate:"omitempty,url"`
	Logo        string           `json:"logo"        validate:"required,url"`
	L2Beat      string           `json:"l2_beat"     validate:"omitempty,url"`
	DeFiLama    string           `json:"defi_lama"   validate:"omitempty"`
	Bridge      string           `json:"bridge"      validate:"omitempty,eth_addr"`
	Explorer    string           `json:"explorer"    validate:"omitempty,url"`
	Stack       string           `json:"stack"       validate:"omitempty"`
	Links       []string         `json:"links"       validate:"omitempty,dive,url"`
	Category    string           `json:"category"    validate:"omitempty,category"`
	Tags        []string         `json:"tags"        validate:"omitempty"`
	Type        string           `json:"type"        validate:"omitempty,oneof=settled sovereign"`
	Compression string           `json:"compression" validate:"omitempty"`
	VM          string           `json:"vm"          validate:"omitempty"`
	Provider    string           `json:"provider"    validate:"omitempty"`
	SettledOn   string           `json:"settled_on"  validate:"omitempty"`
	Providers   []rollupProvider `json:"providers"   validate:"required,min=1,dive"`
}

type rollupProvider struct {
	Namespace string `json:"namespace" validate:"omitempty,base64,namespace"`
	Address   string `json:"address"   validate:"omitempty,address"`
}

func (handler RollupAuthHandler) Create(c echo.Context) error {
	val := c.Get(ApiKeyName)
	apiKey, ok := val.(storage.ApiKey)
	if !ok {
		return handleError(c, errInvalidApiKey, handler.address)
	}

	req, err := bindAndValidate[createRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	rollupId, err := handler.createRollup(c.Request().Context(), req, apiKey.Admin)
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": rollupId,
	})
}

func (handler RollupAuthHandler) createRollup(ctx context.Context, req *createRollupRequest, isAdmin bool) (uint64, error) {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return 0, err
	}

	rollup := storage.Rollup{
		Name:           req.Name,
		Description:    req.Description,
		Website:        req.Website,
		GitHub:         req.GitHub,
		Twitter:        req.Twitter,
		Logo:           req.Logo,
		L2Beat:         req.L2Beat,
		DeFiLama:       req.DeFiLama,
		Explorer:       req.Explorer,
		BridgeContract: req.Bridge,
		Stack:          req.Stack,
		Links:          req.Links,
		Compression:    req.Compression,
		Provider:       req.Provider,
		SettledOn:      req.SettledOn,
		VM:             req.VM,
		Type:           enums.RollupType(req.Type),
		Category:       enums.RollupCategory(req.Category),
		Slug:           slug.Make(req.Name),
		Tags:           req.Tags,
		Verified:       isAdmin,
	}

	if err := tx.SaveRollup(ctx, &rollup); err != nil {
		return 0, tx.HandleError(ctx, err)
	}

	providers, err := handler.createProviders(ctx, rollup.Id, req.Providers...)
	if err != nil {
		return 0, tx.HandleError(ctx, err)
	}

	if err := tx.SaveProviders(ctx, providers...); err != nil {
		return 0, tx.HandleError(ctx, err)
	}

	if err := tx.Flush(ctx); err != nil {
		return 0, err
	}
	return rollup.Id, nil
}

func (handler RollupAuthHandler) createProviders(ctx context.Context, rollupId uint64, data ...rollupProvider) ([]storage.RollupProvider, error) {
	providers := make([]storage.RollupProvider, len(data))
	for i := range data {
		providers[i].RollupId = rollupId
		if data[i].Address != "" {
			_, hashAddress, err := types.Address(data[i].Address).Decode()
			if err != nil {
				return nil, err
			}
			address, err := handler.address.ByHash(ctx, hashAddress)
			if err != nil {
				if handler.address.IsNoRows(err) {
					return nil, errors.Wrap(errUnknownAddress, data[i].Address)
				}
				return nil, err
			}
			providers[i].AddressId = address.Id
		}

		if data[i].Namespace != "" {
			hashNs, err := base64.StdEncoding.DecodeString(data[i].Namespace)
			if err != nil {
				return nil, err
			}
			ns, err := handler.namespace.ByNamespaceIdAndVersion(ctx, hashNs[1:], hashNs[0])
			if err != nil {
				if handler.namespace.IsNoRows(err) {
					return nil, errors.Wrap(errUnknownNamespace, data[i].Namespace)
				}
				return nil, err
			}
			providers[i].NamespaceId = ns.Id
		}
	}
	return providers, nil
}

type updateRollupRequest struct {
	Id          uint64           `param:"id"         validate:"required,min=1"`
	Name        string           `json:"name"        validate:"omitempty,min=1"`
	Description string           `json:"description" validate:"omitempty,min=1"`
	Website     string           `json:"website"     validate:"omitempty,url"`
	GitHub      string           `json:"github"      validate:"omitempty,url"`
	Twitter     string           `json:"twitter"     validate:"omitempty,url"`
	Logo        string           `json:"logo"        validate:"omitempty,url"`
	L2Beat      string           `json:"l2_beat"     validate:"omitempty,url"`
	DeFiLama    string           `json:"defi_lama"   validate:"omitempty"`
	Bridge      string           `json:"bridge"      validate:"omitempty,eth_addr"`
	Explorer    string           `json:"explorer"    validate:"omitempty,url"`
	Stack       string           `json:"stack"       validate:"omitempty"`
	Category    string           `json:"category"    validate:"omitempty,category"`
	Type        string           `json:"type"        validate:"omitempty,oneof=settled sovereign"`
	Compression string           `json:"compression" validate:"omitempty"`
	Provider    string           `json:"provider"    validate:"omitempty"`
	VM          string           `json:"vm"          validate:"omitempty"`
	SettledOn   string           `json:"settled_on"  validate:"omitempty"`
	Links       []string         `json:"links"       validate:"omitempty,dive,url"`
	Providers   []rollupProvider `json:"providers"   validate:"omitempty,min=1,dive"`
	Tags        []string         `json:"tags"        validate:"omitempty"`
}

func (handler RollupAuthHandler) Update(c echo.Context) error {
	val := c.Get(ApiKeyName)
	apiKey, ok := val.(storage.ApiKey)
	if !ok {
		return handleError(c, errInvalidApiKey, handler.address)
	}

	req, err := bindAndValidate[updateRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.updateRollup(c.Request().Context(), req, apiKey.Admin); err != nil {
		return handleError(c, err, handler.rollups)
	}

	return success(c)
}

func (handler RollupAuthHandler) updateRollup(ctx context.Context, req *updateRollupRequest, isAdmin bool) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	if _, err := handler.rollups.GetByID(ctx, req.Id); err != nil {
		return err
	}

	rollup := storage.Rollup{
		Id:             req.Id,
		Name:           req.Name,
		Slug:           slug.Make(req.Name),
		Description:    req.Description,
		Website:        req.Website,
		GitHub:         req.GitHub,
		Twitter:        req.Twitter,
		Logo:           req.Logo,
		L2Beat:         req.L2Beat,
		DeFiLama:       req.DeFiLama,
		Explorer:       req.Explorer,
		BridgeContract: req.Bridge,
		Stack:          req.Stack,
		Compression:    req.Compression,
		Provider:       req.Provider,
		SettledOn:      req.SettledOn,
		VM:             req.VM,
		Type:           enums.RollupType(req.Type),
		Category:       enums.RollupCategory(req.Category),
		Links:          req.Links,
		Tags:           req.Tags,
		Verified:       isAdmin,
	}

	if err := tx.UpdateRollup(ctx, &rollup); err != nil {
		return tx.HandleError(ctx, err)
	}

	if len(req.Providers) > 0 {
		if err := tx.DeleteProviders(ctx, req.Id); err != nil {
			return tx.HandleError(ctx, err)
		}

		providers, err := handler.createProviders(ctx, rollup.Id, req.Providers...)
		if err != nil {
			return tx.HandleError(ctx, err)
		}

		if err := tx.SaveProviders(ctx, providers...); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	return tx.Flush(ctx)
}

type deleteRollupRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

func (handler RollupAuthHandler) Delete(c echo.Context) error {
	req, err := bindAndValidate[deleteRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.deleteRollup(c.Request().Context(), req.Id); err != nil {
		return handleError(c, err, handler.rollups)
	}

	return success(c)
}

func (handler RollupAuthHandler) deleteRollup(ctx context.Context, id uint64) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	if err := tx.DeleteProviders(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.DeleteRollup(ctx, id); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

func (handler RollupAuthHandler) Unverified(c echo.Context) error {
	rollups, err := handler.rollups.Unverified(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.rollups)
	}

	response := make([]responses.Rollup, len(rollups))
	for i := range rollups {
		response[i] = responses.NewRollup(&rollups[i])
	}

	return returnArray(c, response)
}

type verifyRollupRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

func (handler RollupAuthHandler) Verify(c echo.Context) error {
	req, err := bindAndValidate[verifyRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.verify(c.Request().Context(), req.Id); err != nil {
		return handleError(c, err, handler.address)
	}

	return success(c)
}

func (handler RollupAuthHandler) verify(ctx context.Context, id uint64) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	err = tx.UpdateRollup(ctx, &storage.Rollup{
		Id:       id,
		Verified: true,
	})
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

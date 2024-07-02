// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"encoding/base64"

	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/postgres"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
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
	Logo        string           `json:"logo"        validate:"omitempty,url"`
	L2Beat      string           `json:"l2_beat"     validate:"omitempty,url"`
	Bridge      string           `json:"bridge"      validate:"omitempty,eth_addr"`
	Explorer    string           `json:"explorer"    validate:"omitempty,url"`
	Stack       string           `json:"stack"       validate:"omitempty"`
	Links       []string         `json:"links"       validate:"omitempty,dive,url"`
	Providers   []rollupProvider `json:"providers"   validate:"required,min=1"`
}

type rollupProvider struct {
	Namespace string `json:"namespace" validate:"omitempty,base64,namespace"`
	Address   string `json:"address"   validate:"required,address"`
}

func (handler RollupAuthHandler) Create(c echo.Context) error {
	req, err := bindAndValidate[createRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.createRollup(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.rollups)
	}

	return success(c)
}

func (handler RollupAuthHandler) createRollup(ctx context.Context, req *createRollupRequest) error {
	tx, err := postgres.BeginTransaction(ctx, handler.tx)
	if err != nil {
		return err
	}

	rollup := storage.Rollup{
		Name:           req.Name,
		Description:    req.Description,
		Website:        req.Website,
		GitHub:         req.GitHub,
		Twitter:        req.Twitter,
		Logo:           req.Logo,
		L2Beat:         req.L2Beat,
		Explorer:       req.Explorer,
		BridgeContract: req.Bridge,
		Stack:          req.Stack,
		Links:          req.Links,
		Slug:           slug.Make(req.Name),
	}

	if err := tx.SaveRollup(ctx, &rollup); err != nil {
		return tx.HandleError(ctx, err)
	}

	providers, err := handler.createProviders(ctx, rollup.Id, req.Providers...)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.SaveProviders(ctx, providers...); err != nil {
		return tx.HandleError(ctx, err)
	}

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}

	return tx.Flush(ctx)
}

func (handler RollupAuthHandler) createProviders(ctx context.Context, rollupId uint64, data ...rollupProvider) ([]storage.RollupProvider, error) {
	providers := make([]storage.RollupProvider, len(data))
	for i := range data {
		providers[i].RollupId = rollupId
		_, hashAddress, err := types.Address(data[i].Address).Decode()
		if err != nil {
			return nil, err
		}
		address, err := handler.address.ByHash(ctx, hashAddress)
		if err != nil {
			return nil, err
		}
		providers[i].AddressId = address.Id

		if data[i].Namespace != "" {
			hashNs, err := base64.StdEncoding.DecodeString(data[i].Namespace)
			if err != nil {
				return nil, err
			}
			ns, err := handler.namespace.ByNamespaceIdAndVersion(ctx, hashNs[1:], hashNs[0])
			if err != nil {
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
	Bridge      string           `json:"bridge"      validate:"omitempty,eth_addr"`
	Explorer    string           `json:"explorer"    validate:"omitempty,url"`
	Stack       string           `json:"stack"       validate:"omitempty"`
	Links       []string         `json:"links"       validate:"omitempty,dive,url"`
	Providers   []rollupProvider `json:"providers"   validate:"omitempty,min=1"`
}

func (handler RollupAuthHandler) Update(c echo.Context) error {
	req, err := bindAndValidate[updateRollupRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	if err := handler.updateRollup(c.Request().Context(), req); err != nil {
		return handleError(c, err, handler.rollups)
	}

	return success(c)
}

func (handler RollupAuthHandler) updateRollup(ctx context.Context, req *updateRollupRequest) error {
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
		Explorer:       req.Explorer,
		BridgeContract: req.Bridge,
		Stack:          req.Stack,
		Links:          req.Links,
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

	if err := tx.RefreshLeaderboard(ctx); err != nil {
		return tx.HandleError(ctx, err)
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

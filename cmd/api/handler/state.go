// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type StateHandler struct {
	state       storage.IState
	validator   storage.IValidator
	constants   storage.IConstant
	indexerName string
}

func NewStateHandler(state storage.IState, validator storage.IValidator, constants storage.IConstant, indexerName string) *StateHandler {
	return &StateHandler{
		state:       state,
		validator:   validator,
		constants:   constants,
		indexerName: indexerName,
	}
}

// Head godoc
//
//	@Summary		Get current indexer head
//	@Description	Get current indexer head
//	@Tags			general
//	@ID				head
//	@Produce		json
//	@Success		200	{object}	responses.State
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/head [get]
func (sh *StateHandler) Head(c echo.Context) error {
	state, err := sh.state.ByName(c.Request().Context(), sh.indexerName)
	if err != nil {
		return handleError(c, err, sh.state)
	}

	maxValidators, err := getMaxValidatorsCount(c.Request().Context(), sh.constants)
	if err != nil {
		return handleError(c, err, sh.state)
	}

	votingPower, err := sh.validator.TotalVotingPower(c.Request().Context(), maxValidators)
	if err != nil {
		return handleError(c, err, sh.state)
	}
	state.TotalVotingPower = votingPower

	return c.JSON(http.StatusOK, responses.NewState(state))
}

func getMaxValidatorsCount(ctx context.Context, constants storage.IConstant) (int, error) {
	maxValsConsts, err := constants.Get(ctx, types.ModuleNameStaking, "max_validators")
	if err != nil {
		return 0, errors.Wrap(err, "get max validators value")
	}
	return strconv.Atoi(maxValsConsts.Value)
}

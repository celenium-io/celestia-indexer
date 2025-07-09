// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type StateHandler struct {
	state       storage.IState
	validator   storage.IValidator
	indexerName string
}

func NewStateHandler(state storage.IState, validator storage.IValidator, indexerName string) *StateHandler {
	return &StateHandler{
		state:       state,
		validator:   validator,
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

	votingPower, err := sh.validator.TotalVotingPower(c.Request().Context())
	if err != nil {
		return handleError(c, err, sh.state)
	}
	state.TotalVotingPower = votingPower

	return c.JSON(http.StatusOK, responses.NewState(state))
}

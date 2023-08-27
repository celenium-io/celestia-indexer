package handler

import (
	"net/http"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/labstack/echo/v4"
)

type StateHandler struct {
	state storage.IState
}

func NewStateHandler(state storage.IState) *StateHandler {
	return &StateHandler{
		state: state,
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
//	@Router			/v1/head [get]
func (sh *StateHandler) Head(c echo.Context) error {
	state, err := sh.state.List(c.Request().Context(), 1, 0, sdk.SortOrderAsc)
	if err != nil {
		return internalServerError(c, err)
	}
	if len(state) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, responses.NewState(*state[0]))
}

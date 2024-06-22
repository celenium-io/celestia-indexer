// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/labstack/echo/v4"
)

type VestingHandler struct {
	vestings storage.IVestingPeriod
}

func NewVestingHandler(vestings storage.IVestingPeriod) *VestingHandler {
	return &VestingHandler{
		vestings: vestings,
	}
}

type getVestingPeriodsRequest struct {
	Id     uint64 `param:"id"     validate:"required,min=1"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

// Periods godoc
//
//	@Summary		Periods vesting periods by id
//	@Description	Periods vesting periods by id. Returns not empty array only for periodic vestings.
//	@Tags			vesting
//	@ID				get-vesting-periods
//	@Param			id		path	integer	true	"Internal identity"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.VestingPeriod
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/vesting/{id}/periods [get]
func (handler *VestingHandler) Periods(c echo.Context) error {
	req, err := bindAndValidate[getVestingPeriodsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	vestingPeriods, err := handler.vestings.ByVesting(c.Request().Context(), req.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.vestings)
	}

	response := make([]responses.VestingPeriod, len(vestingPeriods))
	for i := range vestingPeriods {
		response[i] = responses.NewVestingPeriod(vestingPeriods[i])
	}

	return returnArray(c, response)
}

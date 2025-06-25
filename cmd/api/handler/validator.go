// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"net/http"
	"strconv"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/math"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	st "github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/labstack/echo/v4"
)

type ValidatorHandler struct {
	validators      storage.IValidator
	blocks          storage.IBlock
	blockSignatures storage.IBlockSignature
	delegations     storage.IDelegation
	constants       storage.IConstant
	jails           storage.IJail
	votes           storage.IVote
	state           storage.IState
	indexerName     string
}

func NewValidatorHandler(
	validators storage.IValidator,
	blocks storage.IBlock,
	blockSignatures storage.IBlockSignature,
	delegations storage.IDelegation,
	constants storage.IConstant,
	jails storage.IJail,
	votes storage.IVote,
	state storage.IState,
	indexerName string,
) *ValidatorHandler {
	return &ValidatorHandler{
		validators:      validators,
		blocks:          blocks,
		blockSignatures: blockSignatures,
		delegations:     delegations,
		constants:       constants,
		jails:           jails,
		votes:           votes,
		state:           state,
		indexerName:     indexerName,
	}
}

type validatorRequest struct {
	Id uint64 `param:"id" validate:"required,min=1"`
}

// Get godoc
//
//	@Summary		Get validator info
//	@Description	Get validator info
//	@Tags			validator
//	@ID				get-validator
//	@Param			id	path	integer	true	"Internal validator id"
//	@Produce		json
//	@Success		200	{object}	responses.Validator
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id} [get]
func (handler *ValidatorHandler) Get(c echo.Context) error {
	req, err := bindAndValidate[validatorRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}

	validator, err := handler.validators.GetByID(c.Request().Context(), req.Id)
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	return c.JSON(http.StatusOK, responses.NewValidator(*validator))
}

type validatorsPagination struct {
	Limit  int   `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int   `query:"offset" validate:"omitempty,min=0"`
	Jailed *bool `query:"jailed" validate:"omitempty"`
}

func (req *validatorsPagination) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
}

// List godoc
//
//	@Summary		List validators
//	@Description	List validators
//	@Tags			validator
//	@ID				list-validator
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Param			jailed	query	boolean	false	"Return only jailed validators"
//	@Produce		json
//	@Success		200	{array}		responses.Validator
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators [get]
func (handler *ValidatorHandler) List(c echo.Context) error {
	req, err := bindAndValidate[validatorsPagination](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	validators, err := handler.validators.ListByPower(c.Request().Context(), storage.ValidatorFilters{
		Limit:  req.Limit,
		Offset: req.Offset,
		Jailed: req.Jailed,
	})
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	response := make([]responses.Validator, len(validators))
	for i := range validators {
		response[i] = *responses.NewValidator(validators[i])
	}

	return returnArray(c, response)
}

type validatorPageableRequest struct {
	Id     uint64 `param:"id"     validate:"required,min=1"`
	Limit  int    `query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset int    `query:"offset" validate:"omitempty,min=0"`
}

func (p *validatorPageableRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// Blocks godoc
//
//	@Summary		Get blocks which was proposed by validator
//	@Description	Get blocks which was proposed by validator
//	@Tags			validator
//	@ID				get-validator-blocks
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested entities"	mininum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						mininum(1)
//	@Produce		json
//	@Success		200	{object}	responses.Block
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id}/blocks [get]
func (handler *ValidatorHandler) Blocks(c echo.Context) error {
	req, err := bindAndValidate[validatorPageableRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	blocks, err := handler.blocks.ByProposer(c.Request().Context(), req.Id, req.Limit, req.Offset)
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	response := make([]responses.Block, len(blocks))
	for i := range blocks {
		response[i] = responses.NewBlock(blocks[i], true)
	}

	return returnArray(c, response)
}

type validatorUptimeRequest struct {
	Id    uint64      `param:"id"    validate:"required,min=1"`
	Limit types.Level `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *validatorUptimeRequest) SetDefault() {
	if r.Limit == 0 {
		r.Limit = 10
	}
}

// Uptime godoc
//
//	@Summary		Get validator's uptime and history of signed block
//	@Description	Get validator's uptime and history of signed block
//	@Tags			validator
//	@ID				get-validator-uptime
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested blocks"	mininum(1)	maximum(100)
//	@Produce		json
//	@Success		200	{object}	responses.ValidatorUptime
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id}/uptime [get]
func (handler *ValidatorHandler) Uptime(c echo.Context) error {
	req, err := bindAndValidate[validatorUptimeRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.blockSignatures)
	}

	startHeight := state.LastHeight - req.Limit - 1
	levels, err := handler.blockSignatures.LevelsByValidator(c.Request().Context(), req.Id, startHeight)
	if err != nil {
		return handleError(c, err, handler.blockSignatures)
	}

	response := responses.NewValidatorUptime(levels, state.LastHeight-1, req.Limit)
	return c.JSON(http.StatusOK, response)
}

type validatorDelegationsRequest struct {
	Id       uint64 `param:"id"        validate:"required,min=1"`
	Limit    int    `query:"limit"     validate:"omitempty,min=1,max=100"`
	Offset   int    `query:"offset"    validate:"omitempty,min=0"`
	ShowZero bool   `query:"show_zero" validate:"omitempty"`
}

func (p *validatorDelegationsRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// Delegators godoc
//
//	@Summary		Get validator's delegators
//	@Description	Get validator's delegators
//	@Tags			validator
//	@ID				validator-delegators
//	@Param			id			path	integer	true	"Internal validator id"
//	@Param			limit		query	integer	false	"Count of requested entities"	minimum(1)	maximum(100)
//	@Param			offset		query	integer	false	"Offset"						minimum(1)
//	@Param			show_zero	query	boolean	false	"Show zero delegations"
//	@Produce		json
//	@Success		200	{array}		responses.Delegation
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id}/delegators [get]
func (handler *ValidatorHandler) Delegators(c echo.Context) error {
	req, err := bindAndValidate[validatorDelegationsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	delegations, err := handler.delegations.ByValidator(
		c.Request().Context(),
		req.Id,
		req.Limit,
		req.Offset,
		req.ShowZero,
	)
	if err != nil {
		return handleError(c, err, handler.delegations)
	}

	response := make([]responses.Delegation, len(delegations))
	for i := range response {
		response[i] = responses.NewDelegation(delegations[i])
	}

	return returnArray(c, response)
}

// Jails godoc
//
//	@Summary		Get validator's jails
//	@Description	Get validator's jails
//	@Tags			validator
//	@ID				validator-jails
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Jail
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id}/jails [get]
func (handler *ValidatorHandler) Jails(c echo.Context) error {
	req, err := bindAndValidate[validatorPageableRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	jails, err := handler.jails.ByValidator(
		c.Request().Context(),
		req.Id,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return handleError(c, err, handler.delegations)
	}

	response := make([]responses.Jail, len(jails))
	for i := range response {
		response[i] = responses.NewJail(jails[i])
	}
	return returnArray(c, response)
}

// Count godoc
//
//	@Summary		Get validator's count by status
//	@Description	Get validator's count by status
//	@Tags			validator
//	@ID				validator-count
//	@Produce		json
//	@Success		200	{object}	responses.ValidatorCount
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/count [get]
func (handler *ValidatorHandler) Count(c echo.Context) error {
	state, err := handler.state.ByName(c.Request().Context(), handler.indexerName)
	if err != nil {
		return handleError(c, err, handler.state)
	}

	jailed, err := handler.validators.JailedCount(c.Request().Context())
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	constant, err := handler.constants.Get(c.Request().Context(), st.ModuleNameStaking, "max_validators")
	if err != nil {
		return handleError(c, err, handler.validators)
	}
	max, err := strconv.ParseInt(constant.Value, 10, 32)
	if err != nil {
		return handleError(c, err, handler.validators)
	}

	active := math.Min(int(max), state.TotalValidators-jailed)
	return c.JSON(http.StatusOK, responses.ValidatorCount{
		Total:    state.TotalValidators,
		Jailed:   jailed,
		Active:   active,
		Inactive: state.TotalValidators - jailed - active,
	})
}

// Votes godoc
//
//	@Summary		Get list of votes for validator
//	@Description	Get list of votes for validator
//	@Tags			validator
//	@ID				validator-votes
//	@Param			id		path	integer	true	"Internal validator id"
//	@Param			limit	query	integer	false	"Count of requested entities"	minimum(1)	maximum(100)
//	@Param			offset	query	integer	false	"Offset"						minimum(1)
//	@Produce		json
//	@Success		200	{array}		responses.Vote
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/validators/{id}/votes [get]
func (handler *ValidatorHandler) Votes(c echo.Context) error {
	req, err := bindAndValidate[validatorPageableRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	jails, err := handler.votes.ByValidatorId(
		c.Request().Context(),
		req.Id,
		storage.VoteFilters{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	)
	if err != nil {
		return handleError(c, err, handler.delegations)
	}

	response := make([]responses.Vote, len(jails))
	for i := range response {
		response[i] = responses.NewVote(jails[i])
	}
	return returnArray(c, response)
}

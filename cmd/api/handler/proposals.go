// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/hex"

	"github.com/celenium-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/celenium-io/celestia-indexer/internal/storage"
	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
)

// ProposalsHandler -
type ProposalsHandler struct {
	proposals storage.IProposal
	votes     storage.IVote
	address   storage.IAddress
}

func NewProposalsHandler(
	proposals storage.IProposal,
	votes storage.IVote,
	address storage.IAddress,
) ProposalsHandler {
	return ProposalsHandler{
		proposals: proposals,
		votes:     votes,
		address:   address,
	}
}

type listProposalsRequest struct {
	Limit    int         `query:"limit"    validate:"omitempty,min=1,max=100"`
	Offset   int         `query:"offset"   validate:"omitempty,min=0"`
	Sort     string      `query:"sort"     validate:"omitempty,oneof=asc desc"`
	Status   StringArray `query:"status"   validate:"omitempty,dive,proposal_status"`
	Type     StringArray `query:"type"     validate:"omitempty,dive,proposal_type"`
	Proposer string      `query:"proposer" validate:"omitempty,address"`
}

func (req *listProposalsRequest) SetDefault() {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Sort == "" {
		req.Sort = desc
	}
}

func (req listProposalsRequest) toFilters(proposerId uint64) storage.ListProposalFilters {
	filters := storage.ListProposalFilters{
		Limit:      req.Limit,
		Offset:     req.Offset,
		ProposerId: proposerId,
		Sort:       pgSort(req.Sort),
		Status:     make([]types.ProposalStatus, len(req.Status)),
		Type:       make([]types.ProposalType, len(req.Type)),
	}

	for i := range req.Status {
		filters.Status[i] = types.ProposalStatus(req.Status[i])
	}
	for i := range req.Type {
		filters.Type[i] = types.ProposalType(req.Type[i])
	}

	return filters
}

// List godoc
//
//		@Summary		List proposal info
//		@Description	List proposal info
//		@Tags			proposal
//		@ID				list-proposal
//		@Param			limit	    query	integer	false	"Count of requested entities"					mininum(1)	maximum(100)
//		@Param			offset	    query	integer	false	"Offset"										mininum(1)
//		@Param			sort	    query	string	false	"Sort order. Default: desc"						Enums(asc, desc)
//	    @Param			proposer	query	string	false	"Proposer celestia address"	                    minlength(47)	maxlength(47)
//		@Param          status      query   string  false   "Comma-separated proposal status list"
//		@Param          type        query   string  false   "Comma-separated proposal type list"
//		@Produce		json
//		@Success		200	{array}		responses.Proposal
//		@Failure		400	{object}	Error
//		@Failure		500	{object}	Error
//		@Router			/proposal [get]
func (handler *ProposalsHandler) List(c echo.Context) error {
	req, err := bindAndValidate[listProposalsRequest](c)
	if err != nil {
		return badRequestError(c, err)
	}
	req.SetDefault()

	var proposerId uint64
	if req.Proposer != "" {
		proposerHash, err := hex.DecodeString(req.Proposer)
		if err != nil {
			return badRequestError(c, err)
		}
		proposerIds, err := handler.address.IdByHash(c.Request().Context(), proposerHash)
		if err != nil {
			return handleError(c, err, handler.address)
		}
		if len(proposerIds) == 1 {
			proposerId = proposerIds[0]
		}
	}

	filters := req.toFilters(proposerId)

	proposals, err := handler.proposals.ListWithFilters(c.Request().Context(), filters)
	if err != nil {
		return handleError(c, err, handler.proposals)
	}
	response := make([]responses.Proposal, len(proposals))
	for i := range proposals {
		response[i] = responses.NewProposal(proposals[i])
	}
	return returnArray(c, response)
}

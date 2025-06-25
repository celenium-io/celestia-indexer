// SPDX-FileCopyrightText: 2025 PK Lab AG <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/celenium-io/celestia-indexer/pkg/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CelestiaApiValidator struct {
	validator *validator.Validate
}

func NewCelestiaApiValidator() *CelestiaApiValidator {
	v := validator.New()
	if err := v.RegisterValidation("address", addressValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("status", statusValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("msg_type", msgTypeValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("namespace", namespaceValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("category", categoryValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("type", typeValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("proposal_status", proposalStatusValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("proposal_type", proposalTypeValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("vote_option", voteOptionValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("voter_type", voterTypeValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("ibc_channel_status", ibcChannelStatusValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("hl_token_type", hyperlaneTokenTypeValidator()); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("hl_transfer_type", hyperlaneTransferTypeValidator()); err != nil {
		panic(err)
	}
	return &CelestiaApiValidator{validator: v}
}

func (v *CelestiaApiValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func isAddress(address string) bool {
	return validateAddress(address, pkgTypes.AddressPrefixCelestia, 47, 128)
}

func isValoperAddress(address string) bool {
	return validateAddress(address, pkgTypes.AddressPrefixValoper, 54)
}

func validateAddress(address string, wantPrefix string, length ...int) bool {
	addrLen := len(address)

	switch len(length) {
	case 1:
		if addrLen != length[0] {
			return false
		}
	case 2:
		if addrLen < length[0] || addrLen > length[1] {
			return false
		}
	default:
		return false
	}

	prefix, _, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return false
	}

	return prefix == wantPrefix
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return isAddress(fl.Field().String()) || isValoperAddress(fl.Field().String())
	}
}

func statusValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseStatus(fl.Field().String())
		return err == nil
	}
}

func msgTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseMsgType(fl.Field().String())
		return err == nil
	}
}

func isNamespace(s string) bool {
	hash, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}
	if len(hash) != 29 {
		return false
	}
	return true
}

func namespaceValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return isNamespace(fl.Field().String())
	}
}

func categoryValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseRollupCategory(fl.Field().String())
		return err == nil
	}
}

func typeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseRollupType(fl.Field().String())
		return err == nil
	}
}

func proposalStatusValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseProposalStatus(fl.Field().String())
		return err == nil
	}
}

func proposalTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseProposalType(fl.Field().String())
		return err == nil
	}
}

func voteOptionValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseVoteOption(fl.Field().String())
		return err == nil
	}
}

func voterTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseVoterType(fl.Field().String())
		return err == nil
	}
}

func ibcChannelStatusValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseIbcChannelStatus(fl.Field().String())
		return err == nil
	}
}

func hyperlaneTokenTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseHLTokenType(fl.Field().String())
		return err == nil
	}
}

func hyperlaneTransferTypeValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		_, err := types.ParseHLTransferType(fl.Field().String())
		return err == nil
	}
}

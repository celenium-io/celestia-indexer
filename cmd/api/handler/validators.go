// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/celenium-io/celestia-indexer/internal/storage"
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
	return &CelestiaApiValidator{validator: v}
}

func (v *CelestiaApiValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func isAddress(address string) bool {
	return validateAddress(address, pkgTypes.AddressPrefixCelestia, 47)
}

func isValoperAddress(address string) bool {
	return validateAddress(address, pkgTypes.AddressPrefixValoper, 54)
}

func validateAddress(address string, wantPrefix string, length int) bool {
	if len(address) != length {
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

type KeyValidator struct {
	apiKeys    storage.IApiKey
	errChecker NoRows
}

func NewKeyValidator(apiKeys storage.IApiKey, errChecker NoRows) KeyValidator {
	return KeyValidator{apiKeys: apiKeys, errChecker: errChecker}
}

const ApiKeyName = "api_key"

func (kv KeyValidator) Validate(key string, c echo.Context) (bool, error) {
	apiKey, err := kv.apiKeys.Get(c.Request().Context(), key)
	if err != nil {
		if kv.errChecker.IsNoRows(err) {
			return false, nil
		}
		return false, err
	}
	c.Logger().Infof("using apikey: %s", apiKey.Description)
	c.Set(ApiKeyName, apiKey)
	return true, nil
}

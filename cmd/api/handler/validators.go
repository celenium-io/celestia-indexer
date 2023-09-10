package handler

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	pkgTypes "github.com/dipdup-io/celestia-indexer/pkg/types"
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
	return &CelestiaApiValidator{validator: v}
}

func (v *CelestiaApiValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func isAddress(address string) bool {
	switch len(address) {
	case 47:
		return validateAddress(address, pkgTypes.AddressPrefixCelestia)
	case 54:
		return validateAddress(address, pkgTypes.AddressPrefixValoper)
	default:
		return false
	}
}

func validateAddress(address string, wantPrefix string) bool {
	prefix, _, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return false
	}

	return prefix == wantPrefix
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return isAddress(fl.Field().String())
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

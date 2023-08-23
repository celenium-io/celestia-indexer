package handler

import (
	"net/http"

	"github.com/btcsuite/btcutil/bech32"
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
	return &CelestiaApiValidator{validator: v}
}

func (v *CelestiaApiValidator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func addressValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		address := fl.Field().String()
		if len(address) != 47 {
			return false
		}
		prefix, _, err := bech32.Decode(address)
		if err != nil {
			return false
		}

		return prefix == "celestia"
	}
}

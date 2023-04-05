package util

import (
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
	once     = new(sync.Once)
)

func registerCustomValidation() {
	validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		combine := strings.Join(Currencies, "")
		return strings.Contains(combine, fl.Field().String())
	})
}

type ErrorResponse struct {
	FailedField string `json:"failedfield"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(i interface{}) []*ErrorResponse {
	once.Do(registerCustomValidation)

	var errors []*ErrorResponse
	err := validate.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

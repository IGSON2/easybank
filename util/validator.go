package util

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type ErrorResponse struct {
	FailedField string `json:"failedfield"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(i interface{}) []*ErrorResponse {
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

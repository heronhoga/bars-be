package utils

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidateStruct(req interface{}) map[string]string {
	err := Validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors[err.Field()] = "This field is required"
		case "min":
			errors[err.Field()] = "Must be at least " + err.Param() + " characters"
		default:
			errors[err.Field()] = "Invalid value"
		}
	}
	return errors
}

package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s interface{}) error
}

type ValidatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &ValidatorImpl{
		validate: validator.New(),
	}
}

func (v *ValidatorImpl) ValidateStruct(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errors := ""
	for _, fieldErr := range validationErrors {
		switch fieldErr.Tag() {
		case "required":
			errors += fmt.Sprintf("%s is required", fieldErr.Field())
		case "email":
			errors += fmt.Sprintf("%s must be a valid email", fieldErr.Field())
		case "min":
			errors += fmt.Sprintf("%s must be at least %s characters", fieldErr.Field(), fieldErr.Param())
		case "max":
			errors += fmt.Sprintf("%s must be at most %s characters", fieldErr.Field(), fieldErr.Param())
		}
	}

	return fmt.Errorf("%v", errors)
}

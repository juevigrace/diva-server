package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	validator *validator.Validate
}

func NewXValidator() *XValidator {
	return &XValidator{
		validator: validator.New(),
	}
}

func (v *XValidator) Validate(body any) error {
	errs := v.validator.Struct(body)
	if errs != nil {
		errMsgs := make([]string, 0)
		for _, err := range errs.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			))
		}

		return errors.New(strings.Join(errMsgs, " and "))
	}

	return nil
}

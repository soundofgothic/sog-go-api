package rjson

import (
	"github.com/enhanced-tools/errors"
	"github.com/go-playground/validator/v10"
)

type Validation struct {
	Errors validator.ValidationErrors
}

type SingleValidationError struct {
	Key    string `json:"key"`
	Reason string `json:"reason"`
}

func NewValidationError(errors validator.ValidationErrors) errors.ErrorOpt {
	return Validation{
		Errors: errors,
	}
}

func (v Validation) MapFormatter() map[string]interface{} {
	errs := make([]SingleValidationError, 0, len(v.Errors))
	for _, e := range v.Errors {
		errs = append(errs, SingleValidationError{
			Key:    e.Field(),
			Reason: e.Tag(),
		})
	}
	return map[string]interface{}{
		"validation": errs,
	}
}

func (v Validation) Type() errors.ErrorOptType {
	return errors.ErrorOptType("validation")
}

func (v Validation) Verbosity() int {
	return 0
}

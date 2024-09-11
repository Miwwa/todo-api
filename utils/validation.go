package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"net/http"
)

type (
	AppValidator struct {
		validator *validator.Validate
	}

	ValidationError struct {
		FailedField string
		Tag         string
		Value       interface{}
	}

	ValidationErrorResponse struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Details []ValidationError `json:"details"`
	}
)

func (v ValidationErrorResponse) Error() string {
	return fmt.Sprintf("%s:%v", v.Message, v.Details)
}

func NewValidator() *AppValidator {
	v := validator.New()
	err := v.RegisterValidation("ulid", validateUlid)
	if err != nil {
		return nil
	}

	return &AppValidator{validator: v}
}

func validateUlid(fl validator.FieldLevel) bool {
	_, err := ulid.Parse(fl.Field().String())
	return err == nil
}

func (v AppValidator) Validate(data interface{}) error {
	var validationErrors []ValidationError

	errs := v.validator.Struct(data)
	if errs == nil {
		return nil
	}

	for _, err := range errs.(validator.ValidationErrors) {
		var elem ValidationError

		elem.FailedField = err.Field()
		elem.Tag = err.Tag()
		elem.Value = err.Value()

		validationErrors = append(validationErrors, elem)
	}

	return ValidationErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "request data validation failed",
		Details: validationErrors,
	}
}

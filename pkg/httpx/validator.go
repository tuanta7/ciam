package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func validateErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be a maximum of %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation on '%s'", e.Field(), e.Tag())
	}
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			var messages []string
			for _, e := range validateErrs {
				messages = append(messages, validateErrorMessage(e))
			}

			return Error{
				Code:    http.StatusBadRequest,
				Message: strings.Join(messages, "; "),
			}
		}
	}

	return nil
}

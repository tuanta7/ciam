package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type JSON map[string]any

func ReadJSON(payload io.Reader, data any) error {
	decoder := json.NewDecoder(payload)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// ParseJSON reads input data to a struct and validates
func ParseJSON(payload io.Reader, data any) error {
	err := ReadJSON(payload, data)
	if err != nil {
		return err
	}

	return validateStruct(data)
}

func WriteJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	return err
}

func ErrorJSON(w http.ResponseWriter, err HTTPError) error {
	code := http.StatusInternalServerError
	if err.Code > 0 {
		code = err.Code
	}

	return WriteJSON(w, code, err)
}

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

func validateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			var messages []string
			for _, e := range validateErrs {
				messages = append(messages, validateErrorMessage(e))
			}

			return errors.New(strings.Join(messages, "; "))
		}
	}

	return nil
}

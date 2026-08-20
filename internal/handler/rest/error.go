package rest

import "net/http"

type HTTPError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Hint        string `json:"hint,omitempty"`
}

func Error(code int, message string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface.
func (he HTTPError) Error() string {
	return he.Message
}

func (he HTTPError) WithHint(hint string) HTTPError {
	he.Hint = hint
	return he
}

func (he HTTPError) WithDescription(description string) HTTPError {
	he.Description = description
	return he
}

func (he HTTPError) WithMessage(message string) HTTPError {
	he.Message = message
	return he
}

func InvalidArgumentError() HTTPError {
	return Error(http.StatusBadRequest, "invalid argument")
}

func InternalError() HTTPError {
	return Error(http.StatusInternalServerError, "internal error").
		WithDescription("something went wrong")
}

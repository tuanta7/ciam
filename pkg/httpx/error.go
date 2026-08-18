package httpx

import "net/http"

type HTTPError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Hint        string `json:"hint,omitempty"`
}

func Error(code int, message string, opts ...Option) HTTPError {
	err := HTTPError{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		opt(&err)
	}

	return err
}

type Option func(e *HTTPError)

func WithHint(hint string) Option {
	return func(e *HTTPError) {
		e.Hint = hint
	}
}

func WithDescription(description string) Option {
	return func(e *HTTPError) {
		e.Description = description
	}
}

func WithMessage(message string) Option {
	return func(e *HTTPError) {
		e.Message = message
	}
}

// Error implements the error interface.
func (he HTTPError) Error() string {
	return he.Message
}

func InvalidArgumentError(opts ...Option) HTTPError {
	return Error(http.StatusBadRequest, "invalid argument", opts...)
}

func InternalError(opts ...Option) HTTPError {
	err := HTTPError{
		Code:        http.StatusInternalServerError,
		Message:     "internal error",
		Description: "something went wrong",
	}

	for _, opt := range opts {
		opt(&err)
	}

	return err
}

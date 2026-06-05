package httpx

import "net/http"

type Error struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Hint        string `json:"hint,omitempty"`
}

func NewError(code int, message string, opts ...Option) Error {
	err := Error{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		opt(&err)
	}

	return err
}

type Option func(e *Error)

func WithHint(hint string) Option {
	return func(e *Error) {
		e.Hint = hint
	}
}

func WithDescription(description string) Option {
	return func(e *Error) {
		e.Description = description
	}
}

func WithMessage(message string) Option {
	return func(e *Error) {
		e.Message = message
	}
}

// Error implements the error interface.
func (he Error) Error() string {
	return he.Message
}

func NewInvalidArgumentError(opts ...Option) Error {
	return NewError(http.StatusBadRequest, "invalid argument", opts...)
}

func NewInternalError(opts ...Option) Error {
	err := Error{
		Code:        http.StatusInternalServerError,
		Message:     "internal error",
		Description: "something went wrong",
	}

	for _, opt := range opts {
		opt(&err)
	}

	return err
}

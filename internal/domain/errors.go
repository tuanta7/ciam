package domain

import "errors"

var (
	ErrClientNotFound = errors.New("client not found")
	ErrInvalidClient  = errors.New("invalid client")

	ErrAuthRequestNotFound = errors.New("auth request not found")
)

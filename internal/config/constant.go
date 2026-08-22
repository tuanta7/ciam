package config

import "time"

const (
	DefaultIDTokenLifetime = time.Hour
)

const (
	DefaultLoginURLTemplate = "/login?auth_request_id=%s"
)

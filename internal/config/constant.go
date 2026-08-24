package config

import "time"

const (
	DefaultIDTokenLifetime = time.Hour
)

const (
	DefaultLoginURLTemplate = "/login?auth_request_id=%s"

	// DefaultAuthRequestLifetime bounds how long a user has to finish logging in.
	// Expired requests are ignored on read; a reaper still has to delete the rows.
	DefaultAuthRequestLifetime = 15 * time.Minute
)

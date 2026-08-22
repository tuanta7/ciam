package utils

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func NowUTCString() string {
	return NowUTC().String()
}

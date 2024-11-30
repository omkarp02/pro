package types

import "time"

const (
	ACCESS_TOKEN_EXPIRY = time.Hour
	ACCESS_TOKEN        = "access"
	REFRESH_TOKEN       = "refresh"
)

var REFRESH_TOKEN_EXPIRY = time.Now().Add(24 * time.Hour)

package constant

import "time"

const (
	ACCESS_TOKEN_EXPIRY  = time.Hour
	REFRESH_TOKEN_EXPIRY = time.Hour
	ACCESS_TOKEN         = "access"
	REFRESH_TOKEN        = "refresh"
	REFRESH_TOKEN_COOKIE = "refreshToken"
)

var REFRESH_TOKEN_COOKIE_EXPIRY = time.Now().Add(24 * time.Hour)

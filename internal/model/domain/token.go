package domain

import "time"

type Tokens struct {
	Access  string
	Refresh string
}

type RefreshTokens struct {
	Refresh   string
	ExpiredAt time.Time
}

package models

import "time"

type JwtGenerator struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

package token

import "time"

type Token struct {
	AuthId         string    `json:"auth_id"`
	ExpirationTime time.Time `json:"expiration_time"`
}

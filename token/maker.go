package token

import "time"

type Maker interface {
	CreateToken(id uint, username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

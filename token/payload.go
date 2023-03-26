package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken     = fmt.Errorf("token is expired")
	ErrInvalidToken     = fmt.Errorf("token is invalid")
	ErrorInvalidKeySize = fmt.Errorf("invalid key size : must be at least %d characters", minSecretKeySize)
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

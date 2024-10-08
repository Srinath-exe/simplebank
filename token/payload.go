package token

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("signature is invalid")
)

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	expiry := now.Add(duration)
	fmt.Println("Expiry time: ", expiry)
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: expiry,
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

package tokens

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

type PayloadData struct {
	Data string `json:"data,omitempty"`
}

type Payload struct {
	Id        uuid.UUID   `json:"id,omitempty"`
	IssuedAt  time.Time   `json:"iat,omitempty"`
	ExpiredAt time.Time   `json:"expires,omitempty"`
	Payload   PayloadData `json:"payload,omitempty"`
}

func NewPayload(p PayloadData, duration time.Duration) (*Payload, error) {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Id:        token_id,
		Payload:   p,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

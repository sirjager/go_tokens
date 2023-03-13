package tokens

import "time"

const (
	Jwt    = "jwt"
	Paseto = "paseto"
)

// TokenBuilder is an interface for managing tokens
type TokenBuilder interface {
	// Create Token if token for specific duration
	CreateToken(data PayloadData, duration time.Duration) (string, *Payload, error)

	// Verify Token if token is valid or not
	VerifyToken(token string) (*Payload, error)
}

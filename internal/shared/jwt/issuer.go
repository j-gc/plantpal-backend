package jwt

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// HS256Issuer issues JWT tokens using HS256 algorithm
type HS256Issuer struct {
	secret []byte
	issuer string
}

// NewHS256Issuer creates a new JWT issuer
func NewHS256Issuer(secret, issuer string) *HS256Issuer {
	return &HS256Issuer{
		secret: []byte(secret),
		issuer: issuer,
	}
}

// Issue creates a new JWT token
func (i *HS256Issuer) Issue(subject string, ttl time.Duration, claims map[string]any) (string, error) {
	now := time.Now()

	tokenClaims := jwtv5.MapClaims{
		"iss": i.issuer,
		"sub": subject,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	// Add custom claims
	for k, v := range claims {
		tokenClaims[k] = v
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, tokenClaims)
	return token.SignedString(i.secret)
}

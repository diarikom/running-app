package njwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// Errors
var ErrInvalidClaim = errors.New("njwt: invalid token claims")
var ErrInvalidSignMethod = errors.New("njwt: unexpected signing method")

type Claim struct {
	Extra   map[string]string `json:"ext,omitempty"`
	Session string            `json:"ses,omitempty"`
	Purpose int               `json:"pur,omitempty"`
	jwt.StandardClaims
}

type Token struct {
	Encoded   string
	CreatedAt int64
	ExpiredAt int64
}

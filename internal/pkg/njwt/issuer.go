package njwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Issuer struct {
	Key             []byte
	DefaultLifetime time.Duration
	Issuer          string
}

type ClaimOpt struct {
	SessionId string
	Subject   string
	Audience  string
	Lifetime  time.Duration
	Purpose   int
	Extras    map[string]string
}

func (j *Issuer) New(opt ClaimOpt) (*Token, error) {
	// If lifetime is zero, set lifetime to default
	if opt.Lifetime == 0 {
		opt.Lifetime = j.DefaultLifetime
	} else {
		opt.Lifetime *= time.Minute
	}

	// Initiate current timestamp and expire time
	t := time.Now().UTC()
	createdAt := t.Unix()
	expiredAt := t.Add(opt.Lifetime).Unix()

	// Initiate access token claims
	claims := &Claim{
		Extra:   opt.Extras,
		Session: opt.SessionId,
		Purpose: opt.Purpose,
		StandardClaims: jwt.StandardClaims{
			Subject:   opt.Subject,
			Issuer:    j.Issuer,
			Audience:  opt.Audience,
			IssuedAt:  createdAt,
			ExpiresAt: expiredAt,
		},
	}

	// Added Claims to Token
	payload := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := payload.SignedString(j.Key)
	if err != nil {
		return nil, err
	}

	// Generate token
	result := Token{
		Encoded:   token,
		CreatedAt: createdAt,
		ExpiredAt: expiredAt,
	}
	return &result, nil
}

func (j *Issuer) Verify(input string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(input, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return j.Key, nil
	})

	// Check parsing err
	if err != nil {
		return nil, err
	}

	// Check claim
	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaim
	}
	return claims, nil
}

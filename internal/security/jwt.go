package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JsonWebTokenSecurity struct {
	Key []byte
}

type JsonWebTokenClaims struct {
	Sub   string
	Scope []string
}

type Claims struct {
	Subject string `json:"sub"`
	Scope   []string
	jwt.RegisteredClaims
}

// Create generates a JSON Web Token (JWT) using the provided claims and expiration time.
// It returns the signed token as a string and an error if any.
func (j *JsonWebTokenSecurity) Create(c *JsonWebTokenClaims, exp int) (string, error) {
	claims := &Claims{
		Subject: c.Sub,
		Scope:   c.Scope,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exp) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}

// Verify verifies the validity of a JSON Web Token (JWT) and returns the ID of the token's subject.
// If the token is invalid or expired, an error is returned.
func (j *JsonWebTokenSecurity) Verify(token string) (c *Claims, err error) {
	var claims Claims

	dec, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.Key, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.New("your token has expired")
		}

		return nil, errors.New("your token has expired, login again")
	}

	if !dec.Valid {
		return nil, errors.New("your token is invalid")
	}

	return &claims, nil
}

package goauth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// JwtMiddleware represents the basic fields needed to performe Jwt authentication.
// HasError is set to true if there's an error parsing the token.
type JwtMiddleware struct {
	Token     *jwt.Token     `json:"-"`
	HasError  bool           `json:"-"`
	GetSecret jwt.Keyfunc    `json:"-"`
	GetClaims GenerateClaims `json:"-"`
}

type GenerateSignedJWT func(r *http.Request) (string, error)

// CreateSecurityContext parses the request headers to get the bearer token.
func (dbtc *JwtMiddleware) CreateSecurityContext(r *http.Request) {
	auth := r.Header.Get("Authorization")
	var err error

	if strings.Contains(auth, "Bearer ") {
		dbtc.Token, err = jwt.Parse(extractTokenFromHeader(auth), dbtc.GetSecret)
		dbtc.HasError = err != nil
	} else {
		dbtc.HasError = true // no token found
	}
}

// HasPermission checks if a token has been parsed, if the token is valid and
// if there was an error during the parsing process
func (dbtc *JwtMiddleware) HasPermission() bool {
	return dbtc.Token != nil && dbtc.Token.Valid && !dbtc.HasError
}

func extractTokenFromHeader(header string) string {
	const bearerPrefix = "Bearer "
	if header != "" && strings.HasPrefix(header, bearerPrefix) {
		return header[len(bearerPrefix):]
	}
	return ""
}

// NewClaimBase returns a basic jwt.RegisteredClaims struct that can be inserted into your own custom claims if need be.
func NewClaimBase(expiredAt *jwt.NumericDate, issuer, subject string, id string, audience []string) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: expiredAt,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
		Subject:   subject,
		ID:        id,
		Audience:  audience,
	}
}

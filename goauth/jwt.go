package goauth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type JwtTokenAuthController struct {
	Token     *jwt.Token  `json:"-"`
	Error     bool        `json:"-"`
	GetSecret jwt.Keyfunc `json:"-"`
	GetClaims LoginUser   `json:"-"`
}

type JwtClaimsEndpointSecurity struct {
	Permissions map[string]string `json:"-"`
	jwt.RegisteredClaims
}

func (dbtc *JwtTokenAuthController) CreateSecurityContext(r *http.Request) {
	auth := r.Header.Get("Authorization")
	var err error

	if strings.Contains(auth, "Bearer ") {
		dbtc.Token, err = jwt.Parse(extractTokenFromHeader(auth), dbtc.GetSecret)
		dbtc.Error = err != nil
	}
}

func (dbtc *JwtTokenAuthController) HasPermission() bool {
	return dbtc.Token != nil && dbtc.Token.Valid && !dbtc.Error
}

func extractTokenFromHeader(header string) string {
	const bearerPrefix = "Bearer "
	if header != "" && strings.HasPrefix(header, bearerPrefix) {
		return header[len(bearerPrefix):]
	}
	return ""
}

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

package goauth

import "net/http"

// ValidateToken checks if the token is valid
type ValidateToken func(token string) bool

// LoginUser takes a username and a password and returns a token
type LoginUser func(username string, password string) string

type DefaultBearerTokenController struct {
	Token               string        `json:"-"`
	ValidateTokenMethod ValidateToken `json:"-"`
	LoginUserMethod     LoginUser     `json:"-"`
}

func (dbtc DefaultBearerTokenController) CreateSecurityContext(r *http.Request) {
	bearerToken := r.Header.Get("Bearer")
	username := r.Header.Get("username")
	password := r.Header.Get("password")

	if username != "" && password != "" {
		dbtc.Token = dbtc.LoginUserMethod(username, password)
	} else if bearerToken != "" {
		dbtc.Token = bearerToken
	}
}

func (dbtc DefaultBearerTokenController) HasPermission() bool {
	return dbtc.Token != "" && dbtc.ValidateTokenMethod(dbtc.Token)
}

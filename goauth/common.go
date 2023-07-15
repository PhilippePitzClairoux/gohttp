package goauth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"regexp"
)

// ExtractUsernamePassword extract username and password value for Authorization
var ExtractUsernamePassword = regexp.MustCompile("username=([\\w\\-]+),password=([\\w\\-]+)")

// ValidateToken checks if the token is valid
type ValidateToken func(token *jwt.Token) bool

// LoginUser confirm user exists and return its claims
type LoginUser func(r *http.Request) bool

type GenerateClaims func() jwt.Claims

// ValidateUser check's if a username and password are valid
type ValidateUser func(data ...any) bool

// ExtractSecurityContext is used by the middleware layer to extract username, password, token, etc to perform auth
type ExtractSecurityContext func(req *http.Request) any

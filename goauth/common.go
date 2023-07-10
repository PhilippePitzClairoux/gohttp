package goauth

import (
	"github.com/golang-jwt/jwt/v5"
	"regexp"
)

// ValidateToken checks if the token is valid
type ValidateToken func(token *jwt.Token) bool

// LoginUser confirm user exists and return its claims
type LoginUser func(username string, password string) jwt.RegisteredClaims

// ExtractUsernamePassword extract username and password value for Authorization
var ExtractUsernamePassword = regexp.MustCompile("username=([\\w\\-]+),password=([\\w\\-]+)")

// ValidateUser check's if a username and password are valid
type ValidateUser func(username string, password string) bool

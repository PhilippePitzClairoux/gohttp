package goauth

import (
	"net/http"
	"regexp"
)

var extractUsernamePassword = regexp.MustCompile("username=([\\w\\-]+),password=([\\w\\-]+)")

type ValidateUser func(username string, password string) bool

type HttpBasicAuthController struct {
	Username     string       `json:"-"`
	Password     string       `json:"-"`
	ValidateUser ValidateUser `json:"-"`
}

func (dhbac *HttpBasicAuthController) CreateSecurityContext(r *http.Request) {
	auth := r.Header.Get("Authorization")
	vals := extractUsernamePassword.FindAllStringSubmatch(auth, -1)

	if len(vals) == 1 && len(vals[0]) == 3 {
		dhbac.Username = vals[0][1]
		dhbac.Password = vals[0][2]
	}
}

func (dhbac *HttpBasicAuthController) HasPermission() bool {
	return dhbac.ValidateUser(dhbac.Username, dhbac.Password)
}

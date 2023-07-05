package goauthentication

import "net/http"

type ValidateUser func(username string, password string) bool

type HttpBasicAuthController struct {
	Username     string       `json:"-"`
	Password     string       `json:"-"`
	ValidateUser ValidateUser `json:"-"`
}

func (dhbac *HttpBasicAuthController) CreateSecurityContext(r *http.Request, header *http.Header) {
	username := header.Get("username")
	password := header.Get("password")

	if username != "" && password != "" {
		dhbac.SetUsernamePassword(username, password)
	}
}

func (dhbac *HttpBasicAuthController) SetUsernamePassword(username string, password string) {
	dhbac.Username = username
	dhbac.Password = password
}

func (dhbac *HttpBasicAuthController) HasPermission() bool {
	return dhbac.ValidateUser(dhbac.Username, dhbac.Password)
}

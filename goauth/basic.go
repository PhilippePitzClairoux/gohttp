package goauth

import (
	"net/http"
)

// HttpBasicAuthController defines the fields that must be used in order to have basic auth in an application
type HttpBasicAuthController struct {
	Username     string       `json:"-"`
	Password     string       `json:"-"`
	ValidateUser ValidateUser `json:"-"`
}

// CreateSecurityContext parses the request headers for an Authorization and extracts a username/password
func (dhbac *HttpBasicAuthController) CreateSecurityContext(r *http.Request) {
	auth := r.Header.Get("Authorization")
	vals := ExtractUsernamePassword.FindAllStringSubmatch(auth, -1)

	if len(vals) == 1 && len(vals[0]) == 3 {
		dhbac.Username = vals[0][1]
		dhbac.Password = vals[0][2]
	}
}

// HasPermission calls the ValidateUser method with the supplied username/password from CreateSecurityContext in order
// to determine if the user has access
func (dhbac *HttpBasicAuthController) HasPermission() bool {
	return dhbac.ValidateUser(dhbac.Username, dhbac.Password)
}

// PostLogin is useless for HttpBasicAuthController since validation is done using username/password every request
func (dhbac *HttpBasicAuthController) PostLogin() interface{} {
	return nil
}

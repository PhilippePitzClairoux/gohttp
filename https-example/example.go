package main

import (
	"crypto/tls"
	"github.com/PhilippePitzClairoux/gohttp"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
	"github.com/PhilippePitzClairoux/gohttp/goerrors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

type TestHandler struct {
	Name       string            `json:"name,omitempty"`
	FamilyName string            `json:"familyName,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

func (r TestHandler) GetMyEntity() TestHandler {
	return r
}

func (r TestHandler) GetsMyEntities(str string) []string {
	return []string{"A", "B", "C"}
}

func (r TestHandler) Post() TestHandler {
	return r
}

func (r TestHandler) Delete(id int) string {
	return "del called!"
}

func (r TestHandler) Patch(str string, float float64) string {
	return "patch called"
}

type AuthController struct {
	*goauth.JwtTokenAuthController
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func (r *AuthController) PostLogin() interface{} {
	var claims = r.GetClaims(r.Username, r.Password)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := r.GetSecret(token)
	if err != nil {
		return goerrors.NewInternalServerError(err.Error())
	}

	str, err := token.SignedString(key)
	if err != nil {
		return goerrors.NewInternalServerError(err.Error())
	}

	return AuthController{
		Token: str,
	}
}

func GetClaims(username string, password string) jwt.RegisteredClaims {
	return goauth.NewClaimBase(
		jwt.NewNumericDate(time.Now().Add(time.Hour*24)),
		"me",
		"me",
		"1",
		[]string{"t", "t2", "t3"},
	)
}

func GetSecret(token *jwt.Token) (interface{}, error) {
	return []byte("123123123123"), nil
}

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{goauth.CreateCertificate()},
		KeyLogWriter:       os.Stdout,
	}

	auth := &AuthController{
		JwtTokenAuthController: &goauth.JwtTokenAuthController{
			HasError:  false,
			GetSecret: GetSecret,
			GetClaims: GetClaims,
			Token:     jwt.New(jwt.SigningMethodHS512),
		},
	}

	srv := gohttp.NewTLSServer(":https", tlsConf)
	err := srv.RegisterAuthenticationMiddleware(auth)
	if err != nil {
		log.Fatalf("Cannot register new endpoints : %s", err)
	}

	err = srv.RegisterEndpoints("/test", TestHandler{})
	if err != nil {
		log.Fatalf("Cannot register new endpoints : %s", err)
	}

	log.Fatal(srv.ListenAndServeTLS("", ""))
}

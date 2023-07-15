package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/PhilippePitzClairoux/gohttp"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
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

func PerformLogin(req *http.Request) bool {
	// check if user is valid or not
	_map := make(map[string]interface{})
	err := json.NewDecoder(req.Body).Decode(&_map)
	if err != nil {
		return false
	}

	return _map["username"].(string) == "admin" &&
		_map["password"].(string) == "admin"
}

func GetClaims() jwt.Claims {
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

func GenerateSignedJWT(request *http.Request) (string, error) {
	secret, _ := GetSecret(nil)
	return jwt.NewWithClaims(jwt.SigningMethodHS512, goauth.NewClaimBase(
		jwt.NewNumericDate(time.Now().Add(time.Hour*24)),
		"me",
		"me",
		"100",
		[]string{"t", "t1", "t2", "t3"},
	)).SignedString(secret)
}

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{goauth.CreateCertificate()},
		KeyLogWriter:       os.Stdout,
	}

	srv := gohttp.NewTLSServer(":https", tlsConf)
	srv.RegisterAuthenticationMiddleware(&goauth.JwtMiddleware{
		HasError:  false,
		GetSecret: GetSecret,
		GetClaims: GetClaims,
		Token:     jwt.New(jwt.SigningMethodHS512),
	}, PerformLogin, GenerateSignedJWT)

	err := srv.RegisterEndpoints("/test/", TestHandler{})
	if err != nil {
		log.Fatalf("Cannot register new endpoints : %s", err)
	}

	log.Fatal(srv.ListenAndServeTLS("", ""))
}

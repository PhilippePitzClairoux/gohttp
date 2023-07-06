package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PhilippePitzClairoux/gohttp"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
	"log"
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
	fmt.Println(r)
	return r
}

func (r TestHandler) Delete(id int) string {
	return "del called!"
}

func (r TestHandler) Patch(str string, float float64) string {
	return "patch called"
}

func main() {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}
	auth := goauth.HttpBasicAuthController{
		ValidateUser: func(username string, password string) bool {
			return username == "admin" && password == "admin"
		},
	}

	srv := gohttp.NewHttpsServer(8080, tlsConf)
	vals, _ := gohttp.NewHttpServerEndpoint("/test", TestHandler{})
	srv.RegisterEndpoints(
		vals,
	)

	srv.RegisterAuthController(&auth)

	vals, _ = gohttp.NewHttpServerEndpoint("/", TestHandler{})
	srv.RegisterEndpoints(
		vals,
	)

	if err := srv.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Fatal("Cannot start HttpServer : ", err)
	}
}

package main

import (
	"github.com/PhilippePitzClairoux/go-http-server/gohttp"
	"github.com/PhilippePitzClairoux/go-http-server/testpackage"
	"log"
)

func main() {
	srv := gohttp.NewHttpServer(8080)
	vals, _ := gohttp.NewHttpServerEndpoint("/test", testpackage.TestHandler{})

	srv.RegisterEndpoints(
		vals,
	)

	if err := srv.ServeAndListen(); err != nil {
		log.Fatal("Cannot start HttpServer : ", err)
	}
}

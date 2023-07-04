package main

import (
	"go-http-server/gohttp"
	"go-http-server/testpackage"
)

func main() {
	srv := gohttp.NewHttpServer(8080)
	vals, _ := gohttp.NewHttpServerEndpoint("/test", testpackage.TestHandler{})

	srv.RegisterEndpoints(
		vals,
	)

	srv.ServeAndListen()
}

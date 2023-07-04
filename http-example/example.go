package main

import (
	"github.com/PhilippePitzClairoux/gohttp"
	"log"
)

type TestHandler struct {
	Name       string            `json:"Name" json:"Name,omitempty"`
	FamilyName string            `json:"FamilyName" json:"FamilyName,omitempty"`
	Properties map[string]string `json:"Properties" json:"Properties,omitempty"`
}

func (r TestHandler) GetMyEntity(str string, i int) TestHandler {
	return r
}

func (r TestHandler) GetsMyEntities(str string) []string {
	return []string{"A", "B", "C"}
}

func (r TestHandler) Post(str string, str2 string) string {
	return "post called!"
}

func (r TestHandler) Delete(id int) string {
	return "del called!"
}

func (r TestHandler) Patch(str string, float float64) string {
	return "patch called"
}

func main() {
	srv := gohttp.NewHttpServer(8080)
	vals, _ := gohttp.NewHttpServerEndpoint("/test", TestHandler{})

	srv.RegisterEndpoints(
		vals,
	)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Cannot start HttpServer : ", err)
	}
}

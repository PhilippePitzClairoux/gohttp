package gohttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type placeHolder struct {
	_type reflect.Kind
	value any
}

type HttpServer struct {
	server              *http.Server
	httpServerEndpoints []*HttpServerEndpoint
}

func NewHttpServer(port int) *HttpServer {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	//scan for endpoints
	//register endpoints
	server := &HttpServer{
		httpServerEndpoints: []*HttpServerEndpoint{},
	}

	server.server = &http.Server{
		Addr:    addr,
		Handler: &internalDispatcher{server},
	}

	return server
}

func (hs *HttpServer) ServeAndListen() {
	if err := hs.server.ListenAndServe(); err != nil {
		fmt.Println("Cannot start HttpServer : ", err)
	}
}

func (hs *HttpServer) RegisterEndpoints(endpoints *[]HttpServerEndpoint) {
	for _, endpoint := range *endpoints {
		hs.RegisterEndpoint(&endpoint)
	}
}

func (hs *HttpServer) RegisterEndpoint(endpoint *HttpServerEndpoint) {
	endpoint._uri = CompileUri(endpoint.name)
	hs.httpServerEndpoints = append(hs.httpServerEndpoints, endpoint)
}

func containsSupportedPlaceHolders(s string) bool {
	for key, _ := range Placeholders {
		if strings.Contains(s, key) {
			return true
		}
	}
	return false
}

func ParseValue(value string, _type reflect.Kind) (any, error) {
	var val any
	var err error

	switch _type {
	case reflect.String:
		val = value
	case reflect.Int:
		val, err = strconv.Atoi(value)
	case reflect.Float64:
		val, err = strconv.ParseFloat(value, 64)
	case reflect.Struct:
		err = json.Unmarshal([]byte(value), &val)
	case reflect.Invalid:
		val = ""
		err = errors.New("cannot parse a struct")

	}

	return val, err
}

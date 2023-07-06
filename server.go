package gohttp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
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
	Server          *http.Server
	sortedEndpoints map[string][]*HttpServerEndpoint
	Auth            goauth.HttpAuthController
}

func NewHttpServer(port int) *HttpServer {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	server := &HttpServer{
		sortedEndpoints: make(map[string][]*HttpServerEndpoint),
	}

	server.Server = &http.Server{
		Addr:    addr,
		Handler: &InternalDispatcher{server},
	}

	return server
}

func NewHttpsServer(port int, conf *tls.Config) *HttpServer {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	server := &HttpServer{
		sortedEndpoints: make(map[string][]*HttpServerEndpoint, 0),
	}

	server.Server = &http.Server{
		Addr:      addr,
		Handler:   &InternalDispatcher{server},
		TLSConfig: conf,
	}

	return server
}

func (hs *HttpServer) ListenAndServe() error {
	fmt.Println("Starting server : ", hs.Server.Addr)
	return hs.Server.ListenAndServe()
}

func (hs *HttpServer) ListenAndServeTLS(cert string, key string) error {
	fmt.Println("Starting server : ", hs.Server.Addr)
	return hs.Server.ListenAndServeTLS(cert, key)
}

func (hs *HttpServer) RegisterAuthController(controller goauth.HttpAuthController) {
	hs.Auth = controller
}

func (hs *HttpServer) RegisterEndpoints(endpoints *[]*HttpServerEndpoint) {
	for _, endpoint := range *endpoints {
		hs.RegisterEndpoint(endpoint)
	}
}

func (hs *HttpServer) RegisterEndpoint(endpoint *HttpServerEndpoint) {
	endpoint.hseUri = CompileUri(endpoint.name)
	hs.sortedEndpoints[endpoint.hseUri.baseUri] = append(hs.sortedEndpoints[endpoint.hseUri.baseUri], endpoint)
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
	case reflect.Bool:
		val, err = strconv.ParseBool(value)
	case reflect.Struct:
	case reflect.Invalid:
	default:
		val = ""
		err = errors.New("invalid type")
	}

	return val, err
}

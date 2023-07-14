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

// placeHolder defines a struct that's used by endpoint URI's to know the types of templated params and help with the parsing
type placeHolder struct {
	_type reflect.Kind
	value any
}

// HttpServer defines various fields that can be used by your server
type HttpServer struct {
	Server          *http.Server
	sortedEndpoints map[string][]*HttpServerEndpoint
	AuthControllers map[string]*goauth.HttpAuthController
}

// NewHttpServer creates a new server that can then be configured
func NewHttpServer(port int) *HttpServer {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	server := &HttpServer{
		sortedEndpoints: make(map[string][]*HttpServerEndpoint),
		AuthControllers: make(map[string]*goauth.HttpAuthController, 0),
	}

	server.Server = &http.Server{
		Addr:    addr,
		Handler: &internalDispatcher{server},
	}

	return server
}

// NewTLSServer creates a new http(s) server with a TLS configuration.
func NewTLSServer(addr string, conf *tls.Config) *HttpServer {

	server := &HttpServer{
		sortedEndpoints: make(map[string][]*HttpServerEndpoint, 0),
		AuthControllers: make(map[string]*goauth.HttpAuthController, 0),
	}

	server.Server = &http.Server{
		Addr:      addr,
		Handler:   &internalDispatcher{server},
		TLSConfig: conf,
	}

	return server
}

// ListenAndServe starts the HTTP server
func (hs *HttpServer) ListenAndServe() error {
	fmt.Println("Starting server : ", hs.Server.Addr)
	return hs.Server.ListenAndServe()
}

// ListenAndServeTLS starts the TLS (https) server
func (hs *HttpServer) ListenAndServeTLS(cert string, key string) error {
	fmt.Println("Starting server : ", hs.Server.Addr)
	return hs.Server.ListenAndServeTLS(cert, key)
}

// RegisterAuthController adds a specific goauth.HttpAuthController to the current server. This controller will be used
// for every basePath you pass to the method.
func (hs *HttpServer) RegisterAuthController(controller goauth.HttpAuthController, authPath string, basePaths ...string) error {
	for _, basePath := range basePaths {
		hs.AuthControllers[basePath] = &controller
	}

	endpoints, err := NewHttpServerEndpoint(authPath, controller)
	if err == nil {
		hs.RegisterEndpoints(endpoints)
		return nil
	}

	return err
}

// RegisterEndpoints registers many endpoints to a server
func (hs *HttpServer) RegisterEndpoints(endpoints *[]*HttpServerEndpoint) {
	for _, endpoint := range *endpoints {
		hs.RegisterEndpoint(endpoint)
	}
}

// RegisterEndpoint registers a single endpoint to a server
func (hs *HttpServer) RegisterEndpoint(endpoint *HttpServerEndpoint) {
	endpoint.hseUri = compileUri(endpoint.name)
	hs.sortedEndpoints[endpoint.hseUri.baseUri] = append(hs.sortedEndpoints[endpoint.hseUri.baseUri], endpoint)
}

// containsSupportedPlaceHolders check's if a parameter is a templated value
func containsSupportedPlaceHolders(s string) bool {
	for key, _ := range Placeholders {
		if strings.Contains(s, key) {
			return true
		}
	}
	return false
}

// parseValue parses the string value to the specified _type
func parseValue(value string, _type reflect.Kind) (any, error) {
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

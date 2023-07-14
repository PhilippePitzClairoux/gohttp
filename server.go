package gohttp

import (
	"crypto/tls"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
	"net/http"
	"reflect"
)

type HttpServer struct {
	*http.Server
	mux  *http.ServeMux
	Auth *goauth.AuthenticationMiddleware
}

// NewHttpServer creates a new server that can then be configured
func NewHttpServer(addr string) *HttpServer {
	httpServer := &HttpServer{
		Server: &http.Server{
			Addr: addr,
		},
		mux: http.NewServeMux(),
	}

	httpServer.Handler = httpServer.mux
	return httpServer
}

// NewTLSServer creates a new http(s) server with a TLS configuration.
func NewTLSServer(addr string, conf *tls.Config) *HttpServer {
	server := NewHttpServer(addr)
	server.TLSConfig = conf

	return server
}

// RegisterEndpoints creates a new endpoint based off a HttpController.
// You can then register these controllers to a server
func (hs *HttpServer) RegisterEndpoints(basePath string, controller HttpController) error {
	hse := make([]*controllerEndpoint, 0)
	ctrlRef := reflect.TypeOf(controller)

	for i := 0; i < ctrlRef.NumMethod(); i++ {
		method := ctrlRef.Method(i)
		supportedMethod := getSupportedMethod(method.Name)

		// do nothing if the method name doesn't start with a supportedMethod
		if supportedMethod != "" {
			//always skip the first method since it's the struct
			val, err := newEndpointFromType(basePath, method.Type)
			if err != nil {
				return err
			}

			// set method
			val.method = supportedMethod
			val.function = method

			hse = append(hse, &val)
		}
	}

	hs.mux.Handle(basePath, controllerEndpoints{
		endpoints:     &hse,
		controllerRef: &controller,
		auth:          hs.Auth,
	})

	return nil
}

func (hs *HttpServer) RegisterAuthenticationMiddleware(t goauth.AuthenticationMiddleware) error {
	hs.Auth = &t
	return hs.RegisterEndpoints("/login", t)
}

package gohttp

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PhilippePitzClairoux/gohttp/goauth"
	"github.com/PhilippePitzClairoux/gohttp/goerrors"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// Placeholders are the templated types you can use in your controllers
var Placeholders = map[string]reflect.Kind{
	"{string}": reflect.String,
	"{int}":    reflect.Int,
	"{float}":  reflect.Float64,
	"{struct}": reflect.Struct,
	"{bool}":   reflect.Bool,
}

// HttpController This is just a generic interface to handle HttpControllers
type HttpController interface {
}

// HttpServerEndpoint defines the necessary fields that an endpoint must have in order to work.
// users of this framework don't create them manually.
type HttpServerEndpoint struct {
	name          string
	method        string
	hseUri        uri
	function      reflect.Method
	controllerRef *HttpController
}

// supportedMethods is a list of method names that functions must have in order to be registered as an endpoint
var supportedMethods []string
var byteBuffer bytes.Buffer
var byteEncoder *gob.Encoder

func init() {
	supportedMethods = []string{"Post", "Get", "Delete", "Put", "Patch"}
	byteEncoder = gob.NewEncoder(&byteBuffer)
}

// internalDispatcher handles the server logic
type internalDispatcher struct {
	Parent *HttpServer
}

// search searches for an endpoint based off a URI
func (id *internalDispatcher) search(uri string) []*HttpServerEndpoint {
	return id.Parent.sortedEndpoints[uri]
}

// ServeTLS serves an HTTPS/TLS server (basically does the same thing ServeHttp would)
func (id *internalDispatcher) ServeTLS(rw http.ResponseWriter, r *http.Request) {
	id.ServeHTTP(rw, r)
}

// ServeHTTP handles authentication, dispatching, controller execution and also error management.
func (id *internalDispatcher) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error = goerrors.NewNotFoundError("Controller not found")
	//var endpoint *HttpServerEndpoint
	requestUri := compileUri(r.RequestURI)
	fmt.Printf("Got a new request : %s %s\n", r.Method, r.RequestURI)

	endpoint, err := id.searchEndpoint(r, &requestUri)
	if err != nil {
		id.handleErrors(rw, err)
		return
	}

	auth := id.Parent.AuthControllers[endpoint.hseUri.baseUri]
	if auth != nil {
		err = goauth.AuthProxyMethod(r, auth)
		if err != nil {
			goto HandleError
		}
	}

	err = id.executeRequest(rw, r, endpoint, &requestUri)
	if err != nil {
		goto HandleError
	}

HandleError:
	if err != nil {
		id.handleErrors(rw, err)
		return
	}
}

// searchEndpoint searches for the endpoint that matches the URI.
// Implementation example :
// If uri = /test/1234/hehe, we're going to search for /test/1234/hehe, /test/1234, /test, /
func (id *internalDispatcher) searchEndpoint(r *http.Request, requestUri *uri) (*HttpServerEndpoint, error) {
	var err error
	var endpoint *HttpServerEndpoint

	if strings.Count(r.RequestURI, "/") == 1 {
		endpoints := id.search(r.RequestURI)
		if len(endpoints) > 0 {
			// found the correct baseUri
			endpoint, err = id.findEndpoint(r, endpoints, requestUri)
		} else if len(id.search("/")) > 0 {
			endpoint, err = id.findEndpoint(r, id.search("/"), requestUri)
		}
	} else {
		// or else we search by removing parameters.
		reg, _ := regexp.Compile(`(/[\w-\\]+)`)
		matches := reg.FindAllString(r.RequestURI, -1)

		for index := len(matches) - 1; index >= 0; index-- {
			endpoints := id.search(strings.Join(matches[:index], ""))

			if len(endpoints) > 0 {
				endpoint, err = id.findEndpoint(r, endpoints, requestUri)
				break
			}
		}
	}
	return endpoint, err
}

// findEndpoint looks for the endpoint that matches the URI once the basePath has been found
func (id *internalDispatcher) findEndpoint(r *http.Request, endpoints []*HttpServerEndpoint, requestUri *uri) (*HttpServerEndpoint, error) {

	for _, endpoint := range endpoints {
		if endpoint.hseUri.uriMatches(requestUri) && endpoint.methodMatches(r.Method) {
			//err = id.executeRequest(rw, r, endpoint, requestUri)
			return endpoint, nil
		}
	}
	return nil, errors.New("no endpoint found")
}

// handleErrors handles errors...
func (id *internalDispatcher) handleErrors(rw http.ResponseWriter, err error) {
	var statusCode int

	if ghe, ok := err.(goerrors.GenericHttpError); ok {
		statusCode = ghe.GetStatusCode()
	} else {
		statusCode = http.StatusNotImplemented
	}

	rw.WriteHeader(statusCode)
	fmt.Println("There was an error : ", err)
}

// executeRequest takes an endpoint and call's the method related to it (keep in mind authentication has already been managed by then).
// the body of the request is also parsed and passed as a parameter if it's present/valid
func (id *internalDispatcher) executeRequest(rw http.ResponseWriter, r *http.Request, endpoint *HttpServerEndpoint, requestUri *uri) error {
	params, err := getValuesForMethodCall(endpoint.hseUri, requestUri)
	if err != nil {
		return err
	}

	// handle request body and potentially add it to the function call
	controller := endpoint.parseRequestPayload(r)
	err = addControllerReference(&controller, *endpoint)
	if err != nil {
		return err
	}

	answer := endpoint.function.Func.Call(
		append([]reflect.Value{reflect.ValueOf(controller)},
			params...,
		),
	)

	byteBuffer.Reset()
	err = byteEncoder.EncodeValue(answer[0])

	rw.Header().Set("Content-Type", "text/json")

	rw.WriteHeader(http.StatusOK)
	jsonBytes, err := json.Marshal(answer[0].Interface()) //byteBuffer.Bytes())
	if _, err = rw.Write(jsonBytes); err != nil {
		return err
	}

	fmt.Println("Dispatched to endpoint : ", endpoint.name, "\n")
	return nil
}

// addControllerReference adds a reference to the endpoint inside the controller
// this is done in order to add security context or other metadata from the request the controller
// might need in order to do its job
func addControllerReference(controller *HttpController, endpoint HttpServerEndpoint) error {
	controllerValue := reflect.Indirect(reflect.ValueOf(*controller))
	endpointValue := reflect.Indirect(reflect.ValueOf(*endpoint.controllerRef))

	// should double check we now have a struct (could still be anything)
	if controllerValue.Kind() != reflect.Struct || endpointValue.Kind() != reflect.Struct {
		return goerrors.NewBadRequestError("invalid type inside payload found")
	}

	controllerFieldTypes := getFields(controllerValue)
	endpointFieldTypes := getFields(endpointValue)

	for _, value := range controllerFieldTypes {
		val := contains(value, &endpointFieldTypes)
		if val != nil && value.CanSet() {
			value.Set(*val)
		}
	}

	return nil
}

// contains checks if the type of value matches a type inside of values.
// This method is strictly used to add a reference to an endpoint inside a controller instance.
func contains(value reflect.Value, values *[]reflect.Value) *reflect.Value {
	for _, val := range *values {
		if value.Type() == val.Type() {
			return &val
		}
	}

	return nil
}

// getFields retuns a list of fields for a structure
func getFields(value reflect.Value) []reflect.Value {
	fields := make([]reflect.Value, 0)
	for index := 0; index < value.NumField(); index++ {
		field := value.Field(index)
		fields = append(fields, field)
	}

	return fields
}

// getValuesForMethodCall parses the request URI to match the parameters of the endpoint URI.
// that way we can call the method related to the endpoint and pass the list that's return as
// a parameter
func getValuesForMethodCall(endpoint uri, request *uri) ([]reflect.Value, error) {
	params := make([]reflect.Value, 0)
	for i, val := range request.params {
		if reflect.TypeOf(endpoint.params[i]).Kind() == reflect.Struct {
			value, err := parseValue(val.(string), endpoint.params[i].(placeHolder)._type)
			if err != nil {
				return nil, err
			}

			params = append(params, reflect.ValueOf(value))
		}
	}

	return params, nil
}

// parseRequestPayload parses the body of the http call and transforms it to a HttpServerEndpoint.
// That way, when a user defines a new controller, they also define a payload that their endpoints will be able
// to get when executing an http request
func (hse *HttpServerEndpoint) parseRequestPayload(req *http.Request) HttpController {
	t := reflect.TypeOf(*hse.controllerRef)
	unmarshalled := reflect.New(t)
	body := req.Body
	content, err := io.ReadAll(body)

	if err == nil {
		_ = json.Unmarshal(content, unmarshalled.Interface())
	}
	_ = body.Close()
	return unmarshalled.Elem().Interface()
}

// methodMatches checks if the HttpServerEndpoint method matches the one we just received
func (hse *HttpServerEndpoint) methodMatches(method string) bool {
	return strings.ToLower(hse.method) == strings.ToLower(method)
}

// NewHttpServerEndpoint creates a new endpoint based off a HttpController.
// You can then register these controllers to a server
func NewHttpServerEndpoint(basePath string, controller HttpController) (*[]*HttpServerEndpoint, error) {
	hse := make([]*HttpServerEndpoint, 0)
	ctrlRef := reflect.TypeOf(controller)

	for i := 0; i < ctrlRef.NumMethod(); i++ {
		method := ctrlRef.Method(i)
		supportedMethod := getSupportedMethod(method.Name)

		// do nothing if the method name doesn't start with a supportedMethod
		if supportedMethod != "" {
			//always skip the first method since it's the struct
			val, err := newEndpointFromType(basePath, method.Type)
			if err != nil {
				return nil, err
			}

			// set method
			val.method = supportedMethod
			val.function = method
			val.controllerRef = &controller

			hse = append(hse, &val)
		}
	}

	return &hse, nil
}

// newEndpointFromType creates an endpoint based off a controller
func newEndpointFromType(name string, p reflect.Type) (HttpServerEndpoint, error) {

	for i := 1; i < p.NumIn(); i++ {
		val, err := getPlaceholderFromType(p.In(i).Kind())
		if err != nil {
			return HttpServerEndpoint{}, err
		}

		if name[len(name)-1] != '/' {
			name += "/"
		}
		name += fmt.Sprintf("%s", val)
	}

	return HttpServerEndpoint{
		name:   name,
		hseUri: compileUri(name),
	}, nil
}

// getPlaceholderFromType returns the matching placeholder string based off the placeholder type in parameter
func getPlaceholderFromType(p reflect.Kind) (string, error) {
	for key, val := range Placeholders {
		if val == p {
			return key, nil
		}
	}
	return "", errors.New("invalid Kind used for placeholder")
}

// getSupportedMethod returns the method string if it's contained in supportedMethods
func getSupportedMethod(s string) string {
	for _, method := range supportedMethods {
		if strings.Index(s, method) == 0 {
			return method
		}
	}
	return ""
}

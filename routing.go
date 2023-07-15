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
	"strconv"
	"strings"
)

// placeHolder defines a struct that's used by endpoint URI's to know the types of templated params and help with the parsing
type placeHolder struct {
	_type reflect.Kind
	value any
}

// Placeholders are the templated types you can use in your controllers
var Placeholders = map[string]reflect.Kind{
	"{string}": reflect.String,
	"{int}":    reflect.Int,
	"{float}":  reflect.Float64,
	"{struct}": reflect.Struct,
	"{bool}":   reflect.Bool,
}

type controllerEndpoint struct {
	name     string
	method   string
	hseUri   uri
	function reflect.Method
}

type controllerEndpoints struct {
	endpoints     *[]*controllerEndpoint
	controllerRef *HttpController
	serverRef     *HttpServer
}

type HttpController interface{}

// supportedMethods is a list of method names that functions must have in order to be registered as an endpoint
var supportedMethods []string
var byteBuffer bytes.Buffer
var byteEncoder *gob.Encoder

func init() {
	supportedMethods = []string{"Post", "Get", "Delete", "Put", "Patch"}
	byteEncoder = gob.NewEncoder(&byteBuffer)
}

// ServeTLS serves an HTTPS/TLS server (basically does the same thing ServeHttp would)
func (ces controllerEndpoints) ServeTLS(rw http.ResponseWriter, r *http.Request) {
	ces.ServeHTTP(rw, r)
}

// ServeHTTP handles authentication, dispatching, controller execution and also error management.
func (ces controllerEndpoints) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error = goerrors.NewNotFoundError("Controller not found")
	//var endpoint *controllerEndpoint
	requestUri := compileUri(r.RequestURI)
	fmt.Printf("Got a new request : %s %s\n", r.Method, r.RequestURI)

	endpoint, err := ces.searchEndpoint(r, &requestUri)
	if err != nil {
		ces.handleErrors(rw, err)
		return
	}

	auth := ces.serverRef.Auth
	if auth != nil {
		err = goauth.AuthProxyMethod(r, auth)
		if err != nil {
			goto HandleError
		}
	}

	err = ces.executeRequest(rw, r, endpoint, &requestUri)
	if err != nil {
		goto HandleError
	}

HandleError:
	if err != nil {
		ces.handleErrors(rw, err)
		return
	}
}

// searchEndpoint searches for the endpoint that matches the URI.
// Implementation example :
// If uri = /test/1234/hehe, we're going to search for /test/1234/hehe, /test/1234, /test, /
func (ces controllerEndpoints) searchEndpoint(r *http.Request, requestUri *uri) (*controllerEndpoint, error) {

	for _, endpoint := range *ces.endpoints {
		if endpoint.hseUri.uriMatches(requestUri) {
			return endpoint, nil
		}
	}

	return nil, goerrors.NewNotFoundError("no matching endpoint found")
}

// findEndpoint looks for the endpoint that matches the URI once the basePath has been found
func (ces controllerEndpoints) findEndpoint(r *http.Request, endpoints []*controllerEndpoint, requestUri *uri) (*controllerEndpoint, error) {

	for _, endpoint := range endpoints {
		if endpoint.hseUri.uriMatches(requestUri) && endpoint.methodMatches(r.Method) {
			//err = id.executeRequest(rw, r, endpoint, requestUri)
			return endpoint, nil
		}
	}
	return nil, errors.New("no endpoint found")
}

// handleErrors handles errors...
func (ces controllerEndpoints) handleErrors(rw http.ResponseWriter, err error) {
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
func (ces controllerEndpoints) executeRequest(rw http.ResponseWriter, r *http.Request, endpoint *controllerEndpoint, requestUri *uri) error {
	params, err := getValuesForMethodCall(endpoint.hseUri, requestUri)
	if err != nil {
		return err
	}

	// handle request body and potentially add it to the function call
	controller := endpoint.parseRequestPayload(r, *ces.controllerRef)
	//err = addControllerReference(&controller, *endpoint)
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

	fmt.Printf("Dispatched to endpoint : %s\n\n", endpoint.name)
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

// parseRequestPayload parses the body of the http call and transforms it to a controllerEndpoint.
// That way, when a user defines a new controller, they also define a payload that their endpoints will be able
// to get when executing an http request
func (hse *controllerEndpoint) parseRequestPayload(req *http.Request, bodyRef HttpController) interface{} {
	t := reflect.TypeOf(bodyRef)
	unmarshalled := reflect.New(t)
	body := req.Body
	content, err := io.ReadAll(body)

	if err == nil {
		_ = json.Unmarshal(content, unmarshalled.Interface())
	}
	_ = body.Close()
	return unmarshalled.Elem().Interface()
}

// methodMatches checks if the controllerEndpoint method matches the one we just received
func (hse *controllerEndpoint) methodMatches(method string) bool {
	return strings.ToLower(hse.method) == strings.ToLower(method)
}

// newEndpointFromType creates an endpoint based off a controller
func newEndpointFromType(name string, p reflect.Type) (controllerEndpoint, error) {

	for i := 1; i < p.NumIn(); i++ {
		val, err := getPlaceholderFromType(p.In(i).Kind())
		if err != nil {
			return controllerEndpoint{}, err
		}

		if name[len(name)-1] != '/' {
			name += "/"
		}
		name += fmt.Sprintf("%s", val)
	}

	return controllerEndpoint{
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

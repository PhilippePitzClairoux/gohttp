package gohttp

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	goerrors "github.com/PhilippePitzClairoux/gohttp/errors"
	"net/http"
	"reflect"
	"strings"
)

var Placeholders = map[string]reflect.Kind{
	"{string}": reflect.String,
	"{int}":    reflect.Int,
	"{float}":  reflect.Float64,
	"{struct}": reflect.Struct,
}

type HttpController interface {
}

type HttpServerEndpoint struct {
	name          string
	method        string
	hseUri        Uri
	function      reflect.Method
	controllerRef HttpController
}

var supportedMethods []string
var byteBuffer bytes.Buffer
var byteEncoder *gob.Encoder

func init() {
	supportedMethods = []string{"Post", "Get", "Delete", "Put", "Patch"}
	byteEncoder = gob.NewEncoder(&byteBuffer)
}

type InternalDispatcher struct {
	Parent *HttpServer
}

func (id *InternalDispatcher) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error = goerrors.NewNotFoundError("Controller not found")

	fmt.Printf("Got a new request : %s %s\n", r.Method, r.RequestURI)
	for _, endpoint := range id.Parent.Endpoints {
		requestUri := CompileUri(r.RequestURI)

		if endpoint.hseUri.uriMatches(&requestUri) && endpoint.methodMatches(r.Method) {
			err = id.executeRequest(rw, r, endpoint, requestUri)
		}
	}

	if err != nil {
		id.handleErrors(rw, err)
	}
}

func (id *InternalDispatcher) handleErrors(rw http.ResponseWriter, err error) {
	var statusCode int

	if ise, ok := err.(goerrors.InternalServerError); ok {
		statusCode = ise.StatusCode
	} else if nfe, ok := err.(goerrors.NotFoundError); ok {
		statusCode = nfe.StatusCode
	} else {
		statusCode = http.StatusNotImplemented
	}

	rw.WriteHeader(statusCode)
	fmt.Println("There was an error : ", err)
}

func (id *InternalDispatcher) executeRequest(rw http.ResponseWriter, r *http.Request, endpoint *HttpServerEndpoint, requestUri Uri) error {
	params, err := getValuesForMethodCall(endpoint.hseUri, requestUri)
	if err != nil {
		return err
	}

	// handle request body and potentially add it to the function call
	controller := endpoint.ParseRequestPayload(r)
	answer := endpoint.function.Func.Call(
		append([]reflect.Value{reflect.ValueOf(controller)},
			params...,
		),
	)

	byteBuffer.Reset()
	err = byteEncoder.EncodeValue(answer[0])

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(byteBuffer.Bytes()); err != nil {
		return err
	}

	fmt.Println("Dispatched to endpoint : ", endpoint.name, "\n")
	return nil
}

func getValuesForMethodCall(endpoint Uri, request Uri) ([]reflect.Value, error) {
	params := make([]reflect.Value, 0)
	for i, val := range request.params {
		if reflect.TypeOf(endpoint.params[i]).Kind() == reflect.Struct {
			value, err := ParseValue(val.(string), endpoint.params[i].(placeHolder)._type)
			if err != nil {
				return nil, err
			}

			params = append(params, reflect.ValueOf(value))
		}
	}

	return params, nil
}

func (hse *HttpServerEndpoint) ParseRequestPayload(req *http.Request) any {
	t := reflect.ValueOf(hse.controllerRef)
	body := req.Body
	content := make([]byte, 0)
	var err error = nil
	var unmarshalled any

	_, err = body.Read(content)
	if err != nil {
		// if there's an error during the parsing of the payload - return empty struct
		fmt.Println("Could not parse body : ", err)
		unmarshalled = hse.controllerRef
	} else {
		unmarshalled = reflect.New(t.Type()).Interface()
		_ = json.Unmarshal(content, &unmarshalled)
	}

	_ = body.Close()
	return unmarshalled
}

func (hse *HttpServerEndpoint) methodMatches(method string) bool {
	return strings.ToLower(hse.method) == strings.ToLower(method)
}

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
			val.controllerRef = controller

			hse = append(hse, &val)
		}
	}

	return &hse, nil
}

func newEndpointFromType(name string, p reflect.Type) (HttpServerEndpoint, error) {

	for i := 1; i < p.NumIn(); i++ {
		val, err := getPlaceholderFromType(p.In(i).Kind())
		if err != nil {
			return HttpServerEndpoint{}, err
		}

		name += fmt.Sprintf("/%s", val)
	}

	return HttpServerEndpoint{
		name:   name,
		hseUri: CompileUri(name),
	}, nil
}

func getPlaceholderFromType(p reflect.Kind) (string, error) {
	for key, val := range Placeholders {
		if val == p {
			return key, nil
		}
	}
	return "", errors.New("invalid Kind used for placeholder")
}

func getSupportedMethod(s string) string {
	for _, method := range supportedMethods {
		if strings.Index(s, method) == 0 {
			return method
		}
	}
	return ""
}

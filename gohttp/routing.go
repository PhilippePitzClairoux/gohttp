package gohttp

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
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
	_uri          Uri
	function      reflect.Method
	controllerRef HttpController
}

var supportedMethods = []string{"Post", "Get", "Delete", "Put", "Patch"}
var byteBuffer bytes.Buffer
var byteEncoder = gob.NewEncoder(&byteBuffer)

type internalDispatcher struct {
	parent *HttpServer
}

func (id *internalDispatcher) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, endpoint := range id.parent.httpServerEndpoints {
		if endpoint._uri.uriMatches(r.RequestURI) && endpoint.methodMatches(r.Method) {
			//handle match or else 404
			//TODO : call endpoint.method, create an array of params based off the request uri and potentially the body of the request
			requestUri := CompileUri(r.RequestURI)

			fmt.Println(requestUri)
			fmt.Println(endpoint._uri)
			params, err := getValuesForMethodCall(endpoint._uri, requestUri)
			if err != nil {
				fmt.Println("Cannot parse request : ", err)
				break
			}

			answer := endpoint.function.Func.Call(
				append([]reflect.Value{reflect.ValueOf(endpoint.controllerRef)},
					params...,
				),
			)

			byteBuffer.Reset()
			err = byteEncoder.EncodeValue(answer[0])
			if err != nil {
				fmt.Println("Cannot encode answer to bytes : ", err)
				break
			}

			rw.WriteHeader(http.StatusOK)
			rw.Write(byteBuffer.Bytes())
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
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

func (hse *HttpServerEndpoint) methodMatches(method string) bool {
	return strings.ToLower(hse.method) == strings.ToLower(method)
}

func NewHttpServerEndpoint(basePath string, controller HttpController) ([]HttpServerEndpoint, error) {
	hse := make([]HttpServerEndpoint, 0)
	ctrlRef := reflect.TypeOf(controller)

	for i := 0; i < ctrlRef.NumMethod(); i++ {
		method := ctrlRef.Method(i)
		supportedMethod := getSupportedMethod(method.Name)

		// do nothing if the method basePath doesn't start with a supported method type
		if supportedMethod != "" {
			//always skip the first method since it's the struct
			val, err := NewEndpointFromType(basePath, method.Type)
			if err != nil {
				return nil, err
			}

			// set method
			val.method = supportedMethod
			val.function = method
			val.controllerRef = controller

			hse = append(hse, val)
		}
	}

	return hse, nil
}

func NewEndpointFromType(name string, p reflect.Type) (HttpServerEndpoint, error) {

	for i := 1; i < p.NumIn(); i++ {
		val, err := getPlaceholderFromType(p.In(i).Kind())
		if err != nil {
			return HttpServerEndpoint{}, err
		}

		name += fmt.Sprintf("/%s", val)
	}

	return HttpServerEndpoint{
		name: name,
		_uri: CompileUri(name),
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

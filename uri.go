package gohttp

import (
	"reflect"
	"strings"
)

// uri of a request
type uri struct {
	fullUri            string
	baseUri            string
	params             []any
	hasTemplatedParams bool
}

// compileUri turns a string uri to a struct
func compileUri(value string) uri {
	cleanValue := cleanUri(value)
	_uri := uri{fullUri: value}
	templatedParams := false
	values := strings.Split(cleanValue, "/")

	for _, param := range values {

		if containsSupportedPlaceHolders(param) {
			_uri.params = append(_uri.params, placeHolder{Placeholders[param], nil})
			templatedParams = true
		} else {
			_uri.params = append(_uri.params, param)
		}
	}

	_uri.baseUri = getBaseUri(values)
	_uri.hasTemplatedParams = templatedParams
	return _uri
}

func getBaseUri(uri []string) string {
	output := "/"
	for _, param := range uri {
		if !containsSupportedPlaceHolders(param) {
			output += param
			break
		}
	}

	return output
}

func (u *uri) uriMatches(target *uri) bool {
	if u.hasTemplatedParams {
		//currentUri := compileUri(uri)

		if len(target.params) != len(u.params) {
			return false
		}

		for i, _ := range target.params {
			//here we basically ignore placeHolders since we can't really validate their value.
			//we could validate the type eventually but for now this should be enough
			//ex: /test/1234 should match /test/{int}   and /test/{int}/{string} should not match /test/2/{string}
			if reflect.TypeOf(u.params[i]).Kind() == reflect.String {
				if u.params[i] != target.params[i] {
					return false
				}
			} else {
				// params only contains placeHolder or string - the casting is technically safe
				valueToParse := target.params[i].(string)
				paramType := u.params[i].(placeHolder)._type
				_, err := parseValue(valueToParse, paramType)

				if err != nil {
					return false
				}
			}
		}
		return true
	}

	return u.fullUri == target.fullUri
}

func cleanUri(s string) string {

	if s == "/" {
		return s
	}

	if s[0] == '/' {
		s = s[1:]
	}

	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}

	return s
}

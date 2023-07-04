package gohttp

import (
	"reflect"
	"strings"
)

type Uri struct {
	fullUri            string
	params             []any
	hasTemplatedParams bool
}

func CompileUri(value string) Uri {
	value = cleanUri(value)
	_uri := Uri{fullUri: value}
	templatedParams := false

	for _, param := range strings.Split(value, "/") {

		if containsSupportedPlaceHolders(param) {
			_uri.params = append(_uri.params, placeHolder{Placeholders[param], nil})
			templatedParams = true
		} else {
			_uri.params = append(_uri.params, param)
		}
	}

	_uri.hasTemplatedParams = templatedParams
	return _uri
}

func (u *Uri) uriMatches(uri string) bool {
	if u.hasTemplatedParams {
		currentUri := CompileUri(uri)

		if len(currentUri.params) != len(u.params) {
			return false
		}

		for i, _ := range currentUri.params {
			//here we basically ignore placeHolders since we can't really validate their value.
			//we could validate the type eventually but for now this should be enough
			//ex: /test/1234 should match /test/{int}   and /test/{int}/yep should not match /test/2/abc
			if reflect.TypeOf(u.params[i]).Kind() == reflect.String {
				if u.params[i] != currentUri.params[i] {
					return false
				}
			} else {
				// params only contains placeHolder or string - the casting is technically safe
				valueToParse := currentUri.params[i].(string)
				paramType := u.params[i].(placeHolder)._type
				_, err := ParseValue(valueToParse, paramType)

				if err != nil {
					return false
				}
			}
		}
		return true
	}

	return u.fullUri == uri
}

func cleanUri(s string) string {
	if s[0] == '/' {
		s = s[1:]
	}

	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}

	return s
}

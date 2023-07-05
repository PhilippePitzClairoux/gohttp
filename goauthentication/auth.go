package goauthentication

import (
	goerrors "github.com/PhilippePitzClairoux/gohttp/goerrors"
	"net/http"
	"reflect"
)

type HttpAuthController interface {
	CreateSecurityContext(r *http.Request, header *http.Header)
	HasPermission() bool
}

func AuthProxy(r *http.Request, controller *HttpAuthController) error {
	value := reflect.ValueOf(controller)
	securityContextMethod := value.MethodByName("CreateSecurityContext")
	hasPermissionMethod := value.MethodByName("HasPermission")
	controllerReference := reflect.ValueOf(controller)

	securityContextMethodArguments := []reflect.Value{
		controllerReference,
		reflect.ValueOf(r),
		reflect.ValueOf(r.Header),
	}

	securityContextMethod.Call(securityContextMethodArguments)
	hasPermission := hasPermissionMethod.Call([]reflect.Value{controllerReference})

	if !hasPermission[0].Bool() {
		return goerrors.NewUnauthorizedError("Unauthorized")
	}

	return nil
}

package goauth

import (
	"fmt"
	"github.com/PhilippePitzClairoux/gohttp/goerrors"
	"github.com/huandu/go-clone"
	"net/http"
	"reflect"
)

type HttpAuthController interface {
	CreateSecurityContext(r *http.Request)
	HasPermission() bool
	PostLogin() interface{}
}

func AuthProxyMethod(r *http.Request, controller *HttpAuthController) error {
	clonedAuthController := clone.Clone(controller).(*HttpAuthController)
	value := reflect.ValueOf(clonedAuthController).Elem()
	securityContextMethod := value.MethodByName("CreateSecurityContext")
	hasPermissionMethod := value.MethodByName("HasPermission")
	//controllerReference := reflect.ValueOf(controller)

	securityContextMethodArguments := []reflect.Value{
		reflect.ValueOf(r),
	}

	securityContextMethod.Call(securityContextMethodArguments)
	hasPermission := hasPermissionMethod.Call([]reflect.Value{})

	fmt.Println("Auth is being checked by : ", reflect.TypeOf(clonedAuthController).Elem().Name())
	fmt.Println("Is user now authenticated : ", hasPermission[0].Bool())

	if hasPermission[0].Kind() == reflect.Bool && !hasPermission[0].Bool() {
		return goerrors.NewUnauthorizedError("Unauthorized")
	}

	return nil
}

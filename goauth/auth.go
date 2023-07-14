package goauth

import (
	"fmt"
	"github.com/PhilippePitzClairoux/gohttp/goerrors"
	"github.com/huandu/go-clone"
	"net/http"
	"reflect"
)

// AuthenticationMiddleware defines the methods that needs to be implemented in order to have a working
// authentication.
type AuthenticationMiddleware interface {
	CreateSecurityContext(r *http.Request)
	HasPermission() bool
}

// AuthProxyMethod this is the method that's called by the routing section of the code in order to validate
// that the user has the right's to call the endpoint (calls CreateSecurityContext and then HasPermission.
// Returns an error if HasPermission returns false
func AuthProxyMethod(r *http.Request, controller *AuthenticationMiddleware) error {
	clonedAuthController := clone.Clone(controller).(*AuthenticationMiddleware)
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

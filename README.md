# Go Http Server
This library tries to make golang http servers
"clean" and easily readable. There's still alot of
things to implement but there is currently the bare minimum.

Please refer to dummy.go and Example.go in order to learn
how this works.

## How to write your first controller!
In order to create a controller that can handle http calls,
you must create a struct that defines functions.
Those functions MUST start with one of the supported supported http method :
```golang
[]string{"Post", "Get", "Delete", "Put", "Patch"}

...

type TestHandler struct {
}

func (TestHandler) GetMyEntity(str string, i int) string {
    return "get called!"
}

func (TestHandler) PostMyEntity(str string, str2 string) string {
    return "post called!"
}

func (TestHandler) DeleteMyEntity(id int) string {
    return "del called!"
}

func (TestHandler) PatchMyEntity(str string, float float64) string {
    return "patch called"
}

...

func main() {
srv := gohttp.NewHttpServer(8080)
vals, _ := gohttp.NewHttpServerEndpoint("/test", testpackage.TestHandler{})

srv.RegisterEndpoints(
vals,
)

srv.ServeAndListen()
}

```
The library will generate endpoints based off the baseUrl passed in `RegisterEndpoints`
and the parameters of usable functions. So for example, `GetMyEntity` will be called
when the GET request matches the following path : `/test/{string}/{int}`.

The code abouve will generate the following endpoints :
```
GET    /test/{string}/{int}
POST   /test/{string}/{string}
DELETE /test/{int}
PATCH  /test/{string}/{float}
```

```bash
wget localhost:8080/test/PARAM_STRING/32
```

## TODO LIST
- Give access to request headers when calling endpoint method
- If the output of a method is a struct, parse it to JSON (same thing for body of a HTTP call)
- Handle various header related stuff (mostly for authentication purposes)
- Optimize code (use more pointers instead of copying most data)
- Reduce cognitive complexity of functions
- Support multi return statements and handle errors
- multi threading ???
- document exported methods
# Go Http Server
This library tries to make golang http servers
"clean" and easily readable. There's still alot of
things to implement but there is currently the bare minimum.

Please refer to http-example and https-example for examples.

## Project is useless
creator didn't know the ecosystem of golang enough when starting this project.
Was fun to do though

## Project features
- Aggregate endpoints for an entity in a single struct
- Create one file/struct per entity
- No more huge main methods defining endpoints
- Automatic dispatching to methods when a request is made
- parametrized endpoints
- authentication (BasicAuth & JWT only for now)

## Server structure
```
HttpServer
 |
 | - -> endpoints
           |
           | - -> basePath (ex: /test)
                     |
                     | - -> Get() *
                     | - -> id
                            | Post(id) **
                            | Delete(id) ***

MyCustomStruct
 |
 | - -> Get() *
 |
 | - -> Post(id string) **
 |
 | - -> Delete(id string) ***
```

## How to write your first controller!
In order to create a controller that can handle http calls,
you must create a struct that defines functions.
These functions MUST start with one of the supported http method :

```golang
[]string{"Post", "Get", "Delete", "Put", "Patch"}

...

type TestHandler struct {
	Name       string `json:"name"`
	properties map[string]string `json:"properties"`
}

func (r TestHandler) GetMyEntity(str string, i int) TestHandler {
    return r
}

func (r TestHandler) PostMyEntity(str string, str2 string) string {
    return "post called!"
}

func (r TestHandler) DeleteMyEntity(id int) string {
    return "del called!"
}

func (r TestHandler) PatchMyEntity(str string, float float64) string {
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
The library will generate endpoints based off the baseUrl passed in `RegisterEndpoints`/`RegisterEndpoint`
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

## Warnings
There can be possible mapping conflicts.
See following example :
```
POSSIBLE CONFLICTS :

GET /test/2/{string}
GET /test/{int}/{string}

GET /test/{int}
GET /test/2

GET /test
GET /{string}
```

Note : base url always have precedence over templated values.
Therefore `GET /test` will be used instead of `GET /{string}`
if `{string}` has the value `test`

## TODO LIST
- Give access to request headers when calling endpoint method
- ~~If the output of a method is a struct, parse it to JSON (same thing for body of a HTTP call)~~
- Handle various header related stuff (mostly for authentication purposes)
- ~~Optimize code (use more pointers instead of copying most data)~~
- ~~Reduce cognitive complexity of functions~~
- Support multi return statements and handle errors
- ~~multi threading ???~~
- ~~document exported methods~~
- ~~optimize search for endpoint~~
- ~~Add methods to handle authentication~~
- Add methods to check permissions

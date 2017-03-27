# swag

[![GoDoc](https://godoc.org/github.com/savaki/swag?status.svg)](https://godoc.org/github.com/savaki/swag)
[![Build Status](https://travis-ci.org/savaki/swag.svg?branch=master)](https://travis-ci.org/savaki/swag)

```swag``` is a lightweight library to generate swagger json for Go projects.  
 
No code generation, no framework constraints, just a simple swagger definition.


## Installation

```bash
go get github.com/savaki/swag
```


## Example

```go
func handlePet(w http.ResponseWriter, _ *http.Request) {
	// your code here
}

type Pet struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []string `json:"tags"`
}

func main() {
    // define our endpoints
    // 
    post := endpoint.New("post", "/", echo,
        endpoint.Summary("Add a new pet to the store"),
        endpoint.Description("Additional information on adding a pet to the store"),
        endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
        endpoint.Response(http.StatusOK, Pet{}, "Successfully added pet"),
    )
    get := endpoint.New("get", "/pet/{petId}", echo,
        endpoint.Summary("Find pet by ID"),
        endpoint.Path("petId", "integer", "ID of pet to return", true),
        endpoint.Response(http.StatusOK, Pet{}, "successful operation"),
    )
    
    // define the swagger api that will contain our endpoints
    // 
    api := swag.New(
      swag.Title("Swagger Petstore"),
      swag.Endpoints(post, get),
    )
    
    // iterate over each endpoint and add them to the default server mux
    // 
    api.Walk(func(path string, endpoint *swagger.Endpoint) {
        h := endpoint.Handler.(http.HandlerFunc)
        http.Handle(path, h)
    })
    
    // use the api to server the swagger.json file
    // 
    enableCors := true
    http.Handle("/swagger", api.Handler(enableCors))
    
    http.ListenAndServe(":8080", nil)
}
```

## Additional Examples

Examples for popular web frameworks can be found in the examples directory:

* [http.Server](examples/builtin/main.go)
* [Echo](examples/echo/main.go)
* [Gin](examples/gin/main.go)
* [Gorilla](examples/gorilla/main.go)
* [httprouter](examples/httprouter/main.go)


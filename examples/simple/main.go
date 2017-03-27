package main

import (
	"io"
	"net/http"

	"github.com/savaki/swaggering"
	"github.com/savaki/swaggering/endpoint"
	"github.com/savaki/swaggering/swagger"
)

func echo(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello World")
}

type Category struct {
	Id   int64  `json:"category"`
	Name string `json:"name"`
}

type Pet struct {
	Id        int64    `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []string `json:"tags"`
}

func main() {
	e := endpoint.New("post", "/", echo,
		endpoint.Summary("Add a new pet to the store"),
		endpoint.Description("Additional information on adding a pet to the store"),
		endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
		endpoint.Response(http.StatusOK, Pet{}, "Successfully added pet"),
	).Build()

	api := swaggering.New(
		swaggering.Endpoints(e),
	).Build()

	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(http.HandlerFunc)
		http.Handle(path, h)
	})

	http.Handle("/swagger", api.Handler(true))
	http.ListenAndServe(":8080", nil)
}

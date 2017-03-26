package main

import (
	"io"
	"net/http"

	"github.com/savaki/swaggering"
)

func echo(w http.ResponseWriter, req *http.Request) {
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
	endpoint := swaggering.NewEndpoint("post", "/", echo).
		Summary("Add a new pet to the store").
		Description("Additional information on adding a pet to the store").
		Body(Pet{}, "Pet object that needs to be added to the store", true).
		Response(http.StatusOK, Pet{}, "Successfully added pet").
		Endpoint

	api := &swaggering.Api{
		BasePath: "/api",
		CORS:     true,
		Endpoints: []swaggering.Endpoint{
			endpoint,
		},
	}

	api.Walk(func(path string, endpoint *swaggering.SwaggerEndpoints) {
		http.Handle(path, endpoint)
	})

	http.Handle("/swagger", api)
	http.ListenAndServe(":8080", nil)
}

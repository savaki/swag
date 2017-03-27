package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savaki/swag"
	"github.com/savaki/swag/endpoint"
	"github.com/savaki/swag/swagger"
)

func handle(c *gin.Context) {
	io.WriteString(c.Writer, "Insert your code here")
}

// Category example from the swagger pet store
type Category struct {
	ID   int64  `json:"category"`
	Name string `json:"name"`
}

// Pet example from the swagger pet store
type Pet struct {
	ID        int64    `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []string `json:"tags"`
}

func main() {
	post := endpoint.New("post", "/pet", "Add a new pet to the store",
		endpoint.Handler(handle),
		endpoint.Description("Additional information on adding a pet to the store"),
		endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
		endpoint.Response(http.StatusOK, Pet{}, "Successfully added pet"),
	)
	get := endpoint.New("get", "/pet/{petId}", "Find pet by ID",
		endpoint.Handler(handle),
		endpoint.Path("petId", "integer", "ID of pet to return", true),
		endpoint.Response(http.StatusOK, Pet{}, "successful operation"),
	)

	api := swag.New(
		swag.Endpoints(post, get),
	)

	router := gin.New()
	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
	})

	enableCors := true
	router.GET("/swagger", gin.WrapH(api.Handler(enableCors)))

	http.ListenAndServe(":8080", router)
}

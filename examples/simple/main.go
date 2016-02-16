package main

import (
	"io"
	"net/http"

	"github.com/savaki/swaggering"
)

type Owner struct {
	Name string
}

type Pet struct {
	String string
	Int    int64 `required:"true"`
	Owner  Owner
	Owners []Owner
	Ptr    *Owner
	Ptrs   []*Owner
}

func echo(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func main() {
	api, _ := swaggering.NewApi(
		swaggering.ApiBasePath("/api"),
	)

	api.EndpointFunc("post", "/pet", echo,
		swaggering.Summary("Add a new pet to the store"),
	)
	api.EndpointFunc("get", "/pet/findByStatus", echo,
		swaggering.Summary("Finds Pets by status"),
		swaggering.Description("Multiple status values can be provided with comma separated strings"),
		swaggering.Param(
			swaggering.ParamDescription("adding the thing to the thing"),
			swaggering.ParamType(Owner{}),
		),
		swaggering.Response(http.StatusOK, "successful operation", swaggering.Type(Pet{})),
		swaggering.Response(http.StatusBadRequest, "Invalid status value"),
	)

	api.Walk(func(path string, endpoints *swaggering.Endpoints) {
		http.Handle(path, endpoints)
	})

	// render swagger with cors
	http.HandleFunc("/swagger", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		api.ServeHTTP(w, req)
	})

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"io"
	"net/http"

	"github.com/savaki/swaggering"
)

func echo(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

type Owner struct {
	Name string
}

type Pet struct {
	Owner Owner
}

func main() {
	api := &swaggering.Api{
		BasePath: "/api",
		CORS:     true,
		Endpoints: []swaggering.Endpoint{
			{
				Method:      "get",
				Path:        "/pet",
				Summary:     "Add a New Pet",
				Description: "PetDescription",
				HandlerFunc: echo,
				Parameter: &swaggering.Parameter{
					Description: "Thingie!",
					Schema:      Owner{},
				},
				Responses: map[int]swaggering.Response{
					http.StatusOK: {
						Description: "Woo hoo!",
						Schema:      Pet{},
					},
				},
			},
		},
	}

	http.Handle("/swagger", api)
	http.ListenAndServe(":8080", nil)
}

//
//type Pet struct {
//	String string
//	Int    int64 `required:"true"`
//	Owner  Owner
//	Owners []Owner
//	Ptr    *Owner
//	Ptrs   []*Owner
//}
//
//
//func main() {
//	api := swaggering.NewApi().
//		WithBasePath("/api")
//
//	api.WithEndpoint("post", "/pet", echo,
//		swaggering.Summary("Add a new pet to the store"),
//	)
//	api.WithEndpoint("get", "/pet/findByStatus", echo,
//		swaggering.Summary("Finds Pets by status"),
//		swaggering.Description("Multiple status values can be provided with comma separated strings"),
//		swaggering.Param(
//			swaggering.ParamDescription("adding the thing to the thing"),
//			swaggering.ParamType(Owner{}),
//		),
//		swaggering.Response2(http.StatusOK, "successful operation", swaggering.ResponseType(Pet{})),
//		swaggering.Response2(http.StatusBadRequest, "Invalid status value"),
//	)
//
//	api.Walk(func(path string, endpoints *swaggering.Endpoints) {
//		http.Handle(path, endpoints)
//	})
//
//	// render swagger with cors
//	http.HandleFunc("/swagger", func(w http.ResponseWriter, req *http.Request) {
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		api.ServeHTTP(w, req)
//	})
//
//	http.ListenAndServe(":8080", nil)
//}

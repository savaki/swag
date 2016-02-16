package swaggering

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
)

type Parameter struct {
	In          string
	Name        string
	Description string
	Required    bool
	Schema      interface{}
}

type Response struct {
	Description string
	Schema      interface{}
}

type Endpoint struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
	Produces    []string
	Consumes    []string
	Parameter   *Parameter
	Responses   map[int]Response

	// Value is a container for arbitrary content to provide support for non net/http web frameworks
	// like gin
	Func interface{}
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if e.Handler != nil {
		e.Handler.ServeHTTP(w, req)
	} else if e.HandlerFunc != nil {
		e.HandlerFunc(w, req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type Api struct {
	Description    string
	Version        string
	TermsOfService string
	Title          string
	Email          string
	LicenseName    string
	LicenseUrl     string
	BasePath       string
	Schemes        []string
	Endpoints      []Endpoint

	// CORS indicates whether the swagger api should generate CORS * headers
	CORS    bool
	once    sync.Once
	swagger *SwaggerApi
}

func (api *Api) Walk(callback func(path string, endpoints *SwaggerEndpoints)) {
	api.init()

	for path, endpoints := range api.swagger.Paths {
		callback(filepath.Join(api.BasePath, path), endpoints)
	}
}

func (api *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	api.init()

	if api.CORS {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.swagger)
}

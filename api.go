package swaggering

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Endpoints represents a container for http method handlers
type Endpoints struct {
	Get    *Endpoint `json:"get,omitempty"`
	Post   *Endpoint `json:"post,omitempty"`
	Put    *Endpoint `json:"put,omitempty"`
	Delete *Endpoint `json:"delete,omitempty"`
	Head   *Endpoint `json:"head,omitempty"`
	Option *Endpoint `json:"option,omitempty"`
}

// ServeHTTP provides a default http.Handler implementation for convenience
func (e *Endpoints) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var endpoint *Endpoint

	if method := req.Method; strings.EqualFold(method, "get") {
		endpoint = e.Get
	} else if strings.EqualFold(method, "post") {
		endpoint = e.Post
	} else if strings.EqualFold(method, "put") {
		endpoint = e.Put
	} else if strings.EqualFold(method, "delete") {
		endpoint = e.Delete
	} else if strings.EqualFold(method, "head") {
		endpoint = e.Head
	} else if strings.EqualFold(method, "option") {
		endpoint = e.Option
	}

	if endpoint == nil || endpoint.Handler == nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		endpoint.Handler.ServeHTTP(w, req)
	}
}

// Api represents our swagger api
type Api struct {
	Swagger string `json:"swagger,omitempty"`
	Info    struct {
		Description    string `json:"description,omitempty"`
		Version        string `json:"version,omitempty"`
		TermsOfService string `json:"termsOfService,omitempty"`
		Title          string `json:"title,omitempty"`
		Contact        struct {
			Email string `json:"email,omitempty"`
		} `json:"contact"`
		License struct {
			Name string `json:"name,omitempty"`
			Url  string `json:"url,omitempty"`
		} `json:"license"`
	} `json:"info"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]*Endpoints `json:"paths,omitempty"`
	Definitions map[string]Definition `json:"definitions,omitempty"`
}

func (api *Api) Add(endpoint *Endpoint) *Api {
	if api.Paths == nil {
		api.Paths = map[string]*Endpoints{}
	}

	if api.Paths[endpoint.Path] == nil {
		api.Paths[endpoint.Path] = &Endpoints{}
	}

	endpoints := api.Paths[endpoint.Path]
	switch strings.ToLower(endpoint.Method) {
	case "get":
		endpoints.Get = endpoint
	case "post":
		endpoints.Post = endpoint
	case "put":
		endpoints.Put = endpoint
	case "delete":
		endpoints.Delete = endpoint
	case "option":
		endpoints.Option = endpoint
	case "head":
		endpoints.Head = endpoint
	}

	return api
}

func (api *Api) Endpoint(method, path string, handler http.Handler, options ...EndpointOption) *Api {
	return api.Add(NewEndpoint(method, path, handler, options...))
}

func (api *Api) EndpointFunc(method, path string, handlerFunc http.HandlerFunc, options ...EndpointOption) *Api {
	return api.Endpoint(method, path, handlerFunc, options...)
}

func (api *Api) Walk(walkFunc func(path string, endpoints *Endpoints)) {
	if api.Paths != nil {
		for path, endpoints := range api.Paths {
			walkFunc(path, endpoints)
		}
	}
}

func (api *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api)
}

type ApiOption func(*Api)

func ApiDescription(v string) ApiOption {
	return func(api *Api) {
		api.Info.Description = v
	}
}

func ApiVersion(v string) ApiOption {
	return func(api *Api) {
		api.Info.Version = v
	}
}

func ApiTitle(v string) ApiOption {
	return func(api *Api) {
		api.Info.Title = v
	}
}

func ApiEmail(v string) ApiOption {
	return func(api *Api) {
		api.Info.Contact.Email = v
	}
}

func ApiBasePath(v string) ApiOption {
	return func(api *Api) {
		api.BasePath = v
	}
}

// NewApi creates a new api instances using default parameters.  Additional parameters can be
// configured using ApiOption
func NewApi(options ...ApiOption) (*Api, error) {
	api := &Api{
		Swagger:  "2.0",
		BasePath: "/",
		Schemes:  []string{"http"},
	}

	api.Info.Contact.Email = "your-email-address"
	api.Info.Description = "Describe your API"
	api.Info.Title = "Your API Title"
	api.Info.Version = "SNAPSHOT"
	api.Info.TermsOfService = "http://swagger.io/terms/"
	api.Info.License.Name = "Apache 2.0"
	api.Info.License.Url = "http://www.apache.org/licenses/LICENSE-2.0.html"

	for _, option := range options {
		option(api)
	}

	return api, nil
}

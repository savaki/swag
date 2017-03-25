package swaggering

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SwaggerItems struct {
	Ref string `json:"$ref,omitempty"`
}

type SwaggerSchema struct {
	Type  string        `json:"type,omitempty"`
	Items *SwaggerItems `json:"items,omitempty"`
	Ref   string        `json:"$ref,omitempty"`
}

type SwaggerParameter struct {
	In          string         `json:"in,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required"`
	Schema      *SwaggerSchema `json:"schema,omitempty"`
	Type        string         `json:"type,omitempty"`
	Format      string         `json:"format,omitempty"`
}

type SwaggerResponse struct {
	Description string         `json:"description,omitempty"`
	Schema      *SwaggerSchema `json:"schema,omitempty"`
}

type SwaggerEndpoint struct {
	Tags        []string                   `json:"tags"`
	Path        string                     `json:"-"`
	Method      string                     `json:"-"`
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	OperationId string                     `json:"operationId"`
	Produces    []string                   `json:"produces,omitempty"`
	Consumes    []string                   `json:"consumes,omitempty"`
	HandlerFunc http.HandlerFunc           `json:"-"`
	Parameters  []SwaggerParameter         `json:"parameters,omitempty"`
	Responses   map[string]SwaggerResponse `json:"responses,omitempty"`
}

func (s *SwaggerEndpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.HandlerFunc(w, req)
}

type SwaggerEndpoints struct {
	Get    *SwaggerEndpoint `json:"get,omitempty"`
	Post   *SwaggerEndpoint `json:"post,omitempty"`
	Put    *SwaggerEndpoint `json:"put,omitempty"`
	Delete *SwaggerEndpoint `json:"delete,omitempty"`
	Head   *SwaggerEndpoint `json:"head,omitempty"`
	Option *SwaggerEndpoint `json:"option,omitempty"`
}

func (e *SwaggerEndpoints) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler

	if method := req.Method; strings.EqualFold(method, "get") {
		handler = e.Get
	} else if strings.EqualFold(method, "post") {
		handler = e.Post
	} else if strings.EqualFold(method, "put") {
		handler = e.Put
	} else if strings.EqualFold(method, "delete") {
		handler = e.Delete
	} else if strings.EqualFold(method, "head") {
		handler = e.Head
	} else if strings.EqualFold(method, "option") {
		handler = e.Option
	}

	if handler == nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		handler.ServeHTTP(w, req)
	}
}

type SwaggerApi struct {
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
	BasePath    string                       `json:"basePath,omitempty"`
	Schemes     []string                     `json:"schemes,omitempty"`
	Paths       map[string]*SwaggerEndpoints `json:"paths,omitempty"`
	Definitions map[string]Object            `json:"definitions,omitempty"`
	Tags        []Tag                        `json:"tags"`
	Host        string                       `json:"host"`
}

func (s *SwaggerApi) clone() *SwaggerApi {
	data, _ := json.Marshal(s)

	clone := &SwaggerApi{}
	json.Unmarshal(data, &clone)
	return clone
}

func (s *SwaggerApi) addEndpoint(endpoint *SwaggerEndpoint) {
	if s.Paths == nil {
		s.Paths = map[string]*SwaggerEndpoints{}
	}

	if s.Paths[endpoint.Path] == nil {
		s.Paths[endpoint.Path] = &SwaggerEndpoints{}
	}

	endpoints := s.Paths[endpoint.Path]
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
}

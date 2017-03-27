package swaggering

import (
	"encoding/json"
	"strings"
)

// Items
type Items struct {
	Ref string `json:"$ref,omitempty"`
}

type Schema struct {
	Type  string `json:"type,omitempty"`
	Items *Items `json:"items,omitempty"`
	Ref   string `json:"$ref,omitempty"`
}

type Parameter struct {
	In          string  `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required"`
	Schema      *Schema `json:"schema,omitempty"`
	Type        string  `json:"type,omitempty"`
	Format      string  `json:"format,omitempty"`
}

type Response struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Endpoint struct {
	Tags        []string            `json:"tags"`
	Path        string              `json:"-"`
	Method      string              `json:"-"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationId string              `json:"operationId"`
	Produces    []string            `json:"produces,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Handler     interface{}         `json:"-"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Endpoints struct {
	Get    *Endpoint `json:"get,omitempty"`
	Post   *Endpoint `json:"post,omitempty"`
	Put    *Endpoint `json:"put,omitempty"`
	Delete *Endpoint `json:"delete,omitempty"`
	Head   *Endpoint `json:"head,omitempty"`
	Option *Endpoint `json:"option,omitempty"`
}

type Contact struct {
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type Info struct {
	Description    string  `json:"description,omitempty"`
	Version        string  `json:"version,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Title          string  `json:"title,omitempty"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
}

type Api struct {
	Swagger     string                `json:"swagger,omitempty"`
	Info        Info                  `json:"info"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]*Endpoints `json:"paths,omitempty"`
	Definitions map[string]object     `json:"definitions,omitempty"`
	Tags        []OldTag              `json:"tags"`
	Host        string                `json:"host"`
}

func (s *Api) clone() *Api {
	data, _ := json.Marshal(s)

	clone := &Api{}
	json.Unmarshal(data, &clone)
	return clone
}

func (s *Api) addEndpoint(endpoint *Endpoint) {
	if s.Paths == nil {
		s.Paths = map[string]*Endpoints{}
	}

	if s.Paths[endpoint.Path] == nil {
		s.Paths[endpoint.Path] = &Endpoints{}
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

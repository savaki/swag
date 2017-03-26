package swaggering

import (
	"net/http"
	"strings"
)

type EndpointBuilder struct {
	Endpoint Endpoint
}

func (b *EndpointBuilder) Handler(h http.Handler) *EndpointBuilder {
	b.Endpoint.Handler = h
	return b
}

func (b *EndpointBuilder) Func(v interface{}) *EndpointBuilder {
	b.Endpoint.Func = v
	return b
}

func (b *EndpointBuilder) Summary(v string) *EndpointBuilder {
	b.Endpoint.Summary = v
	return b
}

func (b *EndpointBuilder) Description(v string) *EndpointBuilder {
	b.Endpoint.Description = v
	return b
}

func (b *EndpointBuilder) Parameter(p Parameter) *EndpointBuilder {
	if b.Endpoint.Parameters == nil {
		b.Endpoint.Parameters = []Parameter{}
	}

	b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	return b
}

func (b *EndpointBuilder) Path(name, typ, description string, required bool) *EndpointBuilder {
	p := Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *EndpointBuilder) Query(name, typ, description string, required bool) *EndpointBuilder {
	p := Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *EndpointBuilder) Body(schema interface{}, description string, required bool) *EndpointBuilder {
	p := Parameter{
		Description: description,
		Schema:      schema,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *EndpointBuilder) Tags(tags ...string) *EndpointBuilder {
	if b.Endpoint.Tags == nil {
		b.Endpoint.Tags = []string{}
	}

	b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	return b
}

func (b *EndpointBuilder) Response(code int, schema interface{}, description string) *EndpointBuilder {
	if b.Endpoint.Responses == nil {
		b.Endpoint.Responses = map[int]Response{}
	}

	b.Endpoint.Responses[code] = Response{
		Description: description,
		Schema:      schema,
	}

	return b
}

func NewEndpoint(method, path string, handler interface{}) *EndpointBuilder {
	method = strings.ToUpper(method)
	return &EndpointBuilder{
		Endpoint: Endpoint{
			Method: method,
			Path:   path,
			Func:   handler,
		},
	}
}

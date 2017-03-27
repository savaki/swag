package swaggering

import (
	"strconv"
	"strings"
)

type EndpointBuilder struct {
	Endpoint Endpoint
}

type EndpointOption func(builder *EndpointBuilder)

func EndpointDescription(v string) EndpointOption {
	return func(b *EndpointBuilder) {
		b.Endpoint.Description = v
	}
}

func OperationId(v string) EndpointOption {
	return func(b *EndpointBuilder) {
		b.Endpoint.OperationId = v
	}
}

func Produces(v ...string) EndpointOption {
	return func(b *EndpointBuilder) {
		b.Endpoint.Produces = v
	}
}

func Consumes(v ...string) EndpointOption {
	return func(b *EndpointBuilder) {
		b.Endpoint.Consumes = v
	}
}

func parameter(p Parameter) EndpointOption {
	return func(b *EndpointBuilder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

func (b *EndpointBuilder) Path(name, typ, description string, required bool) EndpointOption {
	p := Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func (b *EndpointBuilder) Query(name, typ, description string, required bool) EndpointOption {
	p := Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func (b *EndpointBuilder) Body(prototype interface{}, description string, required bool) EndpointOption {
	p := Parameter{
		Description: description,
		Schema:      makeSchema(prototype),
		Required:    required,
	}
	return parameter(p)
}

func Tags(tags ...string) EndpointOption {
	return func(b *EndpointBuilder) {
		if b.Endpoint.Tags == nil {
			b.Endpoint.Tags = []string{}
		}

		b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	}
}

func Respond(code int, prototype interface{}, description string) EndpointOption {
	return func(b *EndpointBuilder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = map[string]Response{}
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = Response{
			Description: description,
			Schema:      makeSchema(prototype),
		}
	}
}

func NewEndpoint(method, path, summary string, handler interface{}, options ...EndpointOption) *EndpointBuilder {
	method = strings.ToUpper(method)
	e := &EndpointBuilder{
		Endpoint: Endpoint{
			Method:   method,
			Path:     path,
			Summary:  summary,
			Handler:  handler,
			Produces: []string{"application/json"},
			Consumes: []string{"application/json"},
		},
	}

	for _, opt := range options {
		opt(e)
	}

	return e
}

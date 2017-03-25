package swaggering

import (
	"net/http"
	"strings"
)

type Builder struct {
	Endpoint Endpoint
}

func (b *Builder) Handler(h http.Handler) *Builder {
	b.Endpoint.Handler = h
	return b
}

func (b *Builder) Func(v interface{}) *Builder {
	b.Endpoint.Func = v
	return b
}

func (b *Builder) Summary(v string) *Builder {
	b.Endpoint.Summary = v
	return b
}

func (b *Builder) Description(v string) *Builder {
	b.Endpoint.Description = v
	return b
}

func (b *Builder) Parameter(p Parameter) *Builder {
	if b.Endpoint.Parameters == nil {
		b.Endpoint.Parameters = []Parameter{}
	}

	b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	return b
}

func (b *Builder) Path(name, typ, description string, required bool) *Builder {
	p := Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *Builder) Query(name, typ, description string, required bool) *Builder {
	p := Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *Builder) Schema(schema interface{}, description string, required bool) *Builder {
	p := Parameter{
		Description: description,
		Schema:      schema,
		Required:    required,
	}
	return b.Parameter(p)
}

func (b *Builder) Tags(tags ...string) *Builder {
	if b.Endpoint.Tags == nil {
		b.Endpoint.Tags = []string{}
	}

	b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	return b
}

func (b *Builder) Response(code int, schema interface{}, description string) *Builder {
	if b.Endpoint.Responses == nil {
		b.Endpoint.Responses = map[int]Response{}
	}

	b.Endpoint.Responses[code] = Response{
		Description: description,
		Schema:      schema,
	}

	return b
}

func New(method, path string, handler interface{}) *Builder {
	method = strings.ToUpper(method)
	return &Builder{
		Endpoint: Endpoint{
			Method: method,
			Path:   path,
			Func:   handler,
		},
	}
}

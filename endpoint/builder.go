package endpoint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/savaki/swag/swagger"
)

type Builder struct {
	Endpoint *swagger.Endpoint
}

func (b *Builder) Build() *swagger.Endpoint {
	return b.Endpoint
}

type Option func(builder *Builder)

func Summary(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Summary = v
	}
}

func Description(v string) Option {
	return func(b *Builder) {
		b.Endpoint.Description = v
	}
}

func OperationId(v string) Option {
	return func(b *Builder) {
		b.Endpoint.OperationId = v
	}
}

func Produces(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Produces = v
	}
}

func Consumes(v ...string) Option {
	return func(b *Builder) {
		b.Endpoint.Consumes = v
	}
}

func parameter(p swagger.Parameter) Option {
	return func(b *Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []swagger.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

func Path(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func Query(name, typ, description string, required bool) Option {
	p := swagger.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func Body(prototype interface{}, description string, required bool) Option {
	p := swagger.Parameter{
		In:          "body",
		Description: description,
		Schema:      swagger.MakeSchema(prototype),
		Required:    required,
	}
	return parameter(p)
}

func Tags(tags ...string) Option {
	return func(b *Builder) {
		if b.Endpoint.Tags == nil {
			b.Endpoint.Tags = []string{}
		}

		b.Endpoint.Tags = append(b.Endpoint.Tags, tags...)
	}
}

func Response(code int, prototype interface{}, description string) Option {
	return func(b *Builder) {
		if b.Endpoint.Responses == nil {
			b.Endpoint.Responses = map[string]swagger.Response{}
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = swagger.Response{
			Description: description,
			Schema:      swagger.MakeSchema(prototype),
		}
	}
}

func New(method, path string, handler interface{}, options ...Option) *Builder {
	if v, ok := handler.(func(w http.ResponseWriter, r *http.Request)); ok {
		handler = http.HandlerFunc(v)
	}

	method = strings.ToUpper(method)
	e := &Builder{
		Endpoint: &swagger.Endpoint{
			Method:   method,
			Path:     path,
			Handler:  handler,
			Produces: []string{"application/json"},
			Consumes: []string{"application/json"},
			Tags:     []string{},
		},
	}

	for _, opt := range options {
		opt(e)
	}

	return e
}

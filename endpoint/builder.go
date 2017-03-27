package endpoint

import (
	"strconv"
	"strings"

	"github.com/savaki/swaggering/types"
)

type Builder struct {
	Endpoint *types.Endpoint
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

func parameter(p types.Parameter) Option {
	return func(b *Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []types.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}

func Path(name, typ, description string, required bool) Option {
	p := types.Parameter{
		Name:        name,
		In:          "path",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func Query(name, typ, description string, required bool) Option {
	p := types.Parameter{
		Name:        name,
		In:          "query",
		Type:        typ,
		Description: description,
		Required:    required,
	}
	return parameter(p)
}

func Body(prototype interface{}, description string, required bool) Option {
	p := types.Parameter{
		Description: description,
		Schema:      makeSchema(prototype),
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
			b.Endpoint.Responses = map[string]types.Response{}
		}

		b.Endpoint.Responses[strconv.Itoa(code)] = types.Response{
			Description: description,
			Schema:      makeSchema(prototype),
		}
	}
}

func New(method, path string, handler interface{}, options ...Option) *Builder {
	method = strings.ToUpper(method)
	e := &Builder{
		Endpoint: &types.Endpoint{
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

package swaggering

import (
	"fmt"
	"net/http"
	"strings"
)

type Parameter struct {
	In          string `json:"in,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Schema      struct {
		Ref string `json:"$ref,omitempty"`
	} `json:"schema,omitempty"`
}

type SchemaItems struct {
	Ref string `json:"$ref"`
}

type Schema struct {
	Ref   string       `json:"$ref,omitempty"`
	Type  string       `json:"type,omitempty"`
	Items *SchemaItems `json:"items,omitempty"`
}

type ResponseItem struct {
	Code        string  `json:"-"`
	Description string  `json:"description"`
	Schema      *Schema `json:"schema,omitempty"`
}

// Endpoint represents the handler for a specific path and method
type Endpoint struct {
	Method      string                  `json:"-"`
	Path        string                  `json:"-"`
	Handler     http.Handler            `json:"-"`
	Summary     string                  `json:"summary,omitempty"`
	Description string                  `json:"description,omitempty"`
	OperationId string                  `json:"operationId"`
	Produces    []string                `json:"produces,omitempty"`
	Consumes    []string                `json:"consumes,omitempty"`
	Parameters  []*Parameter            `json:"parameters"`
	Responses   map[string]ResponseItem `json:"responses"`
}

// EndpointOption provides options to configure an endpoint
type EndpointOption func(*Api, *Endpoint)

// NewEndpoint constructs a new endpoint for the specified method and path.  method is case insensitive
func (api *Api) newEndpoint(method, path string, handler http.Handler, options ...EndpointOption) *Endpoint {
	endpoint := &Endpoint{
		Method:      strings.ToUpper(method),
		Path:        path,
		Handler:     handler,
		Produces:    []string{"application/json"},
		Consumes:    []string{"application/json"},
		Parameters:  []*Parameter{},
		Responses:   map[string]ResponseItem{},
		OperationId: method + path,
	}

	for _, option := range options {
		option(api, endpoint)
	}

	return endpoint
}

// Summary configures an endpoints summary field
func Summary(v string) EndpointOption {
	return func(api *Api, endpoint *Endpoint) {
		endpoint.Summary = v
	}
}

// Description configures an endpoints description field
func Description(v string) EndpointOption {
	return func(api *Api, endpoint *Endpoint) {
		endpoint.Description = v
	}
}

type ParamOption func(api *Api, parameter *Parameter)

func ParamDescription(v string) ParamOption {
	return func(api *Api, p *Parameter) {
		p.Description = v
	}
}

func ParamType(v interface{}) ParamOption {
	return func(api *Api, p *Parameter) {
		obj := defineObject(v)
		p.Schema.Ref = makeRef(obj.Name)

		for _, definition := range define(v) {
			api.AddDefinition(definition)
		}
	}
}

func Param(options ...ParamOption) EndpointOption {
	return func(api *Api, endpoint *Endpoint) {
		if endpoint.Parameters == nil {
			endpoint.Parameters = []*Parameter{}
		}

		parameter := &Parameter{
			In:       "body",
			Name:     "body",
			Required: true,
		}

		for _, option := range options {
			option(api, parameter)
		}

		endpoint.Parameters = append(endpoint.Parameters, parameter)
	}
}

type TypeOption func(*Api, *ResponseItem)

func Response(code int, description string, options ...TypeOption) EndpointOption {
	return func(api *Api, endpoint *Endpoint) {
		if endpoint.Responses == nil {
			endpoint.Responses = map[string]ResponseItem{}
		}

		response := ResponseItem{
			Code:        fmt.Sprintf("%v", code),
			Description: description,
		}

		for _, option := range options {
			option(api, &response)
		}

		endpoint.Responses[response.Code] = response
	}
}

func Type(v interface{}) TypeOption {
	return func(api *Api, response *ResponseItem) {
		if v == nil {
			return
		}

		obj := defineObject(v)

		var schema *Schema
		if obj.IsArray {
			schema = &Schema{
				Type: "array",
				Items: &SchemaItems{
					Ref: makeRef(obj.GoType.Name()),
				},
			}
		} else {
			schema = &Schema{
				Ref: makeRef(obj.GoType.Name()),
			}
		}

		for _, definition := range define(v) {
			api.AddDefinition(definition)
		}

		response.Schema = schema
	}
}

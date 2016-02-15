package swaggering

import (
	"fmt"
	"net/http"
	"strings"
)

type Parameter struct {
}

type Response struct {
	Code        string `json:"-"`
	Description string `json:"description"`
}

// Endpoint represents the handler for a specific path and method
type Endpoint struct {
	Method      string              `json:"-"`
	Path        string              `json:"-"`
	Handler     http.Handler        `json:"-"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Parameters  []Parameter         `json:"parameters"`
	Responses   map[string]Response `json:"responses"`
}

// EndpointOption provides options to configure an endpoint
type EndpointOption func(*Endpoint)

// NewEndpoint constructs a new endpoint for the specified method and path.  method is case insensitive
func NewEndpoint(method, path string, handler http.Handler, options ...EndpointOption) *Endpoint {
	endpoint := &Endpoint{
		Method:      strings.ToUpper(method),
		Path:        path,
		Handler:     handler,
		Produces:    []string{"application/json"},
		Consumes:    []string{"application/json"},
		Parameters:  []Parameter{},
		Responses:   map[string]Response{},
	}

	for _, option := range options {
		option(endpoint)
	}

	return endpoint
}

// Summary configures an endpoints summary field
func Summary(v string) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Summary = v
	}
}

// Description configures an endpoints description field
func Description(v string) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Description = v
	}
}

func Param(kind interface{}) EndpointOption {
	return func(endpoint *Endpoint) {
	}
}

func RespondsWith(code int, description string, kind interface{}) EndpointOption {
	return func(endpoint *Endpoint) {
		if endpoint.Responses == nil {
			endpoint.Responses = map[string]Response{}
		}

		response := Response{
			Code:        fmt.Sprintf("%v", code),
			Description: description,
		}
		endpoint.Responses[response.Code] = response
	}
}

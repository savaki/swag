package swaggering

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func (api *OldApi) initSchema(v interface{}) *Schema {
	schema := &Schema{}

	obj := defineObject(v)
	if obj.IsArray {
		schema.Type = "array"
		schema.Items = &Items{
			Ref: makeRef(obj.Name),
		}

	} else {
		schema.Type = "object"
		schema.Ref = makeRef(obj.Name)
	}

	return schema
}

func (api *OldApi) initParameter(parameter OldParameter) Parameter {
	p := Parameter{
		In:          parameter.In,
		Name:        parameter.Name,
		Description: parameter.Description,
		Required:    parameter.Required,
		Type:        parameter.Type,
		Format:      parameter.Format,
	}

	if parameter.Schema != nil {
		property := inspect(reflect.TypeOf(parameter.Schema), "")

		switch property.GoType.Kind() {
		case reflect.Struct:
			p.Schema = api.initSchema(parameter.Schema)
			if p.In == "" {
				p.In = "body"
			}
			if p.Name == "" {
				p.Name = "body"
			}
		}
	}

	return p
}

func (api *OldApi) initResponse(code int, response OldResponse) (string, Response) {
	sc := fmt.Sprintf("%v", code)
	r := Response{
		Description: response.Description,
	}

	if response.Schema != nil {
		r.Schema = api.initSchema(response.Schema)
	}

	return sc, r
}

func (api *OldApi) initEndpoint(endpoint OldEndpoint) *Endpoint {
	se := &Endpoint{
		Tags:        endpoint.Tags,
		Path:        endpoint.Path,
		Method:      strings.ToLower(endpoint.Method),
		Summary:     endpoint.Summary,
		Description: endpoint.Description,
		OperationId: endpoint.Method + endpoint.Path,
		Produces:    endpoint.Produces,
		Consumes:    endpoint.Consumes,
	}

	// set default produces/consumes
	if se.Produces == nil {
		se.Produces = []string{"application/json"}
	}
	if se.Consumes == nil {
		se.Consumes = []string{"application/json"}
	}

	// assign the handler
	if endpoint.Handler != nil {
		se.Handler = endpoint.Handler.ServeHTTP
	} else {
		se.Handler = endpoint.HandlerFunc
	}

	// handle the parameters
	if endpoint.Parameters != nil {
		se.Parameters = []Parameter{}
		for _, p := range endpoint.Parameters {
			se.Parameters = append(se.Parameters, api.initParameter(p))
		}
	}

	// handle the responses
	se.Responses = map[string]Response{}
	for code, response := range endpoint.Responses {
		sc, r := api.initResponse(code, response)
		se.Responses[sc] = r
	}

	return se
}

func (api *OldApi) initDefinitions() map[string]object {
	definitions := map[string]object{}
	objects := []interface{}{}

	// collect all the objects from all the endpoints
	for _, endpoint := range api.Endpoints {
		if endpoint.Parameters != nil {
			for _, p := range endpoint.Parameters {
				if p.Schema != nil {
					objects = append(objects, p.Schema)
				}
			}
		}

		if endpoint.Responses != nil {
			for _, response := range endpoint.Responses {
				if response.Schema != nil {
					objects = append(objects, response.Schema)
				}
			}
		}
	}

	for _, v := range objects {
		for name, obj := range define(v) {
			if _, exists := definitions[name]; !exists {
				definitions[name] = obj
			}
		}
	}

	return definitions
}

func (api *OldApi) initOnce() {
	// initialize properties that haven't been set
	//
	api.mux = &sync.Mutex{}
	api.byHostAndScheme = map[string]*Api{}

	// render the input into the swagger model
	//
	//api.template = api.initSwagger()

	for _, endpoint := range api.Endpoints {
		se := api.initEndpoint(endpoint)
		api.template.addEndpoint(se)
	}

	api.template.Definitions = api.initDefinitions()
}

func (api *OldApi) init() {
	api.once.Do(api.initOnce)
}

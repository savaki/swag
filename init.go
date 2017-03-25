package swaggering

import (
	"fmt"
	"reflect"
	"strings"
)

func (api *Api) initSwagger() *SwaggerApi {
	s := &SwaggerApi{
		BasePath: "/",
		Swagger:  "2.0",
		Schemes:  []string{"http"},
		Tags:     api.Tags,
		Host:     api.Host,
	}

	s.Info.Contact.Email = "your-email-address"
	s.Info.Description = "Describe your API"
	s.Info.Title = "Your API Title"
	s.Info.Version = "SNAPSHOT"
	s.Info.TermsOfService = "http://swagger.io/terms/"
	s.Info.License.Name = "Apache 2.0"
	s.Info.License.Url = "http://www.apache.org/licenses/LICENSE-2.0.html"

	// override with user provided values

	if api.Description != "" {
		s.Info.Description = api.Description
	}
	if api.Version != "" {
		s.Info.Version = api.Version
	}
	if api.TermsOfService != "" {
		s.Info.TermsOfService = api.TermsOfService
	}
	if api.Title != "" {
		s.Info.Title = api.Title
	}
	if api.Email != "" {
		s.Info.Contact.Email = api.Email
	}
	if api.LicenseName != "" {
		s.Info.License.Name = api.LicenseName
	}
	if api.LicenseUrl != "" {
		s.Info.License.Url = api.LicenseUrl
	}
	if api.BasePath != "" {
		s.BasePath = api.BasePath
	}
	if api.Schemes != nil {
		s.Schemes = api.Schemes
	}

	return s
}

func (api *Api) initSchema(v interface{}) *SwaggerSchema {
	schema := &SwaggerSchema{}

	obj := defineObject(v)
	if obj.IsArray {
		schema.Type = "array"
		schema.Items = &SwaggerItems{
			Ref: makeRef(obj.Name),
		}

	} else {
		schema.Type = "object"
		schema.Ref = makeRef(obj.Name)
	}

	return schema
}

func (api *Api) initParameter(parameter Parameter) SwaggerParameter {
	p := SwaggerParameter{
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

func (api *Api) initResponse(code int, response Response) (string, SwaggerResponse) {
	sc := fmt.Sprintf("%v", code)
	r := SwaggerResponse{
		Description: response.Description,
	}

	if response.Schema != nil {
		r.Schema = api.initSchema(response.Schema)
	}

	return sc, r
}

func (api *Api) initEndpoint(endpoint Endpoint) *SwaggerEndpoint {
	se := &SwaggerEndpoint{
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
		se.HandlerFunc = endpoint.Handler.ServeHTTP
	} else {
		se.HandlerFunc = endpoint.HandlerFunc
	}

	// handle the parameters
	if endpoint.Parameters != nil {
		se.Parameters = []SwaggerParameter{}
		for _, p := range endpoint.Parameters {
			se.Parameters = append(se.Parameters, api.initParameter(p))
		}
	}

	// handle the responses
	se.Responses = map[string]SwaggerResponse{}
	for code, response := range endpoint.Responses {
		sc, r := api.initResponse(code, response)
		se.Responses[sc] = r
	}

	return se
}

func (api *Api) initDefinitions() map[string]Object {
	definitions := map[string]Object{}
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

func (api *Api) initOnce() {
	api.swagger = api.initSwagger()

	for _, endpoint := range api.Endpoints {
		se := api.initEndpoint(endpoint)
		api.swagger.addEndpoint(se)
	}

	api.swagger.Definitions = api.initDefinitions()
}

func (api *Api) init() {
	api.once.Do(api.initOnce)
}

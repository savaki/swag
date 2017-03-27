package swagger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

// Object represents the object entity from the swagger definition
type Object struct {
	IsArray    bool                `json:"-"`
	GoType     reflect.Type        `json:"-"`
	Name       string              `json:"-"`
	Type       string              `json:"type"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// Property represents the property entity from the swagger definition
type Property struct {
	GoType      reflect.Type `json:"-"`
	Type        string       `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
	Enum        []string     `json:"enum,omitempty"`
	Format      string       `json:"format,omitempty"`
	Ref         string       `json:"$ref,omitempty"`
	Example     string       `json:"example,omitempty"`
	Items       *Items       `json:"items,omitempty"`
}

// Contact represents the contact entity from the swagger definition; used by Info
type Contact struct {
	Email string `json:"email,omitempty"`
}

// License represents the license entity from the swagger definition; used by Info
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Info represents the info entity from the swagger definition
type Info struct {
	Description    string  `json:"description,omitempty"`
	Version        string  `json:"version,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Title          string  `json:"title,omitempty"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
}

// Endpoints represents all the swagger endpoints associated with a particular path
type Endpoints struct {
	Delete  *Endpoint `json:"delete,omitempty"`
	Head    *Endpoint `json:"head,omitempty"`
	Get     *Endpoint `json:"get,omitempty"`
	Options *Endpoint `json:"options,omitempty"`
	Post    *Endpoint `json:"post,omitempty"`
	Put     *Endpoint `json:"put,omitempty"`
	Patch   *Endpoint `json:"patch,omitempty"`
	Trace   *Endpoint `json:"trace,omitempty"`
	Connect *Endpoint `json:"connect,omitempty"`
}

// ServeHTTP allows endpoints to serve itself using the builtin http mux
func (e *Endpoints) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var endpoint *Endpoint

	switch req.Method {
	case "DELETE":
		endpoint = e.Delete
	case "HEAD":
		endpoint = e.Head
	case "GET":
		endpoint = e.Get
	case "OPTIONS":
		endpoint = e.Options
	case "POST":
		endpoint = e.Post
	case "PUT":
		endpoint = e.Put
	case "PATCH":
		endpoint = e.Patch
	case "TRACE":
		endpoint = e.Trace
	case "CONNECT":
		endpoint = e.Connect
	}

	if endpoint == nil || endpoint.Handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch v:= endpoint.Handler.(type){
	case func( w http.ResponseWriter, req *http.Request) :
		v(w, req)
	case http.HandlerFunc:
		v(w, req)
	case http.Handler:
		v.ServeHTTP(w, req)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Handler is not a standard http handler")
	}
}

// Walk calls the specified function for each method defined within the Endpoints
func (e *Endpoints) Walk(fn func(endpoint *Endpoint)) {
	if e.Delete != nil {
		fn(e.Delete)
	}
	if e.Head != nil {
		fn(e.Head)
	}
	if e.Get != nil {
		fn(e.Get)
	}
	if e.Options != nil {
		fn(e.Options)
	}
	if e.Post != nil {
		fn(e.Post)
	}
	if e.Put != nil {
		fn(e.Put)
	}
	if e.Patch != nil {
		fn(e.Patch)
	}
	if e.Trace != nil {
		fn(e.Trace)
	}
	if e.Connect != nil {
		fn(e.Connect)
	}
}

// API provides the top level encapsulation for the swagger definition
type API struct {
	Swagger     string                `json:"swagger,omitempty"`
	Info        Info                  `json:"info"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]*Endpoints `json:"paths,omitempty"`
	Definitions map[string]Object     `json:"definitions,omitempty"`
	Tags        []Tag                 `json:"tags"`
	Host        string                `json:"host"`
}

func (a *API) clone() *API {
	return &API{
		Swagger:     a.Swagger,
		Info:        a.Info,
		BasePath:    a.BasePath,
		Schemes:     a.Schemes,
		Paths:       a.Paths,
		Definitions: a.Definitions,
		Tags:        a.Tags,
		Host:        a.Host,
	}
}

func (a *API) addPath(e *Endpoint) {
	if a.Paths == nil {
		a.Paths = map[string]*Endpoints{}
	}

	v, ok := a.Paths[e.Path]
	if !ok {
		v = &Endpoints{}
		a.Paths[e.Path] = v
	}

	switch strings.ToUpper(e.Method) {
	case "DELETE":
		v.Delete = e
	case "GET":
		v.Get = e
	case "HEAD":
		v.Head = e
	case "OPTIONS":
		v.Options = e
	case "POST":
		v.Post = e
	case "PUT":
		v.Put = e
	case "PATCH":
		v.Patch = e
	case "TRACE":
		v.Trace = e
	case "CONNECT":
		v.Connect = e
	default:
		panic(fmt.Errorf("invalid method, %v", e.Method))
	}
}

func (a *API) addDefinition(e *Endpoint) {
	if a.Definitions == nil {
		a.Definitions = map[string]Object{}
	}

	if e.Parameters != nil {
		for _, p := range e.Parameters {
			if p.Schema != nil {
				def := define(p.Schema.Prototype)
				for k, v := range def {
					if _, ok := a.Definitions[k]; !ok {
						a.Definitions[k] = v
					}
				}
			}
		}
	}

	if e.Responses != nil {
		for _, response := range e.Responses {
			if response.Schema != nil {
				def := define(response.Schema.Prototype)
				for k, v := range def {
					if _, ok := a.Definitions[k]; !ok {
						a.Definitions[k] = v
					}
				}
			}
		}
	}
}

// AddEndpoint adds the specified endpoint to the API definition; to generate an endpoint use ```endpoint.New```
func (a *API) AddEndpoint(e *Endpoint) {
	a.addPath(e)
	a.addDefinition(e)
}

// Handler is a factory method that generates an http.HandlerFunc; if enableCors is true, then the handler will generate
// cors headers
func (a *API) Handler(enableCors bool) http.HandlerFunc {
	mux := &sync.Mutex{}
	byHostAndScheme := map[string]*API{}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if enableCors {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.WriteHeader(http.StatusOK)

		// customize the swagger header based on host
		//
		scheme := req.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}

		hostAndScheme := req.Host + ":" + scheme
		mux.Lock()
		v, ok := byHostAndScheme[hostAndScheme]
		if !ok {
			v = a.clone()
			v.Host = req.Host
			v.Schemes = []string{scheme}
			byHostAndScheme[hostAndScheme] = v
		}
		mux.Unlock()

		json.NewEncoder(w).Encode(v)
	}
}

// Walk invoke the callback for each endpoints defined in the swagger doc
func (a *API) Walk(callback func(path string, endpoints *Endpoint)) {
	for path, endpoints := range a.Paths {
		u := filepath.Join(a.BasePath, path)
		endpoints.Walk(func(endpoint *Endpoint) {
			callback(u, endpoint)
		})
	}
}

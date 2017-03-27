package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type Object struct {
	IsArray    bool                `json:"-"`
	GoType     reflect.Type        `json:"-"`
	Name       string              `json:"-"`
	Type       string              `json:"type"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

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

type Contact struct {
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type Info struct {
	Description    string  `json:"description,omitempty"`
	Version        string  `json:"version,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Title          string  `json:"title,omitempty"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
}

type Endpoints struct {
	Delete *Endpoint `json:"delete,omitempty"`
	Head   *Endpoint `json:"head,omitempty"`
	Get    *Endpoint `json:"get,omitempty"`
	Option *Endpoint `json:"option,omitempty"`
	Post   *Endpoint `json:"post,omitempty"`
	Put    *Endpoint `json:"put,omitempty"`
}

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
	if e.Option != nil {
		fn(e.Option)
	}
	if e.Post != nil {
		fn(e.Post)
	}
	if e.Put != nil {
		fn(e.Put)
	}
}

type Api struct {
	Swagger     string                `json:"swagger,omitempty"`
	Info        Info                  `json:"info"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Paths       map[string]*Endpoints `json:"paths,omitempty"`
	Definitions map[string]Object     `json:"definitions,omitempty"`
	Tags        []Tag                 `json:"tags"`
	Host        string                `json:"host"`
}

func (a *Api) clone() *Api {
	return &Api{
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

func (a *Api) addPath(e *Endpoint) {
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
	case "OPTION":
		v.Option = e
	case "POST":
		v.Post = e
	case "PUT":
		v.Put = e
	default:
		panic(fmt.Errorf("invalid method, %v", e.Method))
	}
}

func (a *Api) addDefinition(e *Endpoint) {
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

func (a *Api) AddEndpoint(e *Endpoint) {
	a.addPath(e)
	a.addDefinition(e)
}

func (a *Api) Handler(cors bool) http.HandlerFunc {
	mux := &sync.Mutex{}
	byHostAndScheme := map[string]*Api{}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if cors {
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

func (a *Api) Walk(callback func(path string, endpoints *Endpoint)) {
	for path, endpoints := range a.Paths {
		u := filepath.Join(a.BasePath, path)
		endpoints.Walk(func(endpoint *Endpoint) {
			callback(u, endpoint)
		})
	}
}

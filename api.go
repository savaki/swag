package swaggering

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
)

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Docs        struct {
		Description string `json:"description"`
		Url         string `json:"url"`
	} `json:"externalDocs"`
}

type Parameter struct {
	In          string
	Name        string
	Description string
	Required    bool
	Schema      interface{}
	Type        string
	Format      string
}

type Response struct {
	Description string
	Schema      interface{}
}

type Endpoint struct {
	Tags        []string
	Method      string
	Path        string
	Summary     string
	Description string
	Handler     http.Handler     `json:"-"`
	HandlerFunc http.HandlerFunc `json:"-"`
	Produces    []string
	Consumes    []string
	Parameters  []Parameter
	Responses   map[int]Response

	// Value is a container for arbitrary content to provide support for non net/http web frameworks
	// like gin
	Func interface{} `json:"-"`
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if e.Handler != nil {
		e.Handler.ServeHTTP(w, req)
	} else if e.HandlerFunc != nil {
		e.HandlerFunc(w, req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type Endpoints []Endpoint

func (e Endpoints) Append(endpoints ...Endpoint) Endpoints {
	return append(e, endpoints...)
}

type Api struct {
	Description    string
	Version        string
	TermsOfService string
	Title          string
	Email          string
	LicenseName    string
	LicenseUrl     string
	BasePath       string
	Schemes        []string
	Endpoints      []Endpoint
	Tags           []Tag
	Host           string

	// CORS indicates whether the swagger api should generate CORS * headers
	CORS            bool
	once            sync.Once
	mux             *sync.Mutex
	byHostAndScheme map[string]*SwaggerApi
	template        *SwaggerApi
}

func (api *Api) Walk(callback func(path string, endpoints *SwaggerEndpoints)) {
	api.init()

	for path, endpoints := range api.template.Paths {
		callback(filepath.Join(api.BasePath, path), endpoints)
	}
}

func (api *Api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	api.init()

	if api.CORS {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// customize the swagger header based on host
	//
	scheme := req.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	hostAndScheme := req.Host + ":" + scheme
	api.mux.Lock()
	v, ok := api.byHostAndScheme[hostAndScheme]
	if !ok {
		v = api.template.clone()
		v.Host = req.Host
		v.Schemes = []string{scheme}
		api.byHostAndScheme[hostAndScheme] = v
	}
	api.mux.Unlock()

	json.NewEncoder(w).Encode(v)
}

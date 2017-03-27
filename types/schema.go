package types

type Items struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	Ref    string `json:"$ref,omitempty"`
}

type Schema struct {
	Type      string      `json:"type,omitempty"`
	Items     *Items      `json:"items,omitempty"`
	Ref       string      `json:"$ref,omitempty"`
	Prototype interface{} `json:"-"`
}

type Response struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Parameter struct {
	In          string  `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required"`
	Schema      *Schema `json:"schema,omitempty"`
	Type        string  `json:"type,omitempty"`
	Format      string  `json:"format,omitempty"`
}

type Endpoint struct {
	Tags        []string            `json:"tags"`
	Path        string              `json:"-"`
	Method      string              `json:"-"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationId string              `json:"operationId"`
	Produces    []string            `json:"produces,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Handler     interface{}         `json:"-"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

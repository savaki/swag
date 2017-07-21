// Copyright 2017 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package swagger

import "encoding/json"

// Items represents items from the swagger doc
type Items struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	Ref    string `json:"$ref,omitempty"`
}

// Schema represents a schema from the swagger doc
type Schema struct {
	Type      string      `json:"type,omitempty"`
	Items     *Items      `json:"items,omitempty"`
	Ref       string      `json:"$ref,omitempty"`
	Prototype interface{} `json:"-"`
}

// Header represents a response header
type Header struct {
	Type        string `json:"type"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

// Response represents a response from the swagger doc
type Response struct {
	Description string            `json:"description,omitempty"`
	Schema      *Schema           `json:"schema,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty"`
}

// Parameter represents a parameter from the swagger doc
type Parameter struct {
	In          string  `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required"`
	Schema      *Schema `json:"schema,omitempty"`
	Type        string  `json:"type,omitempty"`
	Format      string  `json:"format,omitempty"`
}

// Endpoint represents an endpoint from the swagger doc
type Endpoint struct {
	Tags        []string            `json:"tags"`
	Path        string              `json:"-"`
	Method      string              `json:"-"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Handler     interface{}         `json:"-"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`

	// swagger spec requires security to be an array of objects
	Security *SecurityRequirement `json:"security,omitempty"`
}

type SecurityRequirement struct {
	Requirements    []map[string][]string
	DisableSecurity bool
}

func (s *SecurityRequirement) MarshalJSON() ([]byte, error) {
	if s.DisableSecurity {
		return []byte("[]"), nil
	}

	return json.Marshal(s.Requirements)
}

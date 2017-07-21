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
package swag_test

import (
	"testing"

	"github.com/savaki/swag"
	"github.com/savaki/swag/swagger"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	api := swag.New(
		swag.Description("blah"),
	)
	assert.Equal(t, "blah", api.Info.Description)
}

func TestVersion(t *testing.T) {
	api := swag.New(
		swag.Version("blah"),
	)
	assert.Equal(t, "blah", api.Info.Version)
}

func TestTermsOfService(t *testing.T) {
	api := swag.New(
		swag.TermsOfService("blah"),
	)
	assert.Equal(t, "blah", api.Info.TermsOfService)
}

func TestTitle(t *testing.T) {
	api := swag.New(
		swag.Title("blah"),
	)
	assert.Equal(t, "blah", api.Info.Title)
}

func TestContactEmail(t *testing.T) {
	api := swag.New(
		swag.ContactEmail("blah"),
	)
	assert.Equal(t, "blah", api.Info.Contact.Email)
}

func TestLicense(t *testing.T) {
	api := swag.New(
		swag.License("name", "url"),
	)
	assert.Equal(t, "name", api.Info.License.Name)
	assert.Equal(t, "url", api.Info.License.URL)
}

func TestBasePath(t *testing.T) {
	api := swag.New(
		swag.BasePath("/"),
	)
	assert.Equal(t, "/", api.BasePath)
}

func TestSchemes(t *testing.T) {
	api := swag.New(
		swag.Schemes("blah"),
	)
	assert.Equal(t, []string{"blah"}, api.Schemes)
}

func TestTag(t *testing.T) {
	api := swag.New(
		swag.Tag("name", "desc",
			swag.TagDescription("ext-desc"),
			swag.TagURL("ext-url"),
		),
	)

	expected := swagger.Tag{
		Name:        "name",
		Description: "desc",
		Docs: swagger.Docs{
			Description: "ext-desc",
			URL:         "ext-url",
		},
	}
	assert.Equal(t, expected, api.Tags[0])
}

func TestHost(t *testing.T) {
	api := swag.New(
		swag.Host("blah"),
	)
	assert.Equal(t, "blah", api.Host)
}

func TestSecurityScheme(t *testing.T) {
	api := swag.New(
		swag.SecurityScheme("basic", swagger.BasicSecurity()),
		swag.SecurityScheme("apikey", swagger.APIKeySecurity("Authorization", "header")),
	)
	assert.Len(t, api.SecurityDefinitions, 2)
	assert.Contains(t, api.SecurityDefinitions, "basic")
	assert.Contains(t, api.SecurityDefinitions, "apikey")
	assert.Equal(t, "header", api.SecurityDefinitions["apikey"].In)
}

func TestSecurity(t *testing.T) {
	api := swag.New(
		swag.Security("basic"),
	)
	assert.Len(t, api.Security.Requirements, 1)
	assert.Contains(t, api.Security.Requirements[0], "basic")
}

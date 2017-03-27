package swaggering_test

import (
	"testing"

	"github.com/savaki/swaggering"
	"github.com/savaki/swaggering/swagger"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	builder := swaggering.New(
		swaggering.Description("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Info.Description)
}

func TestVersion(t *testing.T) {
	builder := swaggering.New(
		swaggering.Version("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Info.Version)
}

func TestTermsOfService(t *testing.T) {
	builder := swaggering.New(
		swaggering.TermsOfService("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Info.TermsOfService)
}

func TestTitle(t *testing.T) {
	builder := swaggering.New(
		swaggering.Title("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Info.Title)
}

func TestContactEmail(t *testing.T) {
	builder := swaggering.New(
		swaggering.ContactEmail("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Info.Contact.Email)
}

func TestLicense(t *testing.T) {
	builder := swaggering.New(
		swaggering.License("name", "url"),
	)
	assert.Equal(t, "name", builder.Api.Info.License.Name)
	assert.Equal(t, "url", builder.Api.Info.License.Url)
}

func TestBasePath(t *testing.T) {
	builder := swaggering.New(
		swaggering.BasePath("/"),
	)
	assert.Equal(t, "/", builder.Api.BasePath)
}

func TestSchemes(t *testing.T) {
	builder := swaggering.New(
		swaggering.Schemes("blah"),
	)
	assert.Equal(t, []string{"blah"}, builder.Api.Schemes)
}

func TestTag(t *testing.T) {
	builder := swaggering.New(
		swaggering.Tag("name", "desc",
			swaggering.TagDescription("ext-desc"),
			swaggering.TagUrl("ext-url"),
		),
	)

	expected := swagger.Tag{
		Name:        "name",
		Description: "desc",
		Docs: swagger.Docs{
			Description: "ext-desc",
			Url:         "ext-url",
		},
	}
	assert.Equal(t, expected, builder.Api.Tags[0])
}

func TestHost(t *testing.T) {
	builder := swaggering.New(
		swaggering.Host("blah"),
	)
	assert.Equal(t, "blah", builder.Api.Host)
}

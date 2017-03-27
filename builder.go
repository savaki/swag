package swag

import "github.com/savaki/swag/swagger"

type Builder struct {
	Api *swagger.Api
}

// Option provides configuration options to the swagger api builder
type Option func(builder *Builder)

// Description sets info.description
func Description(v string) Option {
	return func(builder *Builder) {
		builder.Api.Info.Description = v
	}
}

// Version sets info.version
func Version(v string) Option {
	return func(builder *Builder) {
		builder.Api.Info.Version = v
	}
}

// TermsOfService sets info.termsOfService
func TermsOfService(v string) Option {
	return func(builder *Builder) {
		builder.Api.Info.TermsOfService = v
	}
}

// Title sets info.title
func Title(v string) Option {
	return func(builder *Builder) {
		builder.Api.Info.Title = v
	}
}

// ContactEmail sets info.contact.email
func ContactEmail(v string) Option {
	return func(builder *Builder) {
		builder.Api.Info.Contact.Email = v
	}
}

// License sets both info.license.name and info.license.url
func License(name, url string) Option {
	return func(builder *Builder) {
		builder.Api.Info.License.Name = name
		builder.Api.Info.License.Url = url
	}
}

// BasePath sets basePath
func BasePath(v string) Option {
	return func(builder *Builder) {
		builder.Api.BasePath = v
	}
}

// Schemes sets the scheme
func Schemes(v ...string) Option {
	return func(builder *Builder) {
		builder.Api.Schemes = v
	}
}

// TagOption provides additional customizations to the #Tag option
type TagOption func(tag *swagger.Tag)

// ExternalTagDescription sets externalDocs.description on the tag field
func TagDescription(v string) TagOption {
	return func(t *swagger.Tag) {
		t.Docs.Description = v
	}
}

// ExternalTagUrl sets externalDocs.url on the tag field
func TagUrl(v string) TagOption {
	return func(t *swagger.Tag) {
		t.Docs.Url = v
	}
}

// Tag adds a tag to the swagger api
func Tag(name, description string, options ...TagOption) Option {
	return func(builder *Builder) {
		if builder.Api.Tags == nil {
			builder.Api.Tags = []swagger.Tag{}
		}

		t := swagger.Tag{
			Name:        name,
			Description: description,
		}

		for _, opt := range options {
			opt(&t)
		}

		builder.Api.Tags = append(builder.Api.Tags, t)
	}
}

// Host specifies the host field
func Host(v string) Option {
	return func(builder *Builder) {
		builder.Api.Host = v
	}
}

// Endpoints allows the endpoints to be added dynamically to the Api
func Endpoints(endpoints ...*swagger.Endpoint) Option {
	return func(builder *Builder) {
		for _, e := range endpoints {
			builder.Api.AddEndpoint(e)
		}
	}
}

// New constructs a new api builder
func New(options ...Option) *Builder {
	b := &Builder{
		Api: &swagger.Api{
			BasePath: "/",
			Swagger:  "2.0",
			Schemes:  []string{"http"},
			Info: swagger.Info{
				Contact: swagger.Contact{
					Email: "your-email-address",
				},
				Description:    "Describe your API",
				Title:          "Your API Title",
				Version:        "SNAPSHOT",
				TermsOfService: "http://swagger.io/terms/",
				License: swagger.License{
					Name: "Apache 2.0",
					Url:  "http://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
		},
	}

	for _, opt := range options {
		opt(b)
	}

	return b
}

func (b *Builder) Build() *swagger.Api {
	return b.Api
}

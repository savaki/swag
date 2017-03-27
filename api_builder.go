package swaggering

type apiBuilder struct {
	Api *Api
}

// Option provides configuration options to the swagger api builder
type Option func(builder *apiBuilder)

// Description sets info.description
func Description(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.Description = v
	}
}

// Version sets info.version
func Version(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.Version = v
	}
}

// TermsOfService sets info.termsOfService
func TermsOfService(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.TermsOfService = v
	}
}

// Title sets info.title
func Title(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.Title = v
	}
}

// ContactEmail sets info.contact.email
func ContactEmail(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.Contact.Email = v
	}
}

// License sets both info.license.name and info.license.url
func LicenseNameAndUrl(name, url string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Info.License.Name = name
		builder.Api.Info.License.Url = url
	}
}

// BasePath sets basePath
func BasePath(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.BasePath = v
	}
}

// Schemes sets the scheme
func Schemes(v ...string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Schemes = v
	}
}

// TagOption provides additional customizations to the #Tag option
type TagOption func(*OldTag)

// ExternalTagDescription sets externalDocs.description on the tag field
func ExternalTagDescription(v string) TagOption {
	return func(t *OldTag) {
	}
}

// ExternalTagUrl sets externalDocs.url on the tag field
func ExternalTagUrl(v string) TagOption {
	return func(t *OldTag) {
	}
}

// Tag adds a tag to the swagger api
func Tag(name, description string, options ...TagOption) Option {
	return func(builder *apiBuilder) {
	}
}

// Host specifies the host field
func Host(v string) Option {
	return func(builder *apiBuilder) {
		builder.Api.Host = v
	}
}

// Endpoints allows the endpoints to be added dynamically to the Api
func WithEndpoints(endpoints ...OldEndpoint) Option {
	return func(builder *apiBuilder) {
		//builder.Api.Host = v
	}
}

// New constructs a new api builder
func New(options ...Option) *apiBuilder {
	b := &apiBuilder{
		Api: &Api{
			BasePath: "/",
			Swagger:  "2.0",
			Schemes:  []string{"http"},
			Info: Info{
				Contact: Contact{
					Email: "your-email-address",
				},
				Description:    "Describe your API",
				Title:          "Your API Title",
				Version:        "SNAPSHOT",
				TermsOfService: "http://swagger.io/terms/",
				License: License{
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

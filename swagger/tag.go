package swagger

// Docs represents tag docs from the swagger definition
type Docs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Tag represents a swagger tag
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Docs        Docs   `json:"externalDocs"`
}

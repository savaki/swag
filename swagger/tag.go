package swagger

type Docs struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Docs        Docs   `json:"externalDocs"`
}

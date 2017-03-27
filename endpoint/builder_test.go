package endpoint_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/savaki/swaggering/endpoint"
	"github.com/savaki/swaggering/swagger"
	"github.com/stretchr/testify/assert"
)

func Echo(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello world")
}

func TestNew(t *testing.T) {
	summary := "here's the summary"
	e := endpoint.New("get", "/", Echo,
		endpoint.Summary(summary),
	).Endpoint

	assert.Equal(t, "GET", e.Method)
	assert.Equal(t, "/", e.Path)
	assert.NotNil(t, e.Handler)
	assert.Equal(t, []string{"application/json"}, e.Consumes)
	assert.Equal(t, []string{"application/json"}, e.Produces)
	assert.Equal(t, summary, e.Summary)
	assert.Equal(t, []string{}, e.Tags)
}

func TestTags(t *testing.T) {
	e := endpoint.New("get", "/", Echo,
		endpoint.Tags("blah"),
	).Endpoint

	assert.Equal(t, []string{"blah"}, e.Tags)
}

func TestDescription(t *testing.T) {
	e := endpoint.New("get", "/", Echo,
		endpoint.Description("blah"),
	).Endpoint

	assert.Equal(t, "blah", e.Description)
}

func TestOperationId(t *testing.T) {
	e := endpoint.New("get", "/", Echo,
		endpoint.OperationId("blah"),
	).Endpoint

	assert.Equal(t, "blah", e.OperationId)
}

func TestProduces(t *testing.T) {
	expected := []string{"a", "b"}
	e := endpoint.New("get", "/", Echo,
		endpoint.Produces(expected...),
	).Endpoint

	assert.Equal(t, expected, e.Produces)
}

func TestConsumes(t *testing.T) {
	expected := []string{"a", "b"}
	e := endpoint.New("get", "/", Echo,
		endpoint.Consumes(expected...),
	).Endpoint

	assert.Equal(t, expected, e.Consumes)
}

func TestPath(t *testing.T) {
	expected := swagger.Parameter{
		In:          "path",
		Name:        "id",
		Description: "the description",
		Required:    true,
		Type:        "string",
	}

	e := endpoint.New("get", "/", Echo,
		endpoint.Path(expected.Name, expected.Type, expected.Description, expected.Required),
	).Endpoint

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestQuery(t *testing.T) {
	expected := swagger.Parameter{
		In:          "query",
		Name:        "id",
		Description: "the description",
		Required:    true,
		Type:        "string",
	}

	e := endpoint.New("get", "/", Echo,
		endpoint.Query(expected.Name, expected.Type, expected.Description, expected.Required),
	).Endpoint

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

type Model struct {
	String string `json:"s"`
}

func TestBody(t *testing.T) {
	expected := swagger.Parameter{
		In:          "body",
		Description: "the description",
		Required:    true,
		Schema: &swagger.Schema{
			Type:      "object",
			Ref:       "#/definitions/endpoint_testModel",
			Prototype: Model{},
		},
	}

	e := endpoint.New("get", "/", Echo,
		endpoint.Body(Model{}, expected.Description, expected.Required),
	).Endpoint

	assert.Equal(t, 1, len(e.Parameters))
	assert.Equal(t, expected, e.Parameters[0])
}

func TestResponse(t *testing.T) {
	expected := swagger.Response{
		Description: "successful",
		Schema: &swagger.Schema{
			Type:      "object",
			Ref:       "#/definitions/endpoint_testModel",
			Prototype: Model{},
		},
	}

	e := endpoint.New("get", "/", Echo,
		endpoint.Response(http.StatusOK, Model{}, "successful"),
	).Endpoint

	assert.Equal(t, 1, len(e.Responses))
	assert.Equal(t, expected, e.Responses["200"])
}

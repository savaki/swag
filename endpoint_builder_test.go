package swaggering_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/savaki/swaggering"
	"github.com/stretchr/testify/assert"
)

type Login struct {
	Username string
	Password string
}

type Session struct {
	UserID int
}

func TestBuilder(t *testing.T) {
	b := swaggering.NewEndpoint("get", "/", "the summary", nil,
		swaggering.EndpointDescription("the description"),
	).
		Summary("the summary").
		Description("the description").
		Tags("tag1", "tag2").
		Query("q", "string", "q string", true).
		Path("p", "int", "p string", true).
		Body(Login{}, "login object", true).
		Response(http.StatusOK, Session{}, "successful login")

	data, err := json.Marshal(b.Endpoint)
	assert.Nil(t, err)
	endpoint := swaggering.OldEndpoint{}
	err = json.Unmarshal(data, &endpoint)
	assert.Nil(t, err)

	data, err = ioutil.ReadFile("testdata/builder.json")
	assert.Nil(t, err)
	expected := swaggering.OldEndpoint{}
	err = json.Unmarshal(data, &expected)
	assert.Nil(t, err)

	assert.Equal(t, expected, endpoint)
}

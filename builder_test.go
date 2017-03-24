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
	b := swaggering.New("get", "/", nil)
	b.Summary("the summary")
	b.Description("the description")
	b.Tags("tag1", "tag2")
	b.Query("q", "string", "q string", true)
	b.Path("p", "int", "p string", true)
	b.Schema(Login{}, "login object", true)
	b.Response(http.StatusOK, Session{}, "successful login")

	data, err := json.Marshal(b.Endpoint)
	assert.Nil(t, err)
	endpoint := swaggering.Endpoint{}
	err = json.Unmarshal(data, &endpoint)
	assert.Nil(t, err)

	data, err = ioutil.ReadFile("testdata/builder.json")
	assert.Nil(t, err)
	expected := swaggering.Endpoint{}
	err = json.Unmarshal(data, &expected)
	assert.Nil(t, err)

	assert.Equal(t, expected, endpoint)
}

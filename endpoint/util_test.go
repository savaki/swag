package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamel(t *testing.T) {
	assert.Equal(t, "HelloWorld", camel("hello/world"))
	assert.Equal(t, "UsersUser", camel("/users/{user}"))
}

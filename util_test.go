package swag_test

import (
	"testing"

	"github.com/savaki/swag"
	"github.com/stretchr/testify/assert"
)

func TestColonPath(t *testing.T) {
	assert.Equal(t, "/api/:id", swag.ColonPath("/api/{id}"))
	assert.Equal(t, "/api/:a/:b/:c", swag.ColonPath("/api/{a}/{b}/{c}"))
}

package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Mock struct {
	name string
	pkg  string
}

func (m Mock) PkgPath() string {
	return m.pkg
}

func (m Mock) Name() string {
	return m.name
}

func TestMakeSchema(t *testing.T) {
	name := makeName(Mock{
		name: "Name",
		pkg:  "with-some-dashes",
	})
	assert.Equal(t, "with_some_dashesName", name)
}

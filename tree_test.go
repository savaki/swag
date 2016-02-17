package swaggering

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTree(t *testing.T) {
	Convey("Given an Endpoint", t, func() {
		endpoint := Endpoint{
			Method: "GET",
			Path:   "/hello/a/b/c",
		}

		tree := &Tree{}
		tree.register(endpoint)

		So(tree.child("hello").child("a").child("b").child("c").Path(), ShouldEqual, endpoint.Path)
	})
}

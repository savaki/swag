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

		So(tree, ShouldResemble, &Tree{
			Children: map[string]*Tree{
				"hello": &Tree{
					Children: map[string]*Tree{
						"a": &Tree{
							Children: map[string]*Tree{
								"b": &Tree{
									Children: map[string]*Tree{
										"c": &Tree{
											Endpoints: map[string]Endpoint{
												endpoint.Method: endpoint,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})
	})
}

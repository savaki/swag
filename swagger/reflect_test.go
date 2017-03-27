package swagger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	First string
}

type Pet struct {
	Friend      Person    `json:"friend"`
	Friends     []Person  `json:"friends"`
	Pointer     *Person   `json:"pointer" required:"true"`
	Pointers    []*Person `json:"pointers"`
	Int         int
	IntArray    []int
	String      string
	StringArray []string
}

type Empty struct {
	Nope int `json:"-"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func TestDefine(t *testing.T) {
	t.Run("Given a thing", func(t *testing.T) {
		v := define(Pet{})
		obj, ok := v["swaggerPet"]
		assert.True(t, ok)
		assert.False(t, obj.IsArray)
		assert.Equal(t, 8, len(obj.Properties))

		content := map[string]Object{}
		data, err := ioutil.ReadFile("testdata/pet.json")
		assert.Nil(t, err)
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&content)
		assert.Nil(t, err)
		expected := content["swaggerPet"]

		assert.Equal(t, expected.IsArray, obj.IsArray, "expected IsArray to match")
		assert.Equal(t, expected.Type, obj.Type, "expected Type to match")
		assert.Equal(t, expected.Required, obj.Required, "expected Required to match")
		assert.Equal(t, len(expected.Properties), len(obj.Properties), "expected same number of properties")

		for k, v := range obj.Properties {
			e := expected.Properties[k]
			assert.Equal(t, e.Type, v.Type, "expected %v.Type to match", k)
			assert.Equal(t, e.Description, v.Description, "expected %v.Required to match", k)
			assert.Equal(t, e.Enum, v.Enum, "expected %v.Required to match", k)
			assert.Equal(t, e.Format, v.Format, "expected %v.Required to match", k)
			assert.Equal(t, e.Ref, v.Ref, "expected %v.Required to match", k)
			assert.Equal(t, e.Example, v.Example, "expected %v.Required to match", k)
			assert.Equal(t, e.Items, v.Items, "expected %v.Required to match", k)
		}
	})
}

func TestHonorJsonIgnore(t *testing.T) {
	v := define(Empty{})
	obj, ok := v["swaggerEmpty"]
	assert.True(t, ok)
	assert.False(t, obj.IsArray)
	assert.Equal(t, 0, len(obj.Properties), "expected zero exposed properties")
}

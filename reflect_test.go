package swaggering

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Person struct {
	First string
}

type Pet struct {
	Friend   Person    `json:"friend"`
	Friends  []Person  `json:"friends"`
	Pointer  *Person   `json:"pointer" required:"true"`
	Pointers []*Person `json:"pointers"`
	Ints     []int
	Strings  []string
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func TestDefine(t *testing.T) {
	Convey("Given a thing", t, func() {
		v := define(Pet{})
		data, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(data))
	})
}

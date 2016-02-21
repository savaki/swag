package swaggering

import (
	"fmt"
	"path/filepath"
	"reflect"
)

func makeRef(name string) string {
	return fmt.Sprintf("#/definitions/%v", name)
}

func makeName(t reflect.Type) string {
	return filepath.Base(t.PkgPath()) + t.Name()
}

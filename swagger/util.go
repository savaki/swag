package swagger

import (
	"fmt"
	"path/filepath"
	"strings"
)

func makeRef(name string) string {
	return fmt.Sprintf("#/definitions/%v", name)
}

type reflectType interface {
	PkgPath() string
	Name() string
}

func makeName(t reflectType) string {
	name := filepath.Base(t.PkgPath()) + t.Name()
	return strings.Replace(name, "-", "_", -1)
}

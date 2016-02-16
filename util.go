package swaggering

import "fmt"

func makeRef(name string) string {
	return fmt.Sprintf("#/definitions/%v", name)
}

package endpoint

import (
	"regexp"
	"strings"
)

var (
	reAlphaNumeric = regexp.MustCompile(`[^0-9a-zA-Z]`)
)

func camel(v string) string {
	segments := strings.Split(v, "/")
	results := make([]string, 0, len(segments))

	for _, segment := range segments {
		v := reAlphaNumeric.ReplaceAllString(segment, "")
		if v == "" {
			continue
		}

		results = append(results, strings.ToUpper(v[0:1])+v[1:])
	}

	return strings.Join(results, "")
}

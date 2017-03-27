package swag

import (
	"regexp"
	"strings"
)

var (
	rePath = regexp.MustCompile(`\{([^}]+)}`)
)

// FixPath accepts a swagger path e.g. /api/orgs/{org} and returns a colon identified path e.g. /api/org/:org
func ColonPath(path string) string {
	matches := rePath.FindAllStringSubmatch(path, -1)
	if matches != nil {
		for _, match := range matches {
			path = strings.Replace(path, match[0], ":"+match[1], -1)
		}
	}
	return path
}


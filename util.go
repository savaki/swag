// Copyright 2017 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package swag

import (
	"regexp"
	"strings"
)

var (
	rePath = regexp.MustCompile(`\{([^}]+)}`)
)

// ColonPath accepts a swagger path e.g. /api/orgs/{org} and returns a colon identified path e.g. /api/org/:org
func ColonPath(path string) string {
	matches := rePath.FindAllStringSubmatch(path, -1)
	if matches != nil {
		for _, match := range matches {
			path = strings.Replace(path, match[0], ":"+match[1], -1)
		}
	}
	return path
}

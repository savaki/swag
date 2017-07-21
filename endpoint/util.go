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

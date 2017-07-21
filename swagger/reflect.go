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
package swagger

import (
	"reflect"
	"strings"
)

func inspect(t reflect.Type, jsonTag string) Property {
	p := Property{
		GoType: t,
	}

	if strings.Contains(jsonTag, ",string") {
		p.Type = "string"
		return p
	}

	switch p.GoType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		p.Type = "integer"
		p.Format = "int32"

	case reflect.Int64, reflect.Uint64:
		p.Type = "integer"
		p.Format = "int64"

	case reflect.Float64:
		p.Type = "number"
		p.Format = "double"

	case reflect.Float32:
		p.Type = "number"
		p.Format = "float"

	case reflect.Bool:
		p.Type = "boolean"

	case reflect.String:
		p.Type = "string"

	case reflect.Struct:
		name := makeName(p.GoType)
		p.Ref = makeRef(name)

	case reflect.Ptr:
		p.GoType = t.Elem()
		name := makeName(p.GoType)
		p.Ref = makeRef(name)

	case reflect.Slice:
		p.Type = "array"
		p.Items = &Items{}

		p.GoType = t.Elem() // dereference the slice
		switch p.GoType.Kind() {
		case reflect.Ptr:
			p.GoType = p.GoType.Elem()
			name := makeName(p.GoType)
			p.Items.Ref = makeRef(name)

		case reflect.Struct:
			name := makeName(p.GoType)
			p.Items.Ref = makeRef(name)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			p.Items.Type = "integer"
			p.Items.Format = "int32"

		case reflect.Int64, reflect.Uint64:
			p.Items.Type = "integer"
			p.Items.Format = "int64"

		case reflect.Float64:
			p.Items.Type = "number"
			p.Items.Format = "double"

		case reflect.Float32:
			p.Items.Type = "number"
			p.Items.Format = "float"

		case reflect.String:
			p.Items.Type = "string"
		}
	}

	return p
}

func defineObject(v interface{}) Object {
	var required []string

	var t reflect.Type
	switch value := v.(type) {
	case reflect.Type:
		t = value
	default:
		t = reflect.TypeOf(v)
	}

	properties := map[string]Property{}
	isArray := t.Kind() == reflect.Slice

	if isArray {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		p := inspect(t, "")
		return Object{
			IsArray:  isArray,
			GoType:   t,
			Type:     p.Type,
			Format:   p.Format,
			Name:     t.Kind().String(),
			Required: required,
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// skip unexported fields
		if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
			continue
		}

		// determine the json name of the field
		name := strings.TrimSpace(field.Tag.Get("json"))
		if name == "" || strings.HasPrefix(name, ",") {
			name = field.Name

		} else {
			// strip out things like , omitempty
			parts := strings.Split(name, ",")
			name = parts[0]
		}

		parts := strings.Split(name, ",") // foo,omitempty => foo
		name = parts[0]
		if name == "-" {
			// honor json ignore tag
			continue
		}

		// determine if this field is required or not
		if v := field.Tag.Get("required"); v == "true" {
			if required == nil {
				required = []string{}
			}
			required = append(required, name)
		}

		p := inspect(field.Type, field.Tag.Get("json"))
		properties[name] = p
	}

	return Object{
		IsArray:    isArray,
		GoType:     t,
		Type:       "object",
		Name:       makeName(t),
		Required:   required,
		Properties: properties,
	}
}

func define(v interface{}) map[string]Object {
	objMap := map[string]Object{}

	obj := defineObject(v)
	objMap[obj.Name] = obj

	dirty := true

	for dirty {
		dirty = false
		for _, d := range objMap {
			for _, p := range d.Properties {
				if p.GoType.Kind() == reflect.Struct {
					name := makeName(p.GoType)
					if _, exists := objMap[name]; !exists {
						child := defineObject(p.GoType)
						objMap[child.Name] = child
						dirty = true
					}
				}
			}
		}
	}

	return objMap
}

// MakeSchema takes struct or pointer to a struct and returns a Schema instance suitable for use by the swagger doc
func MakeSchema(prototype interface{}) *Schema {
	schema := &Schema{
		Prototype: prototype,
	}

	obj := defineObject(prototype)
	if obj.IsArray {
		schema.Type = "array"
		schema.Items = &Items{
			Ref: makeRef(obj.Name),
		}

	} else {
		schema.Ref = makeRef(obj.Name)
	}

	return schema
}

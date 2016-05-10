package swaggering

import (
	"reflect"
	"strings"
)

type Object struct {
	IsArray    bool                `json:"-"`
	GoType     reflect.Type        `json:"-"`
	Name       string              `json:"-"`
	Type       string              `json:"type"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Items struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	Ref    string `json:"$ref,omitempty"`
}

type Property struct {
	GoType      reflect.Type `json:"-"`
	Type        string       `json:"type,omitempty"`
	Description string       `json:"description,omitempty"`
	Enum        []string     `json:"enum,omitempty"`
	Format      string       `json:"format,omitempty"`
	Ref         string       `json:"$ref,omitempty"`
	Example     string       `json:"example,omitempty"`
	Items       *Items       `json:"items,omitempty"`
}

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
	isArray := (t.Kind() == reflect.Slice)

	if isArray {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

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

		// determine if this field is required or not
		if v := field.Tag.Get("required"); v == "true" {
			if required == nil {
				required = []string{}
			}
			required = append(required, name)
		}

		properties[name] = inspect(field.Type, field.Tag.Get("json"))
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
	objs := map[string]Object{}

	obj := defineObject(v)
	objs[obj.Name] = obj

	dirty := true

	for dirty {
		dirty = false
		for _, d := range objs {
			for _, p := range d.Properties {
				if p.GoType.Kind() == reflect.Struct {
					name := makeName(p.GoType)
					if _, exists := objs[name]; !exists {
						child := defineObject(p.GoType)
						objs[child.Name] = child
						dirty = true
					}
				}
			}
		}
	}

	return objs
}

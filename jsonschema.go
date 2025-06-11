package gopenrouter

import (
	"fmt"
	"reflect"
	"strings"
)

// GenerateSchema creates a JSON schema from a Go struct.
// It uses reflection to generate a JSON schema from the struct's fields and tags.
func GenerateSchema(v any) (map[string]any, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, but got %s", t.Kind())
	}

	schema := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"required":             []string{},
		"additionalProperties": false,
	}
	properties := schema["properties"].(map[string]any)
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonTagParts := strings.Split(jsonTag, ",")
		fieldName := jsonTagParts[0]

		prop := map[string]any{}
		switch field.Type.Kind() {
		case reflect.String:
			prop["type"] = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			prop["type"] = "integer"
		case reflect.Float32, reflect.Float64:
			prop["type"] = "number"
		case reflect.Bool:
			prop["type"] = "boolean"
		default:
			// For simplicity, this example does not handle nested structs or slices.
			continue
		}

		if desc := field.Tag.Get("jsonschema"); desc != "" {
			prop["description"] = desc
		}

		properties[fieldName] = prop

		// Check for 'required' tag or simply make all fields required by default
		required = append(required, fieldName)
	}
	schema["required"] = required
	return schema, nil
}

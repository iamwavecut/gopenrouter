package gopenrouter

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestGenerateSchema(t *testing.T) {
	type testStruct struct {
		Name         string  `json:"name" jsonschema:"The name of the user"`
		Age          int     `json:"age"`
		Verified     bool    `json:"verified"`
		Score        float64 `json:"score"`
		Ignored      string  `json:"-"`
		NoJSONTag    string
		EmptyJSONTag string `json:""`
	}

	expectedSchema := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"properties": map[string]any{
			"name":     map[string]any{"type": "string", "description": "The name of the user"},
			"age":      map[string]any{"type": "integer"},
			"verified": map[string]any{"type": "boolean"},
			"score":    map[string]any{"type": "number"},
		},
		"required": []string{"age", "name", "score", "verified"},
	}

	// Helper to compare maps by marshalling them to JSON strings
	areEqual := func(a, b map[string]any) bool {
		b1, err1 := json.Marshal(a)
		b2, err2 := json.Marshal(b)
		if err1 != nil || err2 != nil {
			return false
		}
		return string(b1) == string(b2)
	}

	// Test with a struct value
	schema, err := GenerateSchema(testStruct{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Sort required fields for consistent comparison
	sort.Strings(schema["required"].([]string))
	sort.Strings(expectedSchema["required"].([]string))

	if !areEqual(schema, expectedSchema) {
		t.Errorf("schema mismatch:\n- got: %+v\n- want: %+v", schema, expectedSchema)
	}

	// Test with a pointer
	schema, err = GenerateSchema(&testStruct{})
	if err != nil {
		t.Fatalf("unexpected error with pointer: %v", err)
	}

	sort.Strings(schema["required"].([]string))

	if !areEqual(schema, expectedSchema) {
		t.Errorf("schema mismatch with pointer:\n- got: %+v\n- want: %+v", schema, expectedSchema)
	}
}

func TestGenerateSchema_Error(t *testing.T) {
	_, err := GenerateSchema(42)
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedError := "expected a struct, but got int"
	if err.Error() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, err.Error())
	}
}

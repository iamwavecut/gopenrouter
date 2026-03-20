package jsonx

import (
	"encoding/json"
	"maps"
)

func MarshalMap(v any) (map[string]any, error) {
	if v == nil {
		return map[string]any{}, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func MergeMaps(dst map[string]any, src map[string]any) {
	maps.Copy(dst, src)
}

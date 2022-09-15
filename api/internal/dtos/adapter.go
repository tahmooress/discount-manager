package dtos

import (
	"encoding/json"
	"fmt"
)

type stringAdapter string

func (as *stringAdapter) UnmarshalJSON(b []byte) error {
	var temp interface{}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return fmt.Errorf("handlers >> adapterString >> unmarshaler >> %w", err)
	}

	switch v := temp.(type) {
	case string:
		*as = stringAdapter(v)
	default:
		*as = stringAdapter(string(b))
	}

	return nil
}

func (as *stringAdapter) String() string {
	if *as == "null" {
		return ""
	}

	return string(*as)
}

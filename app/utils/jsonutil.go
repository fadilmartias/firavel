package utils

import (
	"github.com/bytedance/sonic"
)

func JSONParse(input any) map[string]any {
	var result map[string]any
	var raw []byte

	switch v := input.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return map[string]any{}
	}

	err := sonic.Unmarshal(raw, &result)
	if err != nil {
		return map[string]any{}
	}
	return result
}

func JSONStringify(v any) string {
	str, err := sonic.MarshalString(v)
	if err != nil {
		return ""
	}
	return str
}

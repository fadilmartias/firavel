package utils

import (
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
)

func StringPtr(v any) *string {
	if v == nil {
		return nil
	}

	// Coba konversi langsung ke string
	if str, ok := v.(string); ok {
		if str == "" || str == "<nil>" {
			return nil
		}
		return &str
	}

	// Kalau bukan string tapi bisa diformat jadi string
	str := fmt.Sprintf("%v", v)
	if str == "" || str == "<nil>" || str == "map[]" {
		return nil
	}

	return &str
}

func IntPtr(v any) *int {
	if v == nil {
		return nil
	}

	// Coba konversi langsung ke int
	if num, ok := v.(int); ok {
		return &num
	}

	// Kalau bukan int tapi bisa diformat jadi int
	str := fmt.Sprintf("%v", v)
	if str == "" || str == "<nil>" || str == "map[]" {
		return nil
	}

	// Coba konversi string ke int
	if num, err := strconv.Atoi(str); err == nil {
		return &num
	}

	return nil
}

func StructToMap(data any) map[string]any {
	var result map[string]any
	jsonBytes, _ := sonic.Marshal(data)
	sonic.Unmarshal(jsonBytes, &result)
	return result
}

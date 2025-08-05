package utils

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
)

// DD seperti dd() di Laravel: dump dan die
func DD(data any, die ...bool) {
	shouldDie := true // default

	if len(die) > 0 {
		shouldDie = die[0]
	}

	bytes, err := sonic.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error while marshaling:", err)
	} else {
		fmt.Println(string(bytes))
	}

	if shouldDie {
		os.Exit(1)
	}
}

func Dump(data any) string {
	bytes, err := sonic.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error while marshaling: %v", err)
	}
	return string(bytes)
}

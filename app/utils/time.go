package utils

import "time"

func Timestamp() string {
	return time.Now().Format("20060102150405")
}

func isValidTimeFormat(s string) bool {
	_, err := time.Parse("15:04", s)
	return err == nil
}

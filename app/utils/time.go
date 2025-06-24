package utils

import "time"

func Timestamp() string {
	return time.Now().Format("20060102150405")
}

package utils

func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetMap(data interface{}) (map[string]interface{}, bool) {
	m, ok := data.(map[string]interface{})
	return m, ok
}

func GetSlice(data interface{}) ([]interface{}, bool) {
	s, ok := data.([]interface{})
	return s, ok
}

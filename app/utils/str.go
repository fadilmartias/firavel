package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SplitWords(s string) []string {
	// Ganti semua pemisah umum dengan spasi
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.TrimSpace(s)

	// Pisahkan camelCase & PascalCase (usVo, UserVoucher)
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}

	return strings.FieldsFunc(string(result), func(r rune) bool {
		return r == ' '
	})
}

func SnakeCase(s string) string {
	words := SplitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

func KebabCase(s string) string {
	words := SplitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "-")
}

func StudlyCase(s string) string {
	words := SplitWords(s)
	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, "")
}

func CamelCase(s string) string {
	studly := StudlyCase(s)
	if len(studly) == 0 {
		return ""
	}
	return strings.ToLower(studly[:1]) + studly[1:]
}

func SlugCase(s string) string {
	slug := strings.ToLower(s)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, "_", " ")
	slug = strings.TrimSpace(slug)
	slug = regexp.MustCompile(`[\s\-_]+`).ReplaceAllString(slug, "-")
	return slug
}

func TitleCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, " ")
}

func LimitWords(s string, limit int) string {
	words := strings.Fields(s)
	if len(words) <= limit {
		return s
	}
	return strings.Join(words[:limit], " ") + "..."
}

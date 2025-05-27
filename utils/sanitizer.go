package utils

import (
	"regexp"
	"strings"
)

// Sanitize removes leading/trailing spaces and non-alphanumeric characters from the input string.
func Sanitize(input string) string {
	input = strings.TrimSpace(input)
	// Allow spaces within the string, remove other special characters
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.ReplaceAllString(input, "")
}

func SanitizeArrayStrings(input string) []string {
	sanitized := make([]string, 0)
	for _, str := range strings.Split(input, ",") {
		sanitized = append(sanitized, Sanitize(str))
	}
	return sanitized
}

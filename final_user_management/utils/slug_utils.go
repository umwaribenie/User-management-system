package utils

import (
	"strings"
	"unicode"
)

// GenerateSlug converts a string (e.g., a name) into a URL-friendly slug.
func GenerateSlug(text string) string {
	var slug strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			slug.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			slug.WriteRune('-')
		}
	}

	// Remove any trailing hyphens
	result := strings.Trim(slug.String(), "-")

	// If the result is empty, return a default slug
	if result == "" {
		return "user"
	}

	return result
}

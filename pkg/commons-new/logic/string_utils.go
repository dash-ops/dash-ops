package commons

import (
	"regexp"
	"strings"
)

// StringProcessor provides string processing utilities
type StringProcessor struct{}

// NewStringProcessor creates a new string processor
func NewStringProcessor() *StringProcessor {
	return &StringProcessor{}
}

// ToUnderscore converts a string to underscore format
// Example: "Hello World" -> "hello_world"
func (s *StringProcessor) ToUnderscore(str string) string {
	if str == "" {
		return ""
	}

	// Convert to lower case
	result := strings.ToLower(str)

	// Replace spaces and tabs with underscores
	spaceRegex := regexp.MustCompile(`[[:space:][:blank:]]+`)
	result = spaceRegex.ReplaceAllString(result, "_")

	// Replace special characters with underscores
	specialCharRegex := regexp.MustCompile(`[^a-z0-9_]`)
	result = specialCharRegex.ReplaceAllString(result, "_")

	// Remove multiple consecutive underscores
	multiUnderscoreRegex := regexp.MustCompile(`_+`)
	result = multiUnderscoreRegex.ReplaceAllString(result, "_")

	// Trim underscores from beginning and end
	result = strings.Trim(result, "_")

	return result
}

// ToKebabCase converts a string to kebab-case format
// Example: "Hello World" -> "hello-world"
func (s *StringProcessor) ToKebabCase(str string) string {
	if str == "" {
		return ""
	}

	// Convert to lower case
	result := strings.ToLower(str)

	// Replace spaces and tabs with hyphens
	spaceRegex := regexp.MustCompile(`[[:space:][:blank:]]+`)
	result = spaceRegex.ReplaceAllString(result, "-")

	// Replace special characters with hyphens
	specialCharRegex := regexp.MustCompile(`[^a-z0-9-]`)
	result = specialCharRegex.ReplaceAllString(result, "-")

	// Remove multiple consecutive hyphens
	multiHyphenRegex := regexp.MustCompile(`-+`)
	result = multiHyphenRegex.ReplaceAllString(result, "-")

	// Trim hyphens from beginning and end
	result = strings.Trim(result, "-")

	return result
}

// Sanitize removes or replaces dangerous characters from a string
func (s *StringProcessor) Sanitize(str string) string {
	if str == "" {
		return ""
	}

	// Remove control characters
	controlCharRegex := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	result := controlCharRegex.ReplaceAllString(str, "")

	// Trim whitespace
	result = strings.TrimSpace(result)

	return result
}

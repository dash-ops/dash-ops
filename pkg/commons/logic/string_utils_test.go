package commons

import (
	"testing"
)

func TestStringProcessor_ToUnderscore(t *testing.T) {
	processor := NewStringProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic conversion",
			input:    "Hello World",
			expected: "hello_world",
		},
		{
			name:     "special characters",
			input:    "Hello-World!@#",
			expected: "hello_world",
		},
		{
			name:     "multiple spaces",
			input:    "Hello    World",
			expected: "hello_world",
		},
		{
			name:     "tabs and spaces",
			input:    "Hello\t\nWorld",
			expected: "hello_world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already underscore",
			input:    "hello_world",
			expected: "hello_world",
		},
		{
			name:     "leading/trailing underscores",
			input:    "_Hello World_",
			expected: "hello_world",
		},
		{
			name:     "numbers and letters",
			input:    "Test123 Service",
			expected: "test123_service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ToUnderscore(tt.input)
			if result != tt.expected {
				t.Errorf("ToUnderscore(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringProcessor_ToKebabCase(t *testing.T) {
	processor := NewStringProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic conversion",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "special characters",
			input:    "Hello_World!@#",
			expected: "hello-world",
		},
		{
			name:     "multiple spaces",
			input:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "already kebab-case",
			input:    "hello-world",
			expected: "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.ToKebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToKebabCase(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringProcessor_Sanitize(t *testing.T) {
	processor := NewStringProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "with control characters",
			input:    "Hello\x00\x1f World",
			expected: "Hello World",
		},
		{
			name:     "with whitespace",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only control characters",
			input:    "\x00\x1f\x7f",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

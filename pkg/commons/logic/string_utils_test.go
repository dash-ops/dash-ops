package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringProcessor_ToUnderscore_WithBasicConversion_ReturnsUnderscoreCase(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello World"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithSpecialCharacters_ReturnsCleanedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello-World!@#"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithMultipleSpaces_ReturnsNormalizedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello    World"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithTabsAndSpaces_ReturnsNormalizedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello\t\nWorld"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := ""

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "", result)
}

func TestStringProcessor_ToUnderscore_WithAlreadyUnderscore_ReturnsSameString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "hello_world"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithLeadingTrailingUnderscores_ReturnsCleanedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "_Hello World_"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "hello_world", result)
}

func TestStringProcessor_ToUnderscore_WithNumbersAndLetters_ReturnsUnderscoreCase(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Test123 Service"

	// Act
	result := processor.ToUnderscore(input)

	// Assert
	assert.Equal(t, "test123_service", result)
}

func TestStringProcessor_ToKebabCase_WithBasicConversion_ReturnsKebabCase(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello World"

	// Act
	result := processor.ToKebabCase(input)

	// Assert
	assert.Equal(t, "hello-world", result)
}

func TestStringProcessor_ToKebabCase_WithSpecialCharacters_ReturnsCleanedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello_World!@#"

	// Act
	result := processor.ToKebabCase(input)

	// Assert
	assert.Equal(t, "hello-world", result)
}

func TestStringProcessor_ToKebabCase_WithMultipleSpaces_ReturnsNormalizedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello    World"

	// Act
	result := processor.ToKebabCase(input)

	// Assert
	assert.Equal(t, "hello-world", result)
}

func TestStringProcessor_ToKebabCase_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := ""

	// Act
	result := processor.ToKebabCase(input)

	// Assert
	assert.Equal(t, "", result)
}

func TestStringProcessor_ToKebabCase_WithAlreadyKebabCase_ReturnsSameString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "hello-world"

	// Act
	result := processor.ToKebabCase(input)

	// Assert
	assert.Equal(t, "hello-world", result)
}

func TestStringProcessor_Sanitize_WithBasicString_ReturnsSameString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello World"

	// Act
	result := processor.Sanitize(input)

	// Assert
	assert.Equal(t, "Hello World", result)
}

func TestStringProcessor_Sanitize_WithControlCharacters_ReturnsCleanedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "Hello\x00\x1f World"

	// Act
	result := processor.Sanitize(input)

	// Assert
	assert.Equal(t, "Hello World", result)
}

func TestStringProcessor_Sanitize_WithWhitespace_ReturnsTrimmedString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "  Hello World  "

	// Act
	result := processor.Sanitize(input)

	// Assert
	assert.Equal(t, "Hello World", result)
}

func TestStringProcessor_Sanitize_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := ""

	// Act
	result := processor.Sanitize(input)

	// Assert
	assert.Equal(t, "", result)
}

func TestStringProcessor_Sanitize_WithOnlyControlCharacters_ReturnsEmptyString(t *testing.T) {
	// Arrange
	processor := NewStringProcessor()
	input := "\x00\x1f\x7f"

	// Act
	result := processor.Sanitize(input)

	// Assert
	assert.Equal(t, "", result)
}

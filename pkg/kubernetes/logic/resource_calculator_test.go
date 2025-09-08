package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceCalculator_ParseCPU(t *testing.T) {
	calculator := NewResourceCalculator()

	tests := []struct {
		name        string
		cpu         string
		expected    int64
		expectError bool
	}{
		{
			name:     "empty string",
			cpu:      "",
			expected: 0,
		},
		{
			name:     "millicores",
			cpu:      "100m",
			expected: 100,
		},
		{
			name:     "decimal cores",
			cpu:      "0.5",
			expected: 500,
		},
		{
			name:     "integer cores",
			cpu:      "2",
			expected: 2000,
		},
		{
			name:     "with whitespace",
			cpu:      " 100m ",
			expected: 100,
		},
		{
			name:        "invalid format",
			cpu:         "100cores",
			expectError: true,
		},
		{
			name:        "invalid number",
			cpu:         "abc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.ParseCPU(tt.cpu)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestResourceCalculator_ParseMemory(t *testing.T) {
	calculator := NewResourceCalculator()

	tests := []struct {
		name        string
		memory      string
		expected    int64
		expectError bool
	}{
		{
			name:     "empty string",
			memory:   "",
			expected: 0,
		},
		{
			name:     "bytes",
			memory:   "1024",
			expected: 1024,
		},
		{
			name:     "kibibytes",
			memory:   "1Ki",
			expected: 1024,
		},
		{
			name:     "mebibytes",
			memory:   "1Mi",
			expected: 1024 * 1024,
		},
		{
			name:     "gibibytes",
			memory:   "1Gi",
			expected: 1024 * 1024 * 1024,
		},
		{
			name:     "kilobytes",
			memory:   "1K",
			expected: 1000,
		},
		{
			name:     "megabytes",
			memory:   "1M",
			expected: 1000 * 1000,
		},
		{
			name:     "gigabytes",
			memory:   "1G",
			expected: 1000 * 1000 * 1000,
		},
		{
			name:     "with whitespace",
			memory:   " 512Mi ",
			expected: 512 * 1024 * 1024,
		},
		{
			name:        "invalid format",
			memory:      "1MB",
			expectError: true,
		},
		{
			name:        "invalid number",
			memory:      "abcMi",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.ParseMemory(tt.memory)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestResourceCalculator_FormatCPU(t *testing.T) {
	calculator := NewResourceCalculator()

	tests := []struct {
		name       string
		millicores int64
		expected   string
	}{
		{
			name:       "zero",
			millicores: 0,
			expected:   "0",
		},
		{
			name:       "millicores",
			millicores: 100,
			expected:   "100m",
		},
		{
			name:       "exact cores",
			millicores: 2000,
			expected:   "2",
		},
		{
			name:       "decimal cores",
			millicores: 1500,
			expected:   "1.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.FormatCPU(tt.millicores)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceCalculator_FormatMemory(t *testing.T) {
	calculator := NewResourceCalculator()

	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "zero",
			bytes:    0,
			expected: "0",
		},
		{
			name:     "bytes",
			bytes:    512,
			expected: "512",
		},
		{
			name:     "kibibytes",
			bytes:    2048,
			expected: "2Ki",
		},
		{
			name:     "mebibytes",
			bytes:    512 * 1024 * 1024,
			expected: "512Mi",
		},
		{
			name:     "gibibytes",
			bytes:    2 * 1024 * 1024 * 1024,
			expected: "2.0Gi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.FormatMemory(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceCalculator_CalculateResourceUtilization(t *testing.T) {
	calculator := NewResourceCalculator()

	tests := []struct {
		name     string
		used     int64
		total    int64
		expected float64
	}{
		{
			name:     "50% utilization",
			used:     50,
			total:    100,
			expected: 50.0,
		},
		{
			name:     "100% utilization",
			used:     100,
			total:    100,
			expected: 100.0,
		},
		{
			name:     "0% utilization",
			used:     0,
			total:    100,
			expected: 0.0,
		},
		{
			name:     "zero total",
			used:     50,
			total:    0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculator.CalculateResourceUtilization(tt.used, tt.total)
			assert.Equal(t, tt.expected, result)
		})
	}
}

package jsonstr

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		compact     bool
		expected    string
		expectError bool
	}{
		{
			name:        "Simple JSON object",
			input:       `{"name":"John","age":30}`,
			compact:     false,
			expected:    `{\"name\":\"John\",\"age\":30}`,
			expectError: false,
		},
		{
			name:        "Pretty JSON object",
			input:       "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			compact:     false,
			expected:    `{\n  \"name\": \"John\",\n  \"age\": 30\n}`,
			expectError: false,
		},
		{
			name:        "Pretty JSON object with compact option",
			input:       "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			compact:     true,
			expected:    `{\"age\":30,\"name\":\"John\"}`,
			expectError: false,
		},
		{
			name:        "JSON with special characters",
			input:       `{"message":"Hello \"world\"","path":"C:\\path\\to\\file"}`,
			compact:     false,
			expected:    `{\"message\":\"Hello \\\"world\\\"\",\"path\":\"C:\\\\path\\\\to\\\\file\"}`,
			expectError: false,
		},
		{
			name:        "JSON array",
			input:       `[1,2,3,4,5]`,
			compact:     false,
			expected:    `[1,2,3,4,5]`,
			expectError: false,
		},
		{
			name: "JSON array with compact option",
			input: `[
  1,
  2,
  3
]`,
			compact:     true,
			expected:    `[1,2,3]`,
			expectError: false,
		},
		{
			name:        "Empty JSON object",
			input:       `{}`,
			compact:     false,
			expected:    `{}`,
			expectError: false,
		},
		{
			name:        "Empty JSON array",
			input:       `[]`,
			compact:     false,
			expected:    `[]`,
			expectError: false,
		},
		{
			name:        "Boolean true value",
			input:       `true`,
			compact:     false,
			expected:    `true`,
			expectError: false,
		},
		{
			name:        "Boolean false value",
			input:       `false`,
			compact:     false,
			expected:    `false`,
			expectError: false,
		},
		{
			name:        "Null value",
			input:       `null`,
			compact:     false,
			expected:    `null`,
			expectError: false,
		},
		{
			name:        "Number value",
			input:       `42`,
			compact:     false,
			expected:    `42`,
			expectError: false,
		},
		{
			name:        "Deeply nested JSON object",
			input:       `{"level1":{"level2":{"level3":{"level4":{"level5":"deep"}}}}}`,
			compact:     false,
			expected:    `{\"level1\":{\"level2\":{\"level3\":{\"level4\":{\"level5\":\"deep\"}}}}}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON - missing value",
			input:       `{"name":"John", "age":}`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - missing quote",
			input:       `{"name:"John", "age":30}`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - missing comma",
			input:       `{"name":"John" "age":30}`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - invalid value",
			input:       `{"name":"John", "age":"thirty"}`,
			compact:     false,
			expected:    `{\"name\":\"John\", \"age\":\"thirty\"}`,
			expectError: false, // This is valid JSON though not the expected type
		},
		{
			name:        "Invalid JSON - incomplete object",
			input:       `{"name":"John"`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - missing closing bracket",
			input:       `{"array": [1, 2, 3 }`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - missing opening bracket",
			input:       `{"array": 1, 2, 3] }`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON - extra comma",
			input:       `{"name":"John", "age":30,}`,
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Empty input",
			input:       "",
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Whitespace only",
			input:       "   \n\t",
			compact:     false,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Non-JSON input",
			input:       "This is not JSON",
			compact:     false,
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Encode([]byte(tc.input), tc.compact)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if result != tc.expected {
					t.Errorf("expected %s but got %s", tc.expected, result)
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		pretty      bool
		validate    func(string) bool
		expectError bool
	}{
		{
			name:   "Simple escaped JSON",
			input:  `{\"name\":\"John\",\"age\":30}`,
			pretty: false,
			validate: func(result string) bool {
				return result == `{"name":"John","age":30}` || result == `{"age":30,"name":"John"}`
			},
			expectError: false,
		},
		{
			name:   "Escaped JSON with pretty output",
			input:  `{\"name\":\"John\",\"age\":30}`,
			pretty: true,
			validate: func(result string) bool {
				// Check if result is properly indented
				var obj map[string]interface{}
				if err := json.Unmarshal([]byte(result), &obj); err != nil {
					return false
				}
				return strings.Contains(result, "\n") && strings.Contains(result, "  ")
			},
			expectError: false,
		},
		{
			name:   "Escaped JSON with newlines",
			input:  `{\n  \"name\": \"John\",\n  \"age\": 30\n}`,
			pretty: false,
			validate: func(result string) bool {
				var obj map[string]interface{}
				return json.Unmarshal([]byte(result), &obj) == nil
			},
			expectError: false,
		},
		{
			name:   "Escaped JSON array",
			input:  `[\"one\",\"two\",\"three\"]`,
			pretty: false,
			validate: func(result string) bool {
				return result == `["one","two","three"]`
			},
			expectError: false,
		},
		{
			name:   "Escaped JSON array with pretty output",
			input:  `[\"one\",\"two\",\"three\"]`,
			pretty: true,
			validate: func(result string) bool {
				return strings.Contains(result, "[\n") && strings.Contains(result, "  \"one\"")
			},
			expectError: false,
		},
		{
			name:   "Escaped boolean value",
			input:  `true`,
			pretty: false,
			validate: func(result string) bool {
				return result == "true"
			},
			expectError: false,
		},
		{
			name:   "Escaped null value",
			input:  `null`,
			pretty: false,
			validate: func(result string) bool {
				return result == "null"
			},
			expectError: false,
		},
		{
			name:   "Escaped number value",
			input:  `42`,
			pretty: false,
			validate: func(result string) bool {
				return result == "42"
			},
			expectError: false,
		},
		{
			name:   "Deeply nested escaped JSON",
			input:  `{\"level1\":{\"level2\":{\"level3\":{\"level4\":{\"level5\":\"deep\"}}}}}`,
			pretty: false,
			validate: func(result string) bool {
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(result), &obj)
				return err == nil && strings.Contains(result, `"level1"`)
			},
			expectError: false,
		},
		{
			name:   "Deeply nested escaped JSON with pretty output",
			input:  `{\"level1\":{\"level2\":{\"level3\":{\"level4\":{\"level5\":\"deep\"}}}}}`,
			pretty: true,
			validate: func(result string) bool {
				return strings.Contains(result, "{\n") && strings.Contains(result, "  \"level1\"")
			},
			expectError: false,
		},
		{
			name:        "Invalid escaped JSON - syntax error",
			input:       `{\"name\":\"John\",\"age\":}`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - missing quote",
			input:       `{\"name:\"John\",\"age\":30}`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - missing escape",
			input:       `{"name":"John","age":30}`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - incomplete object",
			input:       `{\"name\":\"John\"`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - broken escape sequence",
			input:       `{\\\"name\":\"John\"}`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - missing closing bracket",
			input:       `{\"array\": [1, 2, 3 }`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - missing opening bracket",
			input:       `{\"array\": 1, 2, 3] }`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid escaped JSON - extra comma",
			input:       `{\"name\":\"John\", \"age\":30,}`,
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Invalid pretty marshal",
			input:       `{\"name\":\"John\"}`,
			pretty:      true,
			validate:    func(string) bool { return true },
			expectError: false, // Should not fail with a valid JSON
		},
		{
			name:        "Not a JSON string",
			input:       "This is not JSON",
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Empty input",
			input:       "",
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
		{
			name:        "Whitespace only",
			input:       "   \n\t",
			pretty:      false,
			validate:    func(string) bool { return true },
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Decode([]byte(tc.input), tc.pretty)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if !tc.validate(result) {
					t.Errorf("validation failed for result: %s", result)
				}
			}
		})
	}
}

// Test for JSON marshal edge cases
func TestEncodeMarshalEdgeCases(t *testing.T) {
	// Test with valid JSON but with Unicode characters that need escaping
	t.Run("Unicode characters", func(t *testing.T) {
		input := []byte(`{"message":"こんにちは世界"}`)
		result, err := Encode(input, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(result, "こんにちは世界") {
			t.Errorf("expected unicode characters in result but got %s", result)
		}
	})

	// Test with valid JSON containing emojis
	t.Run("Emoji characters", func(t *testing.T) {
		input := []byte(`{"emoji":"😀🙈🚀"}`)
		result, err := Encode(input, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(result, "😀🙈🚀") {
			t.Errorf("expected emoji characters in result but got %s", result)
		}
	})

	// Test with very long string
	t.Run("Long string", func(t *testing.T) {
		// Create a long string
		longStr := strings.Repeat("a", 10000)
		input := []byte(`{"longString":"` + longStr + `"}`)
		result, err := Encode(input, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) < 10000 {
			t.Errorf("expected long result but got length %d", len(result))
		}
	})
}

// Test for JSON unmarshal edge cases
func TestDecodeUnmarshalEdgeCases(t *testing.T) {
	// Test with escaped Unicode characters
	t.Run("Escaped Unicode characters", func(t *testing.T) {
		input := []byte(`{\"message\":\"こんにちは世界\"}`)
		result, err := Decode(input, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(result, "こんにちは世界") {
			t.Errorf("expected unicode characters in result but got %s", result)
		}
	})

	// Test with escaped Emoji characters
	t.Run("Escaped Emoji characters", func(t *testing.T) {
		input := []byte(`{\"emoji\":\"😀🙈🚀\"}`)
		result, err := Decode(input, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(result, "😀🙈🚀") {
			t.Errorf("expected emoji characters in result but got %s", result)
		}
	})

	// Test with very long escaped string
	t.Run("Long escaped string", func(t *testing.T) {
		// Create a long escaped string with proper escaping
		longStr := strings.Repeat("a", 1000)
		input := []byte(`{\"longString\":\"` + longStr + `\"}`)
		_, err := Decode(input, false)
		// We don't care about the result here, just that it doesn't panic or fail
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// Test both functions with nil input
func TestNilInput(t *testing.T) {
	t.Run("Encode with nil input", func(t *testing.T) {
		_, err := Encode(nil, false)
		if err == nil {
			t.Errorf("expected error with nil input but got none")
		}
	})

	t.Run("Decode with nil input", func(t *testing.T) {
		_, err := Decode(nil, false)
		if err == nil {
			t.Errorf("expected error with nil input but got none")
		}
	})
}

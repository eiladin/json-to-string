package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLI performs integration tests on the CLI application
func TestCLI(t *testing.T) {
	// Skip if running short tests
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write test JSON to the file
	testJSON := `{"name":"From File","value":42}`
	if _, err := tmpFile.Write([]byte(testJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Build the binary for testing
	binaryPath := filepath.Join(t.TempDir(), "json-to-string-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		input          string
		validateOutput func(string) bool
		expectError    bool
	}{
		{
			name:  "Encode JSON input from stdin",
			args:  []string{},
			input: `{"name":"John","age":30}`,
			validateOutput: func(output string) bool {
				return strings.Contains(output, `{\"name\":\"John\",\"age\":30}`) ||
					strings.Contains(output, `{\"age\":30,\"name\":\"John\"}`)
			},
			expectError: false,
		},
		{
			name:  "Encode JSON input with compact flag",
			args:  []string{"--compact"},
			input: "{\n  \"name\": \"John\",\n  \"age\": 30\n}",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `{\"age\":30,\"name\":\"John\"}`) ||
					strings.Contains(output, `{\"name\":\"John\",\"age\":30}`)
			},
			expectError: false,
		},
		{
			name:  "Encode JSON from json flag",
			args:  []string{"--json", `{"name":"John","age":30}`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `{\"name\":\"John\",\"age\":30}`) ||
					strings.Contains(output, `{\"age\":30,\"name\":\"John\"}`)
			},
			expectError: false,
		},
		{
			name:  "Decode JSON string",
			args:  []string{"--decode", "--json", `{\"name\":\"John\",\"age\":30}`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `{"name":"John","age":30}`) ||
					strings.Contains(output, `{"age":30,"name":"John"}`)
			},
			expectError: false,
		},
		{
			name:  "Decode JSON string with pretty flag",
			args:  []string{"--decode", "--pretty", "--json", `{\"name\":\"John\",\"age\":30}`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, "{\n  ") &&
					(strings.Contains(output, `"name": "John"`) ||
						strings.Contains(output, `"age": 30`))
			},
			expectError: false,
		},
		{
			name:  "Raw output",
			args:  []string{"--raw", "--json", `{"name":"John"}`},
			input: "",
			validateOutput: func(output string) bool {
				return output == `{\"name\":\"John\"}`
			},
			expectError: false,
		},
		{
			name:  "Encode array input",
			args:  []string{"--json", `[1,2,3,4,5]`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `[1,2,3,4,5]`)
			},
			expectError: false,
		},
		{
			name:  "Encode boolean input",
			args:  []string{"--json", `true`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `true`)
			},
			expectError: false,
		},
		{
			name:  "Encode null input",
			args:  []string{"--json", `null`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `null`)
			},
			expectError: false,
		},
		{
			name:  "Encode empty object",
			args:  []string{"--json", `{}`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `{}`)
			},
			expectError: false,
		},
		{
			name:  "Encode empty array",
			args:  []string{"--json", `[]`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, `[]`)
			},
			expectError: false,
		},
		{
			name:  "Decode with raw flag",
			args:  []string{"--decode", "--raw", "--json", `{\"name\":\"John\"}`},
			input: "",
			validateOutput: func(output string) bool {
				return output == `{"name":"John"}`
			},
			expectError: false,
		},
		{
			name:  "Decode array with pretty flag",
			args:  []string{"--decode", "--pretty", "--json", `[\"one\",\"two\",\"three\"]`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, "[\n") && strings.Contains(output, "  \"one\"")
			},
			expectError: false,
		},
		// Error cases
		{
			name:  "Invalid JSON input - syntax error",
			args:  []string{"--json", `{"name":"John", "age":}`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid JSON input - missing quotes",
			args:  []string{"--json", `{name:John}`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid decode input - not JSON string",
			args:  []string{"--decode", "--json", "This is not JSON"},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid decode input - missing escapes",
			args:  []string{"--decode", "--json", `{"name":"John"}`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Empty input with decode",
			args:  []string{"--decode", "--json", ""},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid flags combination - compact with decode",
			args:  []string{"--decode", "--compact", "--json", `{\"name\":\"John\"}`},
			input: "",
			validateOutput: func(output string) bool {
				return true // The compact flag is ignored when decode is used
			},
			expectError: false, // Not an error, just ignores the compact flag
		},
		{
			name:  "Invalid JSON input - missing closing bracket",
			args:  []string{"--json", `{"array": [1, 2, 3 }`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid JSON input - missing opening bracket",
			args:  []string{"--json", `{"array": 1, 2, 3] }`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Invalid JSON input - extra comma",
			args:  []string{"--json", `{"name":"John", "age":30,}`},
			input: "",
			validateOutput: func(output string) bool {
				return true
			},
			expectError: true,
		},
		{
			name:  "Conflicting input sources - file and json",
			args:  []string{"--file", tmpPath, "--json", `{"name":"John"}`},
			input: "",
			validateOutput: func(output string) bool {
				return strings.Contains(output, "From File") // File takes precedence over JSON string
			},
			expectError: false, // Not an error, just prioritizes one over the other
		},
		{
			name:  "Conflicting input sources - stdin and json",
			args:  []string{"--json", `{"name":"John"}`},
			input: `{"other":"data"}`,
			validateOutput: func(output string) bool {
				return strings.Contains(output, "John") // JSON flag should take precedence
			},
			expectError: false, // Not an error, just prioritizes one over the other
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tc.args...)

			// If we have input, set it up
			if tc.input != "" {
				stdin, err := cmd.StdinPipe()
				if err != nil {
					t.Fatalf("Failed to get stdin pipe: %v", err)
				}
				go func() {
					defer stdin.Close()
					io.WriteString(stdin, tc.input)
				}()
			}

			// Capture output
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Run command
			err := cmd.Run()

			// Check for expected error state
			if tc.expectError {
				if err == nil && stderr.Len() == 0 {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v, stderr: %s", err, stderr.String())
				}

				output := stdout.String()
				// Trim newline for comparison if needed
				output = strings.TrimSpace(output)

				if !tc.validateOutput(output) {
					t.Errorf("output validation failed for: %q", output)
				}
			}
		})
	}
}

// Test error cases for file handling
func TestFileErrors(t *testing.T) {
	// Skip if running short tests
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary for testing
	binaryPath := filepath.Join(t.TempDir(), "json-to-string-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--file", "non-existent-file.json")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err == nil {
			t.Errorf("expected error with non-existent file but got none")
		}
		if !strings.Contains(stderr.String(), "Error reading file") {
			t.Errorf("expected 'Error reading file' in stderr but got: %s", stderr.String())
		}
	})

	// Test with invalid JSON file
	t.Run("Invalid JSON file", func(t *testing.T) {
		// Create a temporary file with invalid JSON
		tmpFile, err := os.CreateTemp("", "invalid-json-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		// Write invalid JSON to the file
		if _, err := tmpFile.Write([]byte(`{"name":"Test", "value":}`)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		cmd := exec.Command(binaryPath, "--file", tmpPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err == nil {
			t.Errorf("expected error with invalid JSON file but got none")
		}
		if !strings.Contains(stderr.String(), "Error") {
			t.Errorf("expected error message in stderr but got: %s", stderr.String())
		}
	})

	// Test with empty file
	t.Run("Empty file", func(t *testing.T) {
		// Create a temporary empty file
		tmpFile, err := os.CreateTemp("", "empty-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		cmd := exec.Command(binaryPath, "--file", tmpPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err == nil {
			t.Errorf("expected error with empty file but got none")
		}
	})

	// Test with readable directory
	t.Run("Directory as file", func(t *testing.T) {
		// Create a temporary directory
		tmpDir, err := os.MkdirTemp("", "test-dir-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		cmd := exec.Command(binaryPath, "--file", tmpDir)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err == nil {
			t.Errorf("expected error when using directory as file but got none")
		}
	})

	// Test with invalid permissions
	t.Run("File with no read permissions", func(t *testing.T) {
		// Skip this test on Windows as permission handling is different
		if os.Getenv("OS") == "Windows_NT" {
			t.Skip("Skipping permission test on Windows")
		}

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "noperm-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		// Write valid JSON to the file
		if _, err := tmpFile.Write([]byte(`{"name":"Test"}`)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		// Remove read permissions
		if err := os.Chmod(tmpPath, 0); err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}

		cmd := exec.Command(binaryPath, "--file", tmpPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		// We don't check the specific error as it might vary by OS
		if err == nil {
			t.Errorf("expected error with no read permissions but got none")
		}
	})

	// Test with file that's too large (simulate by checking if error handling works)
	t.Run("Very large JSON file decode", func(t *testing.T) {
		// Create a temporary file with a large amount of escaped JSON
		tmpFile, err := os.CreateTemp("", "large-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		// Write a large amount of escaped JSON to the file
		largeJSON := `{\"data\":\"`
		for i := 0; i < 10000; i++ {
			largeJSON += "test"
		}
		largeJSON += `\"}`

		if _, err := tmpFile.Write([]byte(largeJSON)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		// Test decoding the large file
		cmd := exec.Command(binaryPath, "--decode", "--file", tmpPath)
		err = cmd.Run()
		// We're not checking for a specific error, just making sure it handles large input
		// This might not actually fail but should test the handling of large input
	})
}

// Test various input format variations
func TestInputFormatVariations(t *testing.T) {
	// Skip if running short tests
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary for testing
	binaryPath := filepath.Join(t.TempDir(), "json-to-string-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Test with JSON array
	t.Run("JSON array input", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--json", `[1,2,3,4,5]`)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with JSON array: %v", err)
		}
		if !strings.Contains(stdout.String(), "[1,2,3,4,5]") {
			t.Errorf("expected array in output but got: %s", stdout.String())
		}
	})

	// Test with JSON boolean
	t.Run("JSON boolean input", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--json", `true`)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with JSON boolean: %v", err)
		}
		if !strings.Contains(stdout.String(), "true") {
			t.Errorf("expected boolean in output but got: %s", stdout.String())
		}
	})

	// Test with JSON number
	t.Run("JSON number input", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--json", `123.45`)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with JSON number: %v", err)
		}
		if !strings.Contains(stdout.String(), "123.45") {
			t.Errorf("expected number in output but got: %s", stdout.String())
		}
	})

	// Test with JSON null
	t.Run("JSON null input", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--json", `null`)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with JSON null: %v", err)
		}
		if !strings.Contains(stdout.String(), "null") {
			t.Errorf("expected null in output but got: %s", stdout.String())
		}
	})

	// Test with empty object
	t.Run("JSON empty object input", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--json", `{}`)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with empty object: %v", err)
		}
		if !strings.Contains(stdout.String(), "{}") {
			t.Errorf("expected empty object in output but got: %s", stdout.String())
		}
	})

	// Test with deeply nested JSON
	t.Run("Deeply nested JSON input", func(t *testing.T) {
		nestedJSON := `{"level1":{"level2":{"level3":{"level4":{"level5":"deep"}}}}}`
		cmd := exec.Command(binaryPath, "--json", nestedJSON)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with nested JSON: %v", err)
		}
		if !strings.Contains(stdout.String(), "level5") {
			t.Errorf("expected nested JSON in output but got: %s", stdout.String())
		}
	})

	// Test with Unicode characters
	t.Run("Unicode characters in JSON", func(t *testing.T) {
		unicodeJSON := `{"message":"こんにちは世界"}`
		cmd := exec.Command(binaryPath, "--json", unicodeJSON)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with Unicode JSON: %v", err)
		}
		if !strings.Contains(stdout.String(), "こんにちは世界") {
			t.Errorf("expected Unicode characters in output but got: %s", stdout.String())
		}
	})

	// Test with emojis
	t.Run("Emoji characters in JSON", func(t *testing.T) {
		emojiJSON := `{"emoji":"😀🙈🚀"}`
		cmd := exec.Command(binaryPath, "--json", emojiJSON)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			t.Errorf("error running with emoji JSON: %v", err)
		}
		if !strings.Contains(stdout.String(), "😀🙈🚀") {
			t.Errorf("expected emoji characters in output but got: %s", stdout.String())
		}
	})
}

// TestTempFile tests the file input functionality
func TestTempFile(t *testing.T) {
	// Skip if running short tests
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a temporary file with JSON
	tmpFile, err := os.CreateTemp("", "json-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write test JSON to the file
	testJSON := `{"name":"File Test","value":42}`
	if _, err := tmpFile.Write([]byte(testJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Build the binary for testing
	binaryPath := filepath.Join(t.TempDir(), "json-to-string-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Test file input
	cmd := exec.Command(binaryPath, "--file", tmpPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if !strings.Contains(output, `\"name\":\"File Test\"`) || !strings.Contains(output, `\"value\":42`) {
		t.Errorf("expected output to contain file test data but got %q", output)
	}
}

// Test no input case
func TestNoInput(t *testing.T) {
	// Skip if running short tests
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary for testing
	binaryPath := filepath.Join(t.TempDir(), "json-to-string-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	cmd := exec.Command(binaryPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err == nil {
		t.Errorf("expected error with no input but got none")
	}

	if !strings.Contains(stderr.String(), "No input provided") {
		t.Errorf("expected 'No input provided' in stderr but got: %s", stderr.String())
	}
}

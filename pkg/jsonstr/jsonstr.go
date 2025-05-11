package jsonstr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Encode takes a JSON byte slice and returns a properly escaped string representation
// If compact is true, it will remove newlines and extra whitespace from the input
func Encode(input []byte, compact bool) (string, error) {
	// Validate that the input is valid JSON
	var temp interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	var jsonStr string
	if compact {
		// If compact mode is enabled, re-marshal the JSON to remove formatting
		compactBytes, err := json.Marshal(temp)
		if err != nil {
			return "", fmt.Errorf("error compacting JSON: %w", err)
		}
		jsonStr = string(compactBytes)
	} else {
		jsonStr = string(input)
	}

	// Convert the JSON to a string with proper escaping
	result, err := json.Marshal(jsonStr)
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %w", err)
	}

	// The result is a JSON string, so we need to remove the outer quotes
	return strings.Trim(string(result), "\""), nil
}

// Decode takes an escaped JSON string and converts it back to JSON
// If pretty is true, it will format the output JSON with indentation
func Decode(input []byte, pretty bool) (string, error) {
	// First, we need to add quotes to make it a valid JSON string
	quotedInput := fmt.Sprintf("\"%s\"", string(input))

	// Unmarshal the string to get the actual JSON string with escapes interpreted
	var jsonString string
	if err := json.Unmarshal([]byte(quotedInput), &jsonString); err != nil {
		return "", fmt.Errorf("invalid JSON string: %w", err)
	}

	// Validate that the result is valid JSON
	var parsedJSON interface{}
	if err := json.Unmarshal([]byte(jsonString), &parsedJSON); err != nil {
		return "", fmt.Errorf("decoded string is not valid JSON: %w", err)
	}

	// Format the output according to the pretty flag
	if pretty {
		prettyBytes, err := json.MarshalIndent(parsedJSON, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error formatting JSON: %w", err)
		}
		return string(prettyBytes), nil
	}

	// Return the compact JSON
	compactBytes, err := json.Marshal(parsedJSON)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}
	return string(compactBytes), nil
}

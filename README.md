# JSON to String

A CLI tool that converts JSON into a properly escaped string representation, and vice versa.

## Installation

```bash
go install github.com/eiladin/json-to-string/cmd/json-to-string@latest
```

Or build from source:

```bash
git clone https://github.com/eiladin/json-to-string.git
cd json-to-string
go build -o json-to-string ./cmd/json-to-string
```

## Usage

### Encoding JSON to String

The tool accepts JSON input in several ways:

#### From a file:

```bash
json-to-string --file input.json
```

#### From a string argument:

```bash
json-to-string --json '{"key": "value"}'
```

#### From stdin (piping):

```bash
echo '{"key": "value"}' | json-to-string
```

#### Removing whitespace and newlines:

Use the `--compact` flag to remove formatting from pretty-printed JSON:

```bash
json-to-string --file input.json --compact
```

### Decoding String to JSON

Use the `--decode` flag to convert a JSON string back to JSON:

```bash
json-to-string --decode --json '{\"key\":\"value\"}'
```

#### Pretty-printing the output:

Use the `--pretty` flag with `--decode` to format the JSON output:

```bash
json-to-string --decode --pretty --file escaped.txt
```

### Raw Output

Use the `--raw` flag to output without a trailing newline (useful for piping):

```bash
echo '{"key":"value"}' | json-to-string --raw | json-to-string --decode --pretty
```

## Examples

### Encoding Example

Input JSON:
```json
{
  "name": "John",
  "age": 30,
  "city": "New York"
}
```

Default Output:
```
{\\n  \"name\": \"John\",\\n  \"age\": 30,\\n  \"city\": \"New York\"\\n}
```

Compact Output (with `--compact` flag):
```
{\"name\":\"John\",\"age\":30,\"city\":\"New York\"}
```

### Decoding Example

Input string:
```
{\"name\":\"John\",\"age\":30,\"city\":\"New York\"}
```

Default Output:
```
{"name":"John","age":30,"city":"New York"}
```

Pretty Output (with `--pretty` flag):
```
{
  "name": "John",
  "age": 30,
  "city": "New York"
}
```

### Piping Example

```bash
# Convert JSON to string and back
echo '{"test":"piping"}' | json-to-string --compact --raw | json-to-string --decode --pretty
```

## Development

### Building

```bash
# Build the application
make build

# Install locally
make install
```

### Testing

The project includes comprehensive unit and integration tests:

```bash
# Run all tests
make test

# Run unit tests only (skips integration tests)
make test-unit

# Run package tests only
make test-pkg

# Generate test coverage report
make test-coverage
```

The test suite includes:

- Unit tests for core functionality
- Integration tests for CLI behavior
- Extensive error condition tests
- Edge cases and invalid input handling
- File I/O error testing

Test coverage is over 85% for the core functionality.

### Other Commands

```bash
# Format Go code
make fmt

# Vet Go code
make vet

# Run all checks (fmt, vet, test)
make check

# Clean build artifacts
make clean
```

## Use Cases

This tool is useful for:
1. Preparing JSON for inclusion in source code
2. Converting escaped JSON strings back to valid JSON
3. Formatting JSON for better readability
4. Safely handling JSON with special characters and escapes 
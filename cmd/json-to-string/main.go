package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/eiladin/json-to-string/pkg/jsonstr"
)

// version is set during build
var version = "dev"

// printUsage prints a custom usage message with examples
func printUsage() {
	fmt.Fprintf(os.Stderr, "json-to-string - Convert JSON to escaped string format and vice versa\n\n")
	fmt.Fprintf(os.Stderr, "Version: %s\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string [options]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\nExamples:\n")

	// Encoding examples
	fmt.Fprintf(os.Stderr, "  Encoding (JSON to String):\n")
	fmt.Fprintf(os.Stderr, "  -----------------------\n")
	fmt.Fprintf(os.Stderr, "  # Encode JSON from a file:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --file input.json\n\n")

	fmt.Fprintf(os.Stderr, "  # Encode JSON from a string argument:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --json '{\"key\": \"value\"}'\n\n")

	fmt.Fprintf(os.Stderr, "  # Encode JSON from stdin (piping):\n")
	fmt.Fprintf(os.Stderr, "  echo '{\"key\": \"value\"}' | json-to-string\n\n")

	fmt.Fprintf(os.Stderr, "  # Encode JSON from file and remove whitespace from pretty-printed JSON:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --compact --file input.json\n\n")

	fmt.Fprintf(os.Stderr, "  # Encode without trailing newline (useful for piping):\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --file input.json --raw\n\n")

	// Decoding examples
	fmt.Fprintf(os.Stderr, "  Decoding (String to JSON):\n")
	fmt.Fprintf(os.Stderr, "  -----------------------\n")
	fmt.Fprintf(os.Stderr, "  # Decode an escaped JSON string back to JSON:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --decode --json '{\\\"key\\\":\\\"value\\\"}'\n\n")

	fmt.Fprintf(os.Stderr, "  # Decode and format the JSON output:\n")
	fmt.Fprintf(os.Stderr, "  json-to-string --decode --pretty --file escaped.txt\n\n")

	fmt.Fprintf(os.Stderr, "  # Chain encode and decode operations (pipe):\n")
	fmt.Fprintf(os.Stderr, "  echo '{\"key\":\"value\"}' | json-to-string --raw | json-to-string --decode --pretty\n")
}

func main() {
	var inputFile string
	var inputString string
	var compact bool
	var decode bool
	var pretty bool
	var rawOutput bool
	var showVersion bool
	var showHelp bool

	flag.StringVar(&inputFile, "file", "", "Input JSON file path")
	flag.StringVar(&inputString, "json", "", "JSON string input")
	flag.BoolVar(&compact, "compact", false, "Remove newlines and extra spaces from pretty-printed JSON")
	flag.BoolVar(&decode, "decode", false, "Decode an escaped JSON string back to JSON")
	flag.BoolVar(&pretty, "pretty", false, "Format the decoded JSON output with indentation (only used with --decode)")
	flag.BoolVar(&rawOutput, "raw", false, "Output without trailing newline (useful for piping)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help with examples")

	// Override the default usage function
	flag.Usage = printUsage

	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	var input []byte
	var err error

	switch {
	case inputFile != "":
		input, err = os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	case inputString != "":
		input = []byte(inputString)
	default:
		// Read from stdin if no file or string provided
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "No input provided. Use --file, --json or pipe data to stdin.")
			printUsage()
			os.Exit(1)
		}
	}

	var result string
	if decode {
		result, err = jsonstr.Decode(input, pretty)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding JSON string: %v\n", err)
			os.Exit(1)
		}
	} else {
		result, err = jsonstr.Encode(input, compact)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	}

	if rawOutput {
		fmt.Print(result)
	} else {
		fmt.Println(result)
	}
}

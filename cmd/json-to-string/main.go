package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/eiladin/json-to-string/pkg/jsonstr"
)

func main() {
	var inputFile string
	var inputString string
	var compact bool
	var decode bool
	var pretty bool
	var rawOutput bool

	flag.StringVar(&inputFile, "file", "", "Input JSON file path")
	flag.StringVar(&inputString, "json", "", "JSON string input")
	flag.BoolVar(&compact, "compact", false, "Remove newlines and extra spaces from pretty-printed JSON")
	flag.BoolVar(&decode, "decode", false, "Decode an escaped JSON string back to JSON")
	flag.BoolVar(&pretty, "pretty", false, "Format the decoded JSON output with indentation (only used with --decode)")
	flag.BoolVar(&rawOutput, "raw", false, "Output without trailing newline (useful for piping)")
	flag.Parse()

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
			flag.Usage()
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

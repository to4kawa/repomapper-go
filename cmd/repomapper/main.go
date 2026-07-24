package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/to4kawa/repomapper-go/internal/mapper"
)

func main() {
	maxTokens := flag.Int("tokens", 0, "max tokens for output (0 = no limit)")
	includeTests := flag.Bool("include-tests", false, "include test symbols in the map")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: repomapper [flags] <path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	absPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	text, err := mapper.Generate(mapper.Options{
		Path:         absPath,
		MaxTokens:    *maxTokens,
		IncludeTests: *includeTests,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Print(text)
}

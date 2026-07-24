package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/to4kawa/repomapper-go/internal/analyzer"
	"github.com/to4kawa/repomapper-go/internal/output"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: repomapper <path>")
		os.Exit(1)
	}

	repoPath := os.Args[1]
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	_, err = git.PlainOpen(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Not a git repository: %v\n", err)
		os.Exit(1)
	}

	a := analyzer.NewGoAnalyzer()
	symbols, err := a.AnalyzeDir(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Analyze error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output.FormatSymbols(symbols, absPath))
}

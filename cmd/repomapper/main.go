package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/to4kawa/repomapper-go/internal/analyzer"
	"github.com/to4kawa/repomapper-go/internal/output"
)

func main() {
	maxTokens := flag.Int("tokens", 0, "max tokens for output (0 = no limit)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: repomapper [flags] <path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	repoPath := flag.Arg(0)
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

	var allSymbols []analyzer.Symbol

	ts := analyzer.NewTreeSitterAnalyzer()
	if syms, err := ts.AnalyzeDir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "treesitter analyze warning: %v\n", err)
	} else {
		allSymbols = append(allSymbols, syms...)
	}

	rs := analyzer.NewRustAnalyzer()
	if syms, err := rs.AnalyzeDir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "rust analyze warning: %v\n", err)
	} else {
		allSymbols = append(allSymbols, syms...)
	}

	allSymbols = analyzer.Rank(allSymbols)

	mapText := output.FormatSymbols(allSymbols, absPath)
	if *maxTokens > 0 {
		mapText = output.LimitByTokens(mapText, *maxTokens)
	}

	fmt.Print(mapText)
}

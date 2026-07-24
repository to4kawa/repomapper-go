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

	var allSymbols []analyzer.Symbol

	// 1. gotreesitter（Go / Python / JS など）
	ts := analyzer.NewTreeSitterAnalyzer()
	tsSymbols, err := ts.AnalyzeDir(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "treesitter analyze warning: %v\n", err)
	} else {
		allSymbols = append(allSymbols, tsSymbols...)
	}

	// 2. Rust（外部バイナリ）
	rs := analyzer.NewRustAnalyzer()
	rsSymbols, err := rs.AnalyzeDir(absPath)
	if err != nil {
		// バイナリがない場合は警告だけ出して続行
		fmt.Fprintf(os.Stderr, "rust analyze warning: %v\n", err)
	} else {
		allSymbols = append(allSymbols, rsSymbols...)
	}

	allSymbols = analyzer.Rank(allSymbols)
	fmt.Print(output.FormatSymbols(allSymbols, absPath))
}

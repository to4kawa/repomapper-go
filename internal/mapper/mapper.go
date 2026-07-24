package mapper

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/to4kawa/repomapper-go/internal/analyzer"
	"github.com/to4kawa/repomapper-go/internal/output"
)

type Options struct {
	Path      string
	MaxTokens int // 0 = no limit
}

func Generate(opts Options) (string, error) {
	_, err := git.PlainOpen(opts.Path)
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}

	var all []analyzer.Symbol

	ts := analyzer.NewTreeSitterAnalyzer()
	if syms, err := ts.AnalyzeDir(opts.Path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: treesitter: %v\n", err)
	} else {
		all = append(all, syms...)
	}

	rs := analyzer.NewRustAnalyzer()
	if syms, err := rs.AnalyzeDir(opts.Path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: rust analyzer: %v\n", err)
	} else {
		all = append(all, syms...)
	}

	all = analyzer.Rank(all)

	text := output.FormatSymbols(all, opts.Path)
	if opts.MaxTokens > 0 {
		text = output.LimitByTokens(text, opts.MaxTokens)
	}
	return text, nil
}

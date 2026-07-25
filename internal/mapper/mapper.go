package mapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/to4kawa/repomapper-go/internal/analyzer"
	"github.com/to4kawa/repomapper-go/internal/output"
)

type Options struct {
	Path         string
	MaxTokens    int  // 0 = no limit
	IncludeTests bool // false = テスト除外（デフォルト）
}

func Generate(opts Options) (string, error) {
	_, err := git.PlainOpen(opts.Path)
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}

	var all []analyzer.Symbol

	// Tree-sitterアナライザー（汎用、Python以外）
	ts := analyzer.NewTreeSitterAnalyzer()
	if syms, err := ts.AnalyzeDir(opts.Path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: treesitter: %v\n", err)
	} else {
		all = append(all, syms...)
	}

	// Python固有アナライザー（汎用の結果を上書き）
	py := analyzer.NewPythonAnalyzer()
	if syms, err := py.AnalyzeDir(opts.Path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: python analyzer: %v\n", err)
	} else {
		// Pythonファイルの重複を除去してから追加
		all = mergePythonSymbols(all, syms)
	}

	// Rustアナライザー
	rs := analyzer.NewRustAnalyzer()
	if syms, err := rs.AnalyzeDir(opts.Path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: rust analyzer: %v\n", err)
	} else {
		all = append(all, syms...)
	}

	all = analyzer.Rank(all)

	if !opts.IncludeTests {
		all = filterOutTests(all)
	}

	text := output.FormatSymbols(all, opts.Path)
	if opts.MaxTokens > 0 {
		text = output.LimitByTokens(text, opts.MaxTokens)
	}
	return text, nil
}

// mergePythonSymbols は汎用アナライザーの結果とPython固有の結果をマージする
func mergePythonSymbols(generic, python []analyzer.Symbol) []analyzer.Symbol {
	// Pythonファイルのパスを収集
	pythonFiles := map[string]bool{}
	for _, s := range python {
		pythonFiles[s.File] = true
	}

	// 汎用アナライザーからPythonファイル以外を保持
	var result []analyzer.Symbol
	for _, s := range generic {
		if !pythonFiles[s.File] {
			result = append(result, s)
		}
	}

	// Python固有の結果を追加
	result = append(result, python...)
	return result
}

func filterOutTests(symbols []analyzer.Symbol) []analyzer.Symbol {
	out := make([]analyzer.Symbol, 0, len(symbols))
	for _, s := range symbols {
		if isTestSymbol(s) {
			continue
		}
		out = append(out, s)
	}
	return out
}

func isTestSymbol(s analyzer.Symbol) bool {
	file := filepath.ToSlash(s.File)
	base := filepath.Base(file)

	// Go: _test.go
	if strings.HasSuffix(base, "_test.go") {
		return true
	}
	// Rust: tests/ 配下、*_test.rs
	if strings.Contains(file, "/tests/") || strings.HasPrefix(file, "tests/") || strings.HasSuffix(base, "_test.rs") {
		return true
	}
	// 関数名ヒューリスティック: TestXxx / test_xxx / xxx_test
	name := s.Name
	if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test") {
		return true
	}
	return false
}

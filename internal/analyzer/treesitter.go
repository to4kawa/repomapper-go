package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

type TreeSitterAnalyzer struct{}

func NewTreeSitterAnalyzer() *TreeSitterAnalyzer {
	return &TreeSitterAnalyzer{}
}

func (a *TreeSitterAnalyzer) AnalyzeFile(path string) ([]Symbol, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	entry := grammars.DetectLanguage(path)
	if entry == nil {
		return nil, nil
	}

	lang := entry.Language()
	parser := gotreesitter.NewParser(lang)
	tree, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	defs := gotreesitter.ExtractDefinitionSpans(tree)

	var symbols []Symbol
	for _, d := range defs {
		line := byteOffsetToLine(src, d.StartByte)
		symbols = append(symbols, Symbol{
			Kind: mapKind(d.Kind),
			Name: d.Name,
			File: path,
			Line: line,
		})
	}

	return symbols, nil
}

func (a *TreeSitterAnalyzer) AnalyzeDir(root string) ([]Symbol, error) {
	var all []Symbol

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" || name == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		syms, err := a.AnalyzeFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			return nil
		}
		all = append(all, syms...)
		return nil
	})

	return all, err
}

func mapKind(k string) string {
	switch strings.ToLower(k) {
	case "function", "func", "method":
		return "func"
	case "struct", "class", "type", "enum":
		return "type"
	case "interface", "trait":
		return "interface"
	default:
		return k
	}
}

// バイトオフセットから行番号を計算（1-indexed）
func byteOffsetToLine(src []byte, offset uint32) int {
	if int(offset) > len(src) {
		offset = uint32(len(src))
	}
	line := 1
	for i := uint32(0); i < offset; i++ {
		if src[i] == '\n' {
			line++
		}
	}
	return line
}

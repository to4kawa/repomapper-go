package output

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/to4kawa/repomapper-go/internal/analyzer"
)

func FormatSymbols(symbols []analyzer.Symbol, root string) string {
	// ファイルごとにグループ化
	grouped := map[string][]analyzer.Symbol{}
	for _, s := range symbols {
		rel, err := filepath.Rel(root, s.File)
		if err != nil {
			rel = s.File
		}
		// Windowsのパス区切りを統一
		rel = filepath.ToSlash(rel)
		grouped[rel] = append(grouped[rel], s)
	}

	// ファイル名でソート
	files := make([]string, 0, len(grouped))
	for f := range grouped {
		files = append(files, f)
	}
	sort.Strings(files)

	var b strings.Builder
	for _, file := range files {
		b.WriteString(file)
		b.WriteString(":\n")

		syms := grouped[file]
		// 種類で軽くソート（type → func → method）
		sort.SliceStable(syms, func(i, j int) bool {
			order := map[string]int{"type": 0, "interface": 1, "func": 2, "method": 3}
			return order[syms[i].Kind] < order[syms[j].Kind]
		})

		for _, s := range syms {
			b.WriteString("  ")
			b.WriteString(formatSymbol(s))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func formatSymbol(s analyzer.Symbol) string {
	switch s.Kind {
	case "method":
		return fmt.Sprintf("func (%s) %s", s.Receiver, s.Name)
	case "func":
		return fmt.Sprintf("func %s", s.Name)
	case "type":
		return fmt.Sprintf("type %s", s.Name)
	case "interface":
		return fmt.Sprintf("type %s interface", s.Name)
	default:
		return fmt.Sprintf("%s %s", s.Kind, s.Name)
	}
}

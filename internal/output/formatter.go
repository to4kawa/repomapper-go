package output

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/to4kawa/repomapper-go/internal/analyzer"
)

func FormatSymbols(symbols []analyzer.Symbol, root string) string {
	// ファイルごとの出現順を保ったままグループ化
	type fileGroup struct {
		path    string
		symbols []analyzer.Symbol
	}

	order := []string{}
	grouped := map[string]*fileGroup{}

	for _, s := range symbols {
		rel, err := filepath.Rel(root, s.File)
		if err != nil {
			rel = s.File
		}
		rel = filepath.ToSlash(rel)

		if _, ok := grouped[rel]; !ok {
			grouped[rel] = &fileGroup{path: rel}
			order = append(order, rel)
		}
		grouped[rel].symbols = append(grouped[rel].symbols, s)
	}

	var b strings.Builder
	for _, path := range order {
		g := grouped[path]
		b.WriteString(path)
		b.WriteString(":\n")
		for _, s := range g.symbols {
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
		if s.Receiver != "" {
			return fmt.Sprintf("func (%s) %s", s.Receiver, s.Name)
		}
		return fmt.Sprintf("func %s", s.Name)
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

package analyzer

import (
	"sort"
	"unicode"
)

func Rank(symbols []Symbol) []Symbol {
	type scored struct {
		sym   Symbol
		score int
	}

	items := make([]scored, 0, len(symbols))
	for _, s := range symbols {
		items = append(items, scored{sym: s, score: scoreOf(s)})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].score > items[j].score
	})

	result := make([]Symbol, len(items))
	for i, it := range items {
		result[i] = it.sym
	}
	return result
}

func scoreOf(s Symbol) int {
	score := 0

	// 種類
	switch s.Kind {
	case "type", "interface":
		score += 5
	case "func", "method":
		score += 3
	default:
		score += 1
	}

	// 公開かどうか（簡易判定）
	if isExported(s.Name) {
		score += 10
	}

	return score
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)[0]
	// Go: 大文字 / Rust: 大文字始まりを公開とみなす（簡易）
	return unicode.IsUpper(r)
}

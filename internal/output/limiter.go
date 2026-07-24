package output

import (
	"strings"
	"unicode"
)

func EstimateTokens(s string) int {
	count := 0
	inToken := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			inToken = false
		} else {
			if !inToken {
				count++
				inToken = true
			}
		}
	}
	return count
}

func LimitByTokens(formatted string, maxTokens int) string {
	if maxTokens <= 0 {
		return formatted
	}

	lines := strings.Split(formatted, "\n")
	var result []string
	used := 0

	for _, line := range lines {
		t := EstimateTokens(line)
		if used+t > maxTokens {
			break
		}
		result = append(result, line)
		used += t
	}

	return strings.Join(result, "\n")
}

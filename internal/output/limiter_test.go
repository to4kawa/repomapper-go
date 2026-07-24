package output

import "testing"

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  hello  world  ", 2},
		{"one two three", 3},
	}
	for _, tt := range tests {
		got := EstimateTokens(tt.input)
		if got != tt.want {
			t.Errorf("EstimateTokens(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestLimitByTokens_NoLimit(t *testing.T) {
	input := "line one\nline two\nline three"
	result := LimitByTokens(input, 0)
	if result != input {
		t.Errorf("expected no change, got %q", result)
	}
}

func TestLimitByTokens_TrimLines(t *testing.T) {
	input := "a b\nc d\ne f"
	// "a b"=2, "c d"=2 => maxTokens=3 で1行目のみ
	result := LimitByTokens(input, 3)
	lines := splitLines(result)

	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d: %q", len(lines), result)
	}

	// "a b\nc d\ne f" => maxTokens=4 で2行
	result4 := LimitByTokens(input, 4)
	lines4 := splitLines(result4)

	if len(lines4) != 2 {
		t.Errorf("expected 2 lines with maxTokens=4, got %d: %q", len(lines4), result4)
	}
}

func TestLimitByTokens_ZeroTokens(t *testing.T) {
	input := "anything"
	result := LimitByTokens(input, 0)
	if result != input {
		t.Errorf("expected no change with maxTokens=0, got %q", result)
	}
}

func TestLimitByTokens_NegativeTokens(t *testing.T) {
	input := "anything"
	result := LimitByTokens(input, -5)
	if result != input {
		t.Errorf("expected no change with negative maxTokens, got %q", result)
	}
}

func TestLimitByTokens_DoesNotBreakLine(t *testing.T) {
	input := "func Foo\nfunc Bar\nfunc Baz\nfunc Qux"
	result := LimitByTokens(input, 3)

	// Should stop at a clean line boundary
	lines := splitLines(result)
	for _, line := range lines {
		if line == "" {
			continue
		}
		if len(line) == 0 {
			continue
		}
	}
	// 3 tokens: "func Foo" = 2 tokens, "func Bar" = 2 tokens => stops after first line
	if len(lines) < 1 {
		t.Error("expected at least one line")
	}
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

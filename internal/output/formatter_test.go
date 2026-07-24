package output

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/to4kawa/repomapper-go/internal/analyzer"
)

func TestFormatSymbols_GroupsByFile(t *testing.T) {
	symbols := []analyzer.Symbol{
		{Kind: "func", Name: "Main", File: filepath.Join("root", "cmd", "main.go"), Line: 1},
		{Kind: "func", Name: "Helper", File: filepath.Join("root", "cmd", "main.go"), Line: 5},
		{Kind: "type", Name: "Config", File: filepath.Join("root", "config.go"), Line: 1},
	}

	result := FormatSymbols(symbols, "root")

	if !strings.Contains(result, "cmd/main.go:") {
		t.Errorf("expected file header cmd/main.go, got:\n%s", result)
	}
	if !strings.Contains(result, "config.go:") {
		t.Errorf("expected file header config.go, got:\n%s", result)
	}
	if !strings.Contains(result, "func Main") {
		t.Errorf("expected func Main, got:\n%s", result)
	}
	if !strings.Contains(result, "func Helper") {
		t.Errorf("expected func Helper, got:\n%s", result)
	}
	if !strings.Contains(result, "type Config") {
		t.Errorf("expected type Config, got:\n%s", result)
	}
}

func TestFormatSymbols_Empty(t *testing.T) {
	result := FormatSymbols(nil, "root")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestFormatSymbols_PreservesOrder(t *testing.T) {
	symbols := []analyzer.Symbol{
		{Kind: "func", Name: "Zebra", File: filepath.Join("root", "a.go"), Line: 1},
		{Kind: "func", Name: "Alpha", File: filepath.Join("root", "a.go"), Line: 2},
	}

	result := FormatSymbols(symbols, "root")
	zIdx := strings.Index(result, "Zebra")
	aIdx := strings.Index(result, "Alpha")

	if zIdx < 0 || aIdx < 0 || zIdx >= aIdx {
		t.Errorf("expected Zebra before Alpha, got:\n%s", result)
	}
}

func TestFormatSymbols_MethodWithReceiver(t *testing.T) {
	symbols := []analyzer.Symbol{
		{Kind: "method", Name: "Serve", Receiver: "*Server", File: filepath.Join("root", "s.go"), Line: 1},
	}

	result := FormatSymbols(symbols, "root")
	if !strings.Contains(result, "func (*Server) Serve") {
		t.Errorf("expected method with receiver, got:\n%s", result)
	}
}

func TestFormatSymbols_Interface(t *testing.T) {
	symbols := []analyzer.Symbol{
		{Kind: "interface", Name: "Reader", File: filepath.Join("root", "r.go"), Line: 1},
	}

	result := FormatSymbols(symbols, "root")
	if !strings.Contains(result, "type Reader interface") {
		t.Errorf("expected interface format, got:\n%s", result)
	}
}

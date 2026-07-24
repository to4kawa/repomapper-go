package mapper

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/to4kawa/repomapper-go/internal/analyzer"
)

func TestGenerate_BasicGoRepo(t *testing.T) {
	dir := t.TempDir()

	// git init
	if out, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	// Create a Go file with known symbols
	goFile := filepath.Join(dir, "app.go")
	content := `package main

type Server struct{}

func (s *Server) Listen() {}

func Run() {}

func internalHelper() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Generate(Options{Path: dir})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !strings.Contains(result, "type Server") {
		t.Errorf("expected type Server in output:\n%s", result)
	}
	if !strings.Contains(result, "func Run") {
		t.Errorf("expected func Run in output:\n%s", result)
	}
	if !strings.Contains(result, "app.go:") {
		t.Errorf("expected file header app.go:\n%s", result)
	}
}

func TestGenerate_MaxTokens(t *testing.T) {
	dir := t.TempDir()

	if out, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	goFile := filepath.Join(dir, "a.go")
	content := `package main

func Alpha() {}
func Beta() {}
func Gamma() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	full, err := Generate(Options{Path: dir})
	if err != nil {
		t.Fatal(err)
	}

	limited, err := Generate(Options{Path: dir, MaxTokens: 3})
	if err != nil {
		t.Fatal(err)
	}

	if len(limited) >= len(full) {
		t.Errorf("expected limited output to be shorter: full=%d limited=%d", len(full), len(limited))
	}
}

func TestGenerate_NotAGitRepo(t *testing.T) {
	dir := t.TempDir()

	_, err := Generate(Options{Path: dir})
	if err == nil {
		t.Error("expected error for non-git directory")
	}
}

func TestGenerate_EmptyRepo(t *testing.T) {
	dir := t.TempDir()

	if out, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	result, err := Generate(Options{Path: dir})
	if err != nil {
		t.Fatalf("Generate failed on empty repo: %v", err)
	}

	// Empty repo should return empty or whitespace-only string
	if strings.TrimSpace(result) != "" {
		t.Errorf("expected empty output for empty repo, got:\n%s", result)
	}
}

func TestGenerate_IncludeTests(t *testing.T) {
	dir := t.TempDir()

	if out, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	// テストファイル
	testFile := filepath.Join(dir, "app_test.go")
	testContent := `package main

func TestHello(t *testing.T) {}
func TestWorld(t *testing.T) {}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 本体ファイル
	mainFile := filepath.Join(dir, "app.go")
	mainContent := `package main

func Run() {}
`
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	// デフォルト（テスト除外）
	outDefault, err := Generate(Options{Path: dir})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(outDefault, "TestHello") {
		t.Errorf("default should exclude test symbols, got:\n%s", outDefault)
	}

	// テスト込み
	outWithTests, err := Generate(Options{Path: dir, IncludeTests: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(outWithTests, "TestHello") {
		t.Errorf("IncludeTests=true should include test symbols, got:\n%s", outWithTests)
	}
}

func TestFilterOutTests(t *testing.T) {
	symbols := []analyzer.Symbol{
		{Kind: "func", Name: "Run", File: "app.go", Line: 1},
		{Kind: "func", Name: "TestHello", File: "app_test.go", Line: 1},
		{Kind: "func", Name: "helper", File: "tests/helper.rs", Line: 1},
		{Kind: "func", Name: "test_something", File: "lib.rs", Line: 1},
		{Kind: "func", Name: "real_func", File: "lib.rs", Line: 5},
		{Kind: "type", Name: "Config", File: "config.go", Line: 1},
	}

	result := filterOutTests(symbols)

	expected := []string{"Run", "real_func", "Config"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d symbols, got %d", len(expected), len(result))
	}
	for i, name := range expected {
		if result[i].Name != name {
			t.Errorf("index %d: expected %s, got %s", i, name, result[i].Name)
		}
	}
}

func TestIsTestSymbol(t *testing.T) {
	tests := []struct {
		sym  analyzer.Symbol
		want bool
	}{
		{analyzer.Symbol{Name: "Run", File: "app.go"}, false},
		{analyzer.Symbol{Name: "TestHello", File: "app.go"}, true},
		{analyzer.Symbol{Name: "Run", File: "app_test.go"}, true},
		{analyzer.Symbol{Name: "helper", File: "tests/helper.rs"}, true},
		{analyzer.Symbol{Name: "foo_test", File: "lib.rs"}, true},
		{analyzer.Symbol{Name: "test_config", File: "lib.rs"}, true},
		{analyzer.Symbol{Name: "real_func", File: "lib.rs"}, false},
	}
	for _, tt := range tests {
		got := isTestSymbol(tt.sym)
		if got != tt.want {
			t.Errorf("isTestSymbol(%+v) = %v, want %v", tt.sym, got, tt.want)
		}
	}
}

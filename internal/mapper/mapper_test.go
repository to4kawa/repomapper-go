package mapper

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

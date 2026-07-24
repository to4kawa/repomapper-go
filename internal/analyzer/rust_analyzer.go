package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RustAnalyzer struct {
	binPath string
}

func NewRustAnalyzer() *RustAnalyzer {
	// 優先順位: 環境変数 → 同じディレクトリ → PATH
	bin := os.Getenv("REPOMAPPER_RUST_BIN")
	if bin == "" {
		// 開発中は相対パスを試す
		candidates := []string{
			"./rust-analyzer/target/debug/repomapper-rust",
			"./rust-analyzer/target/release/repomapper-rust",
			"repomapper-rust",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				bin = c
				break
			}
			// Windows対応
			if _, err := os.Stat(c + ".exe"); err == nil {
				bin = c + ".exe"
				break
			}
		}
	}
	return &RustAnalyzer{binPath: bin}
}

func (a *RustAnalyzer) AnalyzeDir(root string) ([]Symbol, error) {
	if a.binPath == "" {
		return nil, fmt.Errorf("repomapper-rust binary not found")
	}

	cmd := exec.Command(a.binPath, root)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("rust analyzer failed: %w", err)
	}

	var raw []struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
		File string `json:"file"`
		Line int    `json:"line"`
	}
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("invalid json from rust analyzer: %w", err)
	}

	symbols := make([]Symbol, 0, len(raw))
	for _, r := range raw {
		// 相対パスを絶対パスに寄せる
		file := r.File
		if !filepath.IsAbs(file) {
			file = filepath.Join(root, file)
		}
		symbols = append(symbols, Symbol{
			Kind: r.Kind,
			Name: r.Name,
			File: file,
			Line: r.Line,
		})
	}
	return symbols, nil
}

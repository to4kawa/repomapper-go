package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
)

type PythonAnalyzer struct{}

func NewPythonAnalyzer() *PythonAnalyzer {
	return &PythonAnalyzer{}
}

func (a *PythonAnalyzer) AnalyzeFile(path string) ([]Symbol, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	entry := grammars.DetectLanguage(path)
	if entry == nil {
		return nil, nil
	}

	// Python以外はスキップ
	lang := entry.Language()
	if lang == nil {
		return nil, nil
	}

	parser := gotreesitter.NewParser(lang)
	tree, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	var symbols []Symbol
	symbols = append(symbols, extractPythonDefinitions(src, tree.RootNode(), path, 0, lang)...)

	return symbols, nil
}

func (a *PythonAnalyzer) AnalyzeDir(root string) ([]Symbol, error) {
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

// extractPythonDefinitions はASTを走査してPythonの定義を抽出する
func extractPythonDefinitions(src []byte, node *gotreesitter.Node, path string, depth int, lang *gotreesitter.Language) []Symbol {
	var symbols []Symbol

	// トップレベルのみ処理（depth > 0 はネストされた定義）
	if depth > 0 {
		return nil
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		nodeType := child.Type(lang)

		switch nodeType {
		case "function_definition":
			sym := extractPythonFunction(src, child, path, lang)
			if sym != nil {
				symbols = append(symbols, *sym)
			}
		case "class_definition":
			classSyms := extractPythonClass(src, child, path, lang)
			symbols = append(symbols, classSyms...)
		case "assignment", "augmented_assignment":
			sym := extractPythonVariable(src, child, path, lang)
			if sym != nil {
				symbols = append(symbols, *sym)
			}
		case "decorated_definition":
			// デコレータ付き定義は子ノードを処理
			decoratedSyms := extractDecoratedDefinition(src, child, path, lang)
			symbols = append(symbols, decoratedSyms...)
		}
	}

	return symbols
}

// extractPythonFunction は関数定義からシンボルを抽出する
func extractPythonFunction(src []byte, node *gotreesitter.Node, path string, lang *gotreesitter.Language) *Symbol {
	nameNode := node.ChildByFieldName("name", lang)
	if nameNode == nil {
		return nil
	}

	name := nodeContent(src, nameNode)
	line := byteOffsetToLine(src, node.StartByte())

	return &Symbol{
		Kind: "func",
		Name: name,
		File: path,
		Line: line,
	}
}

// extractPythonClass はクラス定義からシンボルを抽出する
func extractPythonClass(src []byte, node *gotreesitter.Node, path string, lang *gotreesitter.Language) []Symbol {
	var symbols []Symbol

	nameNode := node.ChildByFieldName("name", lang)
	if nameNode == nil {
		return nil
	}

	className := nodeContent(src, nameNode)
	line := byteOffsetToLine(src, node.StartByte())

	// クラス自体を追加
	symbols = append(symbols, Symbol{
		Kind: "type",
		Name: className,
		File: path,
		Line: line,
	})

	// クラスボディ内のメソッドを抽出
	bodyNode := node.ChildByFieldName("body", lang)
	if bodyNode != nil {
		for i := 0; i < int(bodyNode.ChildCount()); i++ {
			child := bodyNode.Child(i)
			if child == nil {
				continue
			}

			childType := child.Type(lang)
			if childType == "function_definition" {
				sym := extractPythonMethod(src, child, path, className, lang)
				if sym != nil {
					symbols = append(symbols, *sym)
				}
			} else if childType == "decorated_definition" {
				methodSyms := extractDecoratedDefinition(src, child, path, lang)
				symbols = append(symbols, methodSyms...)
			}
		}
	}

	return symbols
}

// extractPythonMethod はメソッド定義からシンボルを抽出する
func extractPythonMethod(src []byte, node *gotreesitter.Node, path string, className string, lang *gotreesitter.Language) *Symbol {
	nameNode := node.ChildByFieldName("name", lang)
	if nameNode == nil {
		return nil
	}

	name := nodeContent(src, nameNode)
	line := byteOffsetToLine(src, node.StartByte())

	return &Symbol{
		Kind:     "method",
		Name:     name,
		Receiver: className,
		File:     path,
		Line:     line,
	}
}

// extractPythonVariable は代入文からグローバル変数を抽出する
func extractPythonVariable(src []byte, node *gotreesitter.Node, path string, lang *gotreesitter.Language) *Symbol {
	// 大文字の定数のみ抽出（PRINTER_NAME等）
	leftNode := node.Child(0)
	if leftNode == nil {
		return nil
	}

	name := nodeContent(src, leftNode)

	// 大文字で構成される定数のみ対象
	if !isPythonConstant(name) {
		return nil
	}

	line := byteOffsetToLine(src, node.StartByte())

	return &Symbol{
		Kind: "constant",
		Name: name,
		File: path,
		Line: line,
	}
}

// extractDecoratedDefinition はデコレータ付き定義を処理する
func extractDecoratedDefinition(src []byte, node *gotreesitter.Node, path string, lang *gotreesitter.Language) []Symbol {
	var symbols []Symbol
	var annotations []string

	// デコレータを抽出
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type(lang) == "decorator" {
			annotation := extractDecorator(src, child)
			if annotation != "" {
				annotations = append(annotations, annotation)
			}
		}
	}

	// 子ノードから関数/クラス定義を探す
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		childType := child.Type(lang)
		switch childType {
		case "function_definition":
			sym := extractPythonFunction(src, child, path, lang)
			if sym != nil {
				sym.Annotations = annotations
				symbols = append(symbols, *sym)
			}
		case "class_definition":
			classSyms := extractPythonClass(src, child, path, lang)
			for j := range classSyms {
				classSyms[j].Annotations = annotations
			}
			symbols = append(symbols, classSyms...)
		}
	}

	return symbols
}

// extractDecorator はデコレータノードからアノテーション文字列を抽出する
func extractDecorator(src []byte, node *gotreesitter.Node) string {
	// @app.route('/') のような完全な表現を取得
	content := nodeContent(src, node)
	return content
}

// isPythonConstant は大文字で構成される定数かどうかを判定する
func isPythonConstant(name string) bool {
	if len(name) == 0 {
		return false
	}
	// 最初の文字が大文字かアンダースコアで始まるか確認
	if name[0] >= 'a' && name[0] <= 'z' {
		return false
	}
	// 大文字、数字、アンダースコアのみで構成されているか確認
	for _, c := range name {
		if c != '_' && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// nodeContent はノードの内容を文字列として返す
func nodeContent(src []byte, node *gotreesitter.Node) string {
	start := node.StartByte()
	end := node.EndByte()
	if int(end) > len(src) {
		end = uint32(len(src))
	}
	return strings.TrimSpace(string(src[start:end]))
}

package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Symbol struct {
	Kind     string // "func", "method", "type", "interface"
	Name     string
	Receiver string // メソッドの場合のみ
	File     string
	Line     int
}

type GoAnalyzer struct{}

func NewGoAnalyzer() *GoAnalyzer {
	return &GoAnalyzer{}
}

func (a *GoAnalyzer) AnalyzeFile(path string) ([]Symbol, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error in %s: %w", path, err)
	}

	var symbols []Symbol

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			sym := Symbol{
				Name: x.Name.Name,
				File: path,
				Line: fset.Position(x.Pos()).Line,
			}
			if x.Recv != nil && len(x.Recv.List) > 0 {
				sym.Kind = "method"
				sym.Receiver = exprToString(x.Recv.List[0].Type)
			} else {
				sym.Kind = "func"
			}
			symbols = append(symbols, sym)

		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					sym := Symbol{
						Name: ts.Name.Name,
						File: path,
						Line: fset.Position(ts.Pos()).Line,
					}
					switch ts.Type.(type) {
					case *ast.InterfaceType:
						sym.Kind = "interface"
					default:
						sym.Kind = "type"
					}
					symbols = append(symbols, sym)
				}
			}
		}
		return true
	})

	return symbols, nil
}

func (a *GoAnalyzer) AnalyzeDir(root string) ([]Symbol, error) {
	var all []Symbol

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// vendor や .git はスキップ
			name := d.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// _test.go は一旦スキップ（必要なら後で入れる）
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		syms, err := a.AnalyzeFile(path)
		if err != nil {
			// パースエラーは無視して続行
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			return nil
		}
		all = append(all, syms...)
		return nil
	})

	return all, err
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	default:
		return fmt.Sprintf("%T", expr)
	}
}
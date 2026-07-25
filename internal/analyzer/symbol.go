package analyzer

type Symbol struct {
	Kind        string   // "func", "method", "type", "interface", "variable", "constant"
	Name        string
	Receiver    string   // メソッドの場合のみ
	File        string
	Line        int
	Annotations []string // @app.route 等のデコレータ
}

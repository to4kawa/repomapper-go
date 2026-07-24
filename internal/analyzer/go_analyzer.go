package analyzer

type Symbol struct {
	Kind     string // "func", "method", "type", "interface"
	Name     string
	Receiver string // メソッドの場合のみ
	File     string
	Line     int
}

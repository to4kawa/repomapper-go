use serde::Serialize;
use std::env;
use std::fs;
use syn::{visit::Visit, Item};
use walkdir::WalkDir;

#[derive(Serialize)]
struct Symbol {
    kind: String,
    name: String,
    file: String,
    line: usize,
}

fn main() {
    let path = env::args().nth(1).expect("usage: repomapper-rust <dir>");
    let mut symbols = Vec::new();

    for entry in WalkDir::new(&path).into_iter().filter_map(|e| e.ok()) {
        let p = entry.path();
        // テストディレクトリ・テストファイルをスキップ
        if p.components().any(|c| c.as_os_str() == "tests") {
            continue;
        }
        if p.file_name()
            .and_then(|n| n.to_str())
            .map(|n| n.ends_with("_test.rs"))
            .unwrap_or(false)
        {
            continue;
        }
        if p.extension().map(|e| e == "rs").unwrap_or(false) {
            if let Ok(src) = fs::read_to_string(p) {
                if let Ok(file) = syn::parse_file(&src) {
                    let mut visitor = SymbolVisitor {
                        file: p.display().to_string(),
                        symbols: &mut symbols,
                        src: &src,
                    };
                    visitor.visit_file(&file);
                }
            }
        }
    }

    println!("{}", serde_json::to_string(&symbols).unwrap());
}

struct SymbolVisitor<'a> {
    file: String,
    symbols: &'a mut Vec<Symbol>,
    src: &'a str,
}

impl<'a> Visit<'a> for SymbolVisitor<'a> {
    fn visit_item(&mut self, item: &'a Item) {
        match item {
            Item::Fn(f) => {
                self.symbols.push(Symbol {
                    kind: "func".into(),
                    name: f.sig.ident.to_string(),
                    file: self.file.clone(),
                    line: line_of(self.src, f.sig.ident.span().start().line),
                });
            }
            Item::Struct(s) => {
                self.symbols.push(Symbol {
                    kind: "type".into(),
                    name: s.ident.to_string(),
                    file: self.file.clone(),
                    line: line_of(self.src, s.ident.span().start().line),
                });
            }
            Item::Enum(e) => {
                self.symbols.push(Symbol {
                    kind: "type".into(),
                    name: e.ident.to_string(),
                    file: self.file.clone(),
                    line: line_of(self.src, e.ident.span().start().line),
                });
            }
            Item::Trait(t) => {
                self.symbols.push(Symbol {
                    kind: "interface".into(),
                    name: t.ident.to_string(),
                    file: self.file.clone(),
                    line: line_of(self.src, t.ident.span().start().line),
                });
            }
            Item::Impl(i) => {
                // impl内のメソッドも取りたい場合はここで掘る
                for item in &i.items {
                    if let syn::ImplItem::Fn(m) = item {
                        self.symbols.push(Symbol {
                            kind: "method".into(),
                            name: m.sig.ident.to_string(),
                            file: self.file.clone(),
                            line: line_of(self.src, m.sig.ident.span().start().line),
                        });
                    }
                }
            }
            _ => {}
        }
        syn::visit::visit_item(self, item);
    }
}

fn line_of(_src: &str, line: usize) -> usize {
    line // synのspanは1-indexed
}

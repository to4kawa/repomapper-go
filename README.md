# repomapper-go

Aider由来のRepoMap機能をGoで再実装したツール。

リポジトリ内の重要なシンボルを抽出し、LLMに渡しやすい簡潔なマップを生成する。

## 現状

- [x] Gitリポジトリの認識
- [x] ファイル一覧の取得
- [x] シンボル抽出（Go / Rust）
- [x] 重要度ランキング
- [x] トークン制限付きマップ生成（`-tokens`）
- [x] MCPサーバー対応（stdio / repo_map ツール）

## 使い方（CLI）

```bash
# 制限なし
go run ./cmd/repomapper <path>

# トークン制限あり
go run ./cmd/repomapper -tokens 8000 <path>
```

## MCP（予定）

stdioでMCPサーバーとして起動し、repo_map ツールを提供する。

```bash
go run ./cmd/mcp
```

## 開発

```bash
go mod tidy
go run ./cmd/repomapper .
go test ./...

# Rust analyzer
cargo build --manifest-path rust-analyzer/Cargo.toml
```

## 設計

詳細は [design.yaml](./design.yaml) を参照。

## ライセンス

MIT
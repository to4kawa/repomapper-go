# repomapper-go

Aider由来のRepoMap機能をGoで再実装したツール。

リポジトリ内の重要なシンボル（関数・型・メソッドなど）を抽出し、LLMに渡しやすい簡潔なマップを生成する。

## 現状

- [x] Gitリポジトリの認識
- [x] ファイル一覧の取得
- [x] シンボル抽出（Go）
- [ ] 重要度ランキング
- [ ] トークン制限付きマップ生成
- [ ] MCPサーバー対応

## 使い方（現在）

```bash
go run ./cmd/repomapper <repository-path>
```

例:

```bash
go run ./cmd/repomapper .
go run ./cmd/repomapper C:\path\to\repo
```

## 開発

```bash
# 依存取得
go mod tidy

# 実行
go run ./cmd/repomapper <path>

# テスト
go test ./...
```

## 設計

詳細は [design.yaml](./design.yaml) を参照。

## ライセンス

MIT
# repomapper-go

Aider由来のRepoMap機能をGoで再実装したツール。

リポジトリ内の重要なシンボルを抽出し、LLMに渡しやすい簡潔なマップを生成する。

## 現状

- [x] Gitリポジトリの認識
- [x] ファイル一覧の取得
- [x] シンボル抽出（Go / Rust）
- [x] 重要度ランキング
- [x] トークン制限付きマップ生成（`-tokens`）
- [x] テスト関数の除外（`-include-tests` で含める）
- [x] MCPサーバー対応（stdio / `repo_map` ツール）

## 使い方（CLI）

```bash
# 制限なし
go run ./cmd/repomapper <path>

# トークン制限あり
go run ./cmd/repomapper -tokens 8000 <path>

# テスト関数を含める（デフォルトは除外）
go run ./cmd/repomapper -include-tests <path>
```

> **注意**: フラグはパスの前に指定すること。

## MCPサーバー

stdioで起動する。

```bash
go run ./cmd/mcp
# または
go build -o repomapper-mcp ./cmd/mcp
# または PATH にインストール
go install ./cmd/mcp
```

### 提供ツール

| ツール | 説明 |
|--------|------|
| `repo_map` | リポジトリのシンボルマップを返す |

引数:
- `path` (string, 必須): リポジトリパス
- `tokens` (int, 任意): 最大トークン数（0または省略で制限なし）
- `include_tests` (bool, 任意): テスト関数を含める（デフォルト: false）

### クライアント設定例

```json
{
  "mcpServers": {
    "repomapper-go": {
      "command": "repomapper-mcp",
      "args": []
    }
  }
}
```

## 開発

```bash
go mod tidy
go run ./cmd/repomapper .
go run ./cmd/mcp
go test ./...

# Rust解析を使う場合
cargo build --release --manifest-path rust-analyzer/Cargo.toml
```

## 設計

詳細は [design.yaml](./design.yaml) を参照。

## ライセンス

MIT

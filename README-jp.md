# alfred-note-md-template

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

> Markdown テンプレート（画像・キャプション込み）を Alfred から note.com のエディタへ貼り付けます。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-note-md-template/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-note-md-template/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-note-md-template/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-note-md-template/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

| 項目 | 内容 |
|---|---|
| 開発対象 | Alfred 5 Script Filter ワークフロー |
| 開発環境 | 個人〜小規模チーム（1〜3人） |
| 実装言語 | Go（サードパーティ依存なし） |
| ライセンス | MIT |
| AI ツール | Claude Code / GitHub Copilot / Gemini CLI |

## What it does

note.com の記事テンプレートをローカルの `.md` ファイルとして管理できます —
画像とイタリック体のキャプションは標準 Markdown 記法で書けます。Alfred で
`note` と入力してテンプレートを選ぶと、テキスト・各画像・キャプションの順に
note.com のエディタへブロックごとに貼り付けます。

```markdown
はじめにのテキスト。

![キャプション](./images/photo.png)

*この画像のキャプションです。*

続きのテキスト。
```

## Setup

1. `.md` テンプレートをフォルダに置きます（デフォルト: `~/Documents/Note Templates`）。
2. デフォルト以外を使う場合は、Alfred Preferences でこのワークフローを開き
   **Configure Workflow** から **Templates Directory** を設定します。
3. ブラウザで note.com のエディタを開きます。
4. Alfred で `note` と入力し、テンプレートを選んで Enter を押します。

インタプリタやサードパーティランタイムのインストールは不要です — ワークフローは
単一のコンパイル済みバイナリとして配布されます。

初回のペースト時に、macOS が Alfred への Accessibility 権限の許可を求める
ことがあります（Cmd+V / Return のキー操作をシミュレートするために必要）。
詳細とトラブルシューティングは [docs/usage.md](docs/usage.md) を参照してください。

## Features

- 単一の静的 Go バイナリ — ランタイムのベンダリング不要、インタプリタの起動コストなし
- Universal（amd64+arm64）ビルド — Intel / Apple Silicon の両方でネイティブ動作
- フルテストスイート — `go test` で Alfred なしにテスト実行可能
- CI/CD — GitHub Actions でリント・テスト・ビルド・リリースを自動化
- AI 対応 — `AI_CONTEXT.md` + `CLAUDE.md` で AI アシスタントのコンテキストを管理

## Requirements

- Alfred 5（Script Filter には Powerpack が必要）
- macOS（クリップボードと `osascript` によるキー操作シミュレーションに必要）

## Quick Start (developers)

```bash
git clone https://github.com/y-marui/alfred-note-md-template
cd alfred-note-md-template

# Script Filter をローカルでシミュレート
templates_dir=~/Documents/Note\ Templates go run ./cmd/note-md-template-alfred ""

# テストを実行
make test

# ワークフローパッケージをビルド
make build-workflow
# → dist/note.com-template-paste-<version>.alfredworkflow
```

`dist/*.alfredworkflow` をダブルクリックして Alfred にインストールします。

## Usage

```
note              テンプレート一覧を表示
note <query>      名前でテンプレートを絞り込む
```

詳しい使い方は [docs/usage.md](docs/usage.md) を参照してください。

## Project Structure

```
alfred-note-md-template/
├── cmd/
│   ├── note-md-template-alfred/       # Script Filter バイナリ
│   └── note-md-template-paste-alfred/ # Run Script アクションバイナリ
├── internal/
│   ├── mdtemplate/     # Markdown テンプレート → ブロックパーサー
│   ├── templatelist/   # テンプレート一覧・絞り込み
│   ├── paste/          # ブロックごとのペースト順序制御
│   ├── clipboard/      # クリップボード書き込み (pbcopy / osascript)
│   ├── keystroke/       # Cmd+V / Return のシミュレーション
│   └── scriptfilter/   # Alfred Script Filter JSON 型
├── workflow/           # Alfred パッケージ (info.plist, icon.png)
├── scripts/            # build-workflow.sh, extract-changelog.sh
└── docs/               # アーキテクチャ・利用ドキュメント
```

## Documentation

| ドキュメント | 内容 |
|---|---|
| [DEVELOPING.md](DEVELOPING.md) | 開発フロー・命名規則・コードレビュー |
| [docs/architecture.md](docs/architecture.md) | アーキテクチャ全体設計 |
| [docs/file-map.md](docs/file-map.md) | ファイルレベルの依存関係マップ |
| [docs/specification.md](docs/specification.md) | 機能仕様・データフロー |
| [docs/ui-design.md](docs/ui-design.md) | Alfred結果アイテムのUI設計指針 |
| [docs/usage.md](docs/usage.md) | エンドユーザー向け利用ガイド |
| [docs/decisions/](docs/decisions/) | アーキテクチャ決定記録 (ADR) |

## AI-Assisted Development

このプロジェクトは AI 支援開発に対応しています。

| ツール | 役割 |
|---|---|
| Claude Code | アーキテクチャ設計・大規模変更・リファクタリング |
| GitHub Copilot | バグ修正・細かな実装・単体テスト作成 |
| Gemini CLI | ドキュメント管理 |

セッションコンテキスト: [`AI_CONTEXT.md`](AI_CONTEXT.md)、[`CLAUDE.md`](CLAUDE.md)

## Release

```bash
# 1. workflow/info.plist のバージョンを更新
# 2. タグを付けてプッシュ
git tag v1.2.3
git push --tags
# GitHub Actions が .alfredworkflow をビルドして GitHub Release を作成
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版 [README.md](README.md) があります。編集時は同一コミットで更新してください。*

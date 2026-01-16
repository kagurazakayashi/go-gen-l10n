![go-gen-l10n](ico/icon.png)

# go-gen-l10n

ARB 駆動 Go ローカライゼーションコードジェネレーター

[English](README.md) | [简体中文](README.zh-Hans.md) | [繁體中文](README.zh-Hant.md) | **日本語**

## 紹介

`go-gen-l10n` は ARB 駆動の Go ローカライゼーション（i18n / l10n）コードジェネレーターです。ディレクトリ内の `app_*.arb` ファイルを自動的にスキャンし、型安全な Go インターフェースと言語実装コードを生成します。Go プロジェクトに国際化を簡単に統合できます。

### ARB とは？

ARB（Application Resource Bundle）は Flutter プロジェクトで普及した JSON ベースのローカライゼーションリソース形式です。各 `.arb` ファイルは 1 つの言語の翻訳リソースを表し、シンプルで理解しやすい構造を持っています。

### 特徴

- ARB ファイルから Go ローカライゼーションコードを自動生成
- 型安全な `AppLocalizations` インターフェース — 翻訳キーの正確性をコンパイル時に検証
- 言語ごとに個別ファイルに分割（`app_localizations.go` + `app_localizations_zh.go` + ...）
- `-lang` フラグでプログラム自体の出力言語を切り替え（生成コードのコメントも翻訳）
- 出力ディレクトリとパッケージ名を設定可能

## クイックスタート

### ステップ 1：ARB 翻訳ファイルを作成

Go プロジェクトのディレクトリにフォルダ（例：`l10n`）を作成し、ARB ファイルを追加します。

ARB ファイルの命名規則：`app_<言語コード>.arb`

例として、日本語と英語の翻訳ファイルを作成します：

#### l10n/app_ja.arb

```json
{
  "@@locale": "ja",
  "hello": "こんにちは",
  "welcome": "ようこそ！"
}
```

#### l10n/app_en.arb

```json
{
  "@@locale": "en",
  "hello": "Hello",
  "welcome": "Welcome!"
}
```

> **注意：** プログラムの言語コードはファイル名から決定されます（`app_ja.arb` → `ja`）。`@@locale` フィールドの値は参照しません。`@@locale` は ARB 形式の慣例的な宣言であり、`@` で始まる他のメタデータキーと同様に無視されます。

### ステップ 2：コードジェネレーターを実行

```bash
go-gen-l10n -dir ./l10n -pkg l10n -lang ja
```

パラメータ：

| パラメータ | デフォルト | 説明                                                  |
| ---------- | ---------- | ----------------------------------------------------- |
| `-dir`     | `./l10n`   | ARB ファイルを含むディレクトリ                        |
| `-pkg`     | `l10n`     | 生成コードの Go パッケージ名                          |
| `-lang`    | `en`       | プログラム出力と生成コードのコメント言語（en/zh/...） |

### ステップ 2.5：生成されるファイル構造

実行後、`-dir` で指定したディレクトリに以下のファイルが生成されます：

```text
l10n/
├── app_en.arb                  # 英語翻訳ソース（手動作成）
├── app_ja.arb                  # 日本語翻訳ソース（手動作成）
├── app_localizations.go        # 基本ファイル：インターフェース定義 + GetLocalizations 関数
├── app_localizations_en.go     # 英語実装（自動生成）
└── app_localizations_ja.go     # 日本語実装（自動生成）
```

> **ヒント：** 自動生成されたファイルを `.gitignore` に追加して、ビルド成果物のコミットを避けましょう：
>
> ```gitignore
> l10n/app_localizations*.go
> *.syso
> go-gen-l10n
> go-gen-l10n.exe
> ```

### ステップ 3：コードで使用

```go
package main

import (
    "fmt"
    "yourproject/l10n"
)

func main() {
    // 日本語の翻訳インスタンスを取得
    l := l10n.GetLocalizations("ja")
    fmt.Println(l.Hello())   // 出力：こんにちは
    fmt.Println(l.Welcome()) // 出力：ようこそ！
}
```

`GetLocalizations()` は指定された言語コードに対応する翻訳インスタンスを返します。サポートされていない言語が指定された場合は、デフォルト言語（最初に読み込まれた ARB ファイルの言語）にフォールバックします。

> **ヒント：** `//go:generate` を使用して、ビルド前に自動的にジェネレーターを実行できます。
>
> プロジェクト内の任意の `.go` ファイル（例：`main.go`）の先頭に以下を追加します：
>
> ```go
> //go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang ja
> ```
>
> 実行可能ファイルを `$PATH` ではなくプロジェクトディレクトリに配置した場合は、相対パスを使用します：
>
> ```go
> // Windows
> //go:generate .\go-gen-l10n.exe -dir ./l10n -pkg l10n -lang ja
>
> // macOS / Linux
> //go:generate ./go-gen-l10n -dir ./l10n -pkg l10n -lang ja
> ```
>
> その後、`go generate ./...` を実行すると自動的にジェネレーターが呼び出されます。

### コマンドリファレンス

```bash
# ヘルプを表示
go-gen-l10n -help

# デフォルト設定で実行（英語出力）
go-gen-l10n

# ディレクトリとパッケージ名を指定
go-gen-l10n -dir ./translations -pkg i18n

# 日本語出力を使用
go-gen-l10n -lang ja

# 中国語（繁体字）出力を使用
go-gen-l10n -lang zh_Hant
```

## デプロイと統合方法

### 方法 1：リリースをダウンロード（初心者におすすめ）

最も簡単な方法で、Go のコンパイル環境は不要です。

1. [Releases](https://github.com/kagurazakayashi/go-gen-l10n/releases) ページを開きます
2. お使いの OS に合ったアーカイブをダウンロードします。
3. 解凍して実行可能ファイル `go-gen-l10n`（Windows では `go-gen-l10n.exe`）を取得します
4. （オプション）実行ファイルのディレクトリをシステム PATH に追加すると、どこからでも呼び出せます

### 方法 2：ソースからビルド

Go がインストールされている場合（**バージョン 1.24.4 以上が必要**）、ソースから直接ビルドできます：

#### macOS / Linux

```bash
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go build -o go-gen-l10n .
./go-gen-l10n -help
```

#### Windows（コマンドプロンプトまたは PowerShell）

```batch
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go generate
go build -o go-gen-l10n.exe .
go-gen-l10n.exe -help
```

> **説明：**
>
> - `go mod tidy` は依存関係を整理し、`go.mod` を最新の状態に保ちます
> - `go generate` は、Windows アプリケーションの作成に必要なアイコンなどのリソースを生成します。
> - `go build -o <ファイル名>` はコンパイルし、指定されたパスに実行ファイルを出力します

#### コンパイルしたバイナリを `$GOPATH/bin` にインストール

```bash
# macOS / Linux / Windows 共通
go install
```

これで任意の場所から `go-gen-l10n` を実行できます。

### 方法 3：Git サブモジュールとしてプロジェクトに統合

この方法では、ジェネレーターを Go プロジェクトに埋め込み、チームメンバーが個別にインストールする必要がなくなります。

#### 1. Git サブモジュールとして追加

```bash
# プロジェクトのルートディレクトリで実行
git submodule add https://github.com/kagurazakayashi/go-gen-l10n.git tools/go-gen-l10n
```

#### 2. go.mod でローカルパスを参照

```bash
go mod edit -replace github.com/kagurazakayashi/go-gen-l10n=./tools/go-gen-l10n
```

#### 3. ローカルにインストール

```bash
go mod tidy
go install ./tools/go-gen-l10n
```

#### 4. `//go:generate` で自動実行

エントリファイル（例：`main.go`）の先頭に以下を追加します：

```go
//go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang ja
```

その後、`go generate ./...` を実行すると自動的にジェネレーターが呼び出されます。

## ユニットテスト

```bash
go test ./...
```

| テスト範囲             | テスト名                                                   | デモ内容                                                     |
| ---------------------- | ---------------------------------------------------------- | ------------------------------------------------------------ |
| ARB 解析               | `TestArbMap`                                               | `@` プレフィックスキーのフィルタリング、非文字列値のスキップ |
| 文字列変換             | `TestToCamelCase`                                          | `snake_case` と `kebab-case` から `PascalCase` への変換      |
| 特殊文字               | `TestArbMapTranslationQuoting`                             | 引用符、改行、バックスラッシュの Go 文字列エスケープ         |
| テンプレート           | `TestBaseTemplateRendering`、`TestLocaleTemplateRendering` | テンプレート出力構造の検証                                   |
| 欠落翻訳               | `TestLocaleTemplateMissingTranslation`                     | 翻訳がない場合のキー名へのフォールバック                     |
| エンドツーエンド       | `TestFullGeneration`                                       | 一時 ARB ファイル作成、コード生成、出力検証                  |
| 生成パッケージ         | `TestGetLocalizationsKnownLocales`                         | `GetLocalizations()` がロケールごとに正しい型を返す          |
| ロケールフォールバック | `TestGetLocalizationsFallback`                             | 不明なロケールがデフォルト言語にフォールバック               |
| 翻訳値                 | `TestAppLocalizationsEnMethods` など                       | 各ロケールの実際の翻訳テキスト                               |

## ライセンス

```LICENSE
Copyright (c) 2026 KagurazakaYashi
go-gen-l10n is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
```

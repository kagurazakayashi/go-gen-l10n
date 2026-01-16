![go-gen-l10n](ico/icon.png)

# go-gen-l10n

ARB-driven Go localization code generator

**English** | [简体中文](README.zh-Hans.md) | [繁體中文](README.zh-Hant.md) | [日本語](README_ja.md)

## Introduction

`go-gen-l10n` is an ARB-driven Go localization (i18n / l10n) code generator. It automatically scans `app_*.arb` files in a directory and generates type-safe Go interfaces and locale implementations, making it easy to integrate internationalization into your Go project.

### What is ARB?

ARB (Application Resource Bundle) is a JSON-based localization resource format popularized by the Flutter project. Each `.arb` file represents the translation resources for a single language, with a simple and easy-to-understand structure.

### Features

- Automatically generates Go localization code from ARB files
- Type-safe `AppLocalizations` interface — translation key correctness verified at compile time
- Per-locale file splitting (`app_localizations.go` + `app_localizations_zh.go` + ...)
- `-lang` flag to switch program output language (generated code comments are also translated)
- Configurable output directory and package name

## Quick Start

### Step 1: Create ARB Translation Files

In your Go project directory, create a folder (e.g., `l10n`) and add ARB files.

ARB file naming convention: `app_<locale_code>.arb`

For example, create English and Chinese translation files:

#### l10n/app_en.arb

```json
{
  "@@locale": "en",
  "hello": "Hello",
  "welcome": "Welcome!"
}
```

#### l10n/app_zh.arb

```json
{
  "@@locale": "zh",
  "hello": "你好",
  "welcome": "欢迎！"
}
```

> **Note:** The locale code is determined by the filename (`app_zh.arb` → `zh`), not the `@@locale` field. `@@locale` is an ARB convention — this tool treats it like any other `@`-prefixed meta key and ignores it.

### Step 2: Run the Code Generator

```bash
go-gen-l10n -dir ./l10n -pkg l10n -lang en
```

Flags:

| Flag    | Default  | Description                                                         |
| ------- | -------- | ------------------------------------------------------------------- |
| `-dir`  | `./l10n` | Directory containing ARB files                                      |
| `-pkg`  | `l10n`   | Output Go package name                                              |
| `-lang` | `en`     | Language for program output and generated code comments (en/zh/...) |

### Step 2.5: Generated File Structure

After execution, the following files are generated in the directory specified by `-dir`:

```text
l10n/
├── app_en.arb                  # English translation source (created manually)
├── app_zh.arb                  # Chinese translation source (created manually)
├── app_localizations.go        # Base file: interface definition + GetLocalizations function
├── app_localizations_en.go     # English implementation (auto-generated)
└── app_localizations_zh.go     # Chinese implementation (auto-generated)
```

> **Tip:** Add generated files to `.gitignore` to avoid committing build artifacts:
>
> ```gitignore
> l10n/app_localizations*.go
> *.syso
> go-gen-l10n
> go-gen-l10n.exe
> ```

### Step 3: Use in Your Code

```go
package main

import (
    "fmt"
    "yourproject/l10n"
)

func main() {
    // Get the English localization instance
    l := l10n.GetLocalizations("en")
    fmt.Println(l.Hello())   // Output: Hello
    fmt.Println(l.Welcome()) // Output: Welcome!
}
```

`GetLocalizations()` returns the translation instance for the given locale code. If the requested locale is not supported, it falls back to the default locale (the first loaded ARB file's language).

> **Tip:** You can use `//go:generate` to automatically run the generator before building.
>
> Add the following comment at the top of any `.go` file in your project (e.g., `main.go`):
>
> ```go
> //go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang en
> ```
>
> If you placed the binary in your project directory instead of `$PATH`, use a relative path:
>
> ```go
> // Windows
> //go:generate .\go-gen-l10n.exe -dir ./l10n -pkg l10n -lang en
>
> // macOS / Linux
> //go:generate ./go-gen-l10n -dir ./l10n -pkg l10n -lang en
> ```
>
> Then run `go generate ./...` to invoke the generator automatically.

### Command Reference

```bash
# Show help
go-gen-l10n -help

# Run with defaults (English output)
go-gen-l10n

# Specify directory and package name
go-gen-l10n -dir ./translations -pkg i18n

# Use Japanese output
go-gen-l10n -lang ja

# Use Traditional Chinese output
go-gen-l10n -lang zh_Hant
```

## Deployment and Integration

### Option 1: Download Release (Recommended for Beginners)

This is the easiest approach — no Go compilation environment needed.

1. Visit the [Releases](https://github.com/kagurazakayashi/go-gen-l10n/releases) page
2. Download the archive matching your OS
3. Extract to get the executable `go-gen-l10n` (`go-gen-l10n.exe` on Windows)
4. (Optional) Add the executable directory to your system PATH for global access

### Option 2: Build from Source

If you have Go installed (**1.24.4 or later required**), build directly from source:

#### macOS / Linux

```bash
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go build -o go-gen-l10n .
./go-gen-l10n -help
```

#### Windows (Command Prompt or PowerShell)

```batch
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go generate
go build -o go-gen-l10n.exe .
go-gen-l10n.exe -help
```

> **Explanation:**
>
> - `go mod tidy` syncs dependencies and keeps `go.mod` up to date
> - `go generate` runs the resources required to create a Windows application, such as icons. Running this before compilation is recommended
> - `go build -o <filename>` compiles and outputs the executable to the specified path

#### Install the compiled binary to `$GOPATH/bin`

```bash
# Works on macOS / Linux / Windows
go install
```

You can then run `go-gen-l10n` from anywhere.

### Option 3: Integrate as a Git Submodule

This approach embeds the generator in your Go project so teammates don't need to install it separately.

#### 1. Add as a Git submodule

```bash
# Run in your project root
git submodule add https://github.com/kagurazakayashi/go-gen-l10n.git tools/go-gen-l10n
```

#### 2. Reference the local path in go.mod

```bash
go mod edit -replace github.com/kagurazakayashi/go-gen-l10n=./tools/go-gen-l10n
```

#### 3. Install locally

```bash
go mod tidy
go install ./tools/go-gen-l10n
```

#### 4. Use `//go:generate` for automatic execution

Add the following comment at the top of your entry file (e.g., `main.go`):

```go
//go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang en
```

Then run `go generate ./...` to invoke the generator automatically.

## Unit Tests

```bash
go test ./...
```

| Area                 | Tests                                                      | What it demonstrates                                          |
| -------------------- | ---------------------------------------------------------- | ------------------------------------------------------------- |
| ARB parsing          | `TestArbMap`                                               | How `@`-prefixed keys are filtered, non-string values skipped |
| String conversion    | `TestToCamelCase`                                          | How `snake_case` and `kebab-case` keys become `PascalCase`    |
| Special characters   | `TestArbMapTranslationQuoting`                             | Proper Go string escaping for quotes, newlines, backslashes   |
| Templates            | `TestBaseTemplateRendering`, `TestLocaleTemplateRendering` | Template output structure verification                        |
| Missing translations | `TestLocaleTemplateMissingTranslation`                     | Fallback to raw key name when translation is absent           |
| End-to-end           | `TestFullGeneration`                                       | Creates temp ARB files, generates code, validates output      |
| Generated package    | `TestGetLocalizationsKnownLocales`                         | `GetLocalizations()` returns correct type per locale          |
| Locale fallback      | `TestGetLocalizationsFallback`                             | Unknown locales fall back to default language                 |
| Translation values   | `TestAppLocalizationsEnMethods` etc.                       | Actual translation text for each locale                       |

## License

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

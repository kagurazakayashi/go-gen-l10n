![go-gen-l10n](./ico/icon.ico)

# go-gen-l10n

ARB 驅動的 Go 多語言程式碼生成工具

[English](README.md) | [簡體中文](README.zh-Hans.md) | **繁體中文** | [日本語](README_ja.md)

## 簡介

`go-gen-l10n` 是一個 ARB 驅動的 Go 多語言（i18n / l10n）程式碼生成工具。它會自動掃描目錄中的 `app_*.arb` 檔案，生成型別安全的 Go 介面與語言實現程式碼，讓你在 Go 專案中輕鬆整合國際化支援。

### 什麼是 ARB？

ARB（Application Resource Bundle）是一種基於 JSON 的本地化資源格式，由 Flutter 專案推廣使用。每個 `.arb` 檔案代表一個語言的翻譯資源，結構簡潔易懂。

### 特性

- 從 ARB 檔案自動生成 Go 本地化程式碼
- 型別安全的 `AppLocalizations` 介面，編譯時檢查翻譯鍵是否正確
- 按語言拆分為獨立檔案（`app_localizations.go` + `app_localizations_zh.go` + ...）
- 支援 `-lang` 引數切換程式自身的輸出語言（生成程式碼的註釋也會一起翻譯）
- 命令列引數可配置輸出目錄與包名

## 快速開始

### 步驟 1：建立 ARB 翻譯檔案

在你的 Go 專案目錄中，新建一個資料夾（如 `l10n`），並在其中建立 ARB 檔案。

ARB 檔案的命名規則：`app_<語言程式碼>.arb`

例如，建立英文和中文的翻譯檔案：

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
  "welcome": "歡迎！"
}
```

> **注意：** 程式的語系識別碼由檔案名稱決定（`app_zh.arb` → `zh`），而非 `@@locale` 欄位的值。`@@locale` 是 ARB 格式的慣例宣告，本工具會將其與其它 `@` 開頭的後設資料鍵一併忽略。

### 步驟 2：執行程式碼生成器

```bash
go-gen-l10n -dir ./l10n -pkg l10n -lang zh
```

引數說明：

| 引數    | 預設值   | 說明                                            |
| ------- | -------- | ----------------------------------------------- |
| `-dir`  | `./l10n` | ARB 檔案所在目錄                                |
| `-pkg`  | `l10n`   | 生成程式碼的 Go 包名                            |
| `-lang` | `en`     | 程式輸出資訊與生成程式碼註釋的語言（en/zh/...） |

### 步驟 2.5：生成後的檔案結構

執行上面命令後，`-dir` 目錄中會生成如下檔案：

```text
l10n/
├── app_en.arb                  # 英文翻譯原始檔（你手動建立的）
├── app_zh.arb                  # 中文翻譯原始檔（你手動建立的）
├── app_localizations.go        # 基礎檔案：介面定義 + GetLocalizations 函式
├── app_localizations_en.go     # 英文實現（程式自動生成）
└── app_localizations_zh.go     # 中文實現（程式自動生成）
```

> **建議：** 將自動生成的檔案新增到 `.gitignore` 中，避擴音交編譯產物：
>
> ```gitignore
> l10n/app_localizations*.go
> *.syso
> go-gen-l10n
> go-gen-l10n.exe
> ```

### 步驟 3：在程式碼中使用

```go
package main

import (
    "fmt"
    "yourproject/l10n"
)

func main() {
    // 獲取中文翻譯例項
    l := l10n.GetLocalizations("zh")
    fmt.Println(l.Hello())   // 輸出：你好
    fmt.Println(l.Welcome()) // 輸出：歡迎！
}
```

`GetLocalizations()` 會根據傳入的語言程式碼返回對應的翻譯例項。如果傳入未支援的語言程式碼，會返回預設語言（即第一個載入的 ARB 檔案對應的語言）的翻譯。

> **提示：** 你可以透過 `//go:generate` 指令在專案構建前自動執行程式碼生成器。

### 常用命令速查

```bash
# 檢視幫助
go-gen-l10n -help

# 使用預設引數生成（英文輸出）
go-gen-l10n

# 指定目錄和包名
go-gen-l10n -dir ./translations -pkg i18n

# 使用日文輸出
go-gen-l10n -lang ja

# 使用繁體中文輸出
go-gen-l10n -lang zh_Hant
```

## 部署和整合方式

### 方式一：下載 Release（推薦新手使用）

這是最簡單的使用方式，無需安裝 Go 編譯環境。

1. 開啟 [Releases](https://github.com/kagurazakayashi/go-gen-l10n/releases) 頁面
2. 根據你的作業系統下載對應的壓縮包。
3. 解壓後得到可執行檔案 `go-gen-l10n`（Windows 平臺為 `go-gen-l10n.exe`）
4. （可選）將可執行檔案所在目錄新增到系統 PATH，方便全域性呼叫

### 方式二：手動編譯

如果你已經安裝了 Go（**需要 1.24.4 或更高版本**），可以直接從原始碼編譯：

#### macOS / Linux

```bash
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go build -o go-gen-l10n .
./go-gen-l10n -help
```

#### Windows（命令提示符或 PowerShell）

```batch
git clone https://github.com/kagurazakayashi/go-gen-l10n.git
cd go-gen-l10n
go mod tidy
go generate
go build -o go-gen-l10n.exe .
go-gen-l10n.exe -help
```

> **說明：**
>
> - `go mod tidy` 整理依賴關係，確保 `go.mod` 檔案是最新的
> - `go generate` 會執行建立 Windows 程式所需的圖示等資源
> - `go build -o <檔名>` 編譯並輸出可執行檔案到你指定的路徑

#### 將編譯好的程式安裝到 `$GOPATH/bin`

```bash
# macOS / Linux / Windows 通用
go install
```

之後可以在任意位置直接執行 `go-gen-l10n`。

### 方式三：作為 Git 子模組引入專案

這種方式適合將程式碼生成器整合到你的 Go 專案中，讓團隊其他成員無需單獨安裝。

#### 1. 新增為 Git 子模組

```bash
# 在你的專案根目錄下執行
git submodule add https://github.com/kagurazakayashi/go-gen-l10n.git tools/go-gen-l10n
```

#### 2. 在 go.mod 中引用本地路徑

```bash
go mod edit -replace github.com/kagurazakayashi/go-gen-l10n=./tools/go-gen-l10n
```

#### 3. 安裝到本地

```bash
go mod tidy
go install ./tools/go-gen-l10n
```

#### 4. 在程式碼中使用 `//go:generate` 自動執行

在你的專案入口檔案（如 `main.go`）頂部新增：

```go
//go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang zh
```

之後執行 `go generate ./...` 即可自動呼叫程式碼生成器。

## 單元測試

```bash
go test ./...
```

| 測試範圍   | 測試名稱                                                   | 演示內容                                     |
| ---------- | ---------------------------------------------------------- | -------------------------------------------- |
| ARB 解析   | `TestArbMap`                                               | `@` 前綴鍵過濾、非字串值跳過                 |
| 字串轉換   | `TestToCamelCase`                                          | `snake_case` 和 `kebab-case` 轉 `PascalCase` |
| 特殊字元   | `TestArbMapTranslationQuoting`                             | 引號、換行、反斜線的 Go 字串轉義             |
| 範本渲染   | `TestBaseTemplateRendering`、`TestLocaleTemplateRendering` | 範本輸出結構驗證                             |
| 缺失翻譯   | `TestLocaleTemplateMissingTranslation`                     | 翻譯缺失時回退到原始 key 名稱                |
| 端到端     | `TestFullGeneration`                                       | 建立暫存 ARB 檔案、產生程式碼、驗證輸出      |
| 產生的套件 | `TestGetLocalizationsKnownLocales`                         | `GetLocalizations()` 按語言回傳正確型別      |
| 語言回退   | `TestGetLocalizationsFallback`                             | 未知語言程式碼回退到預設語言                 |
| 翻譯值     | `TestAppLocalizationsEnMethods` 等                         | 每種語言的實際翻譯文字                       |

## 許可證

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

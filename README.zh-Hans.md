![go-gen-l10n](ico/icon.png)

# go-gen-l10n

ARB 驱动的 Go 多语言代码生成工具

[English](README.md) | **简体中文** | [繁體中文](README.zh-Hant.md) | [日本語](README_ja.md)

## 简介

`go-gen-l10n` 是一个 ARB 驱动的 Go 多语言（i18n / l10n）代码生成工具。它会自动扫描目录中的 `app_*.arb` 文件，生成类型安全的 Go 接口与语言实现代码，让你在 Go 项目中轻松集成国际化支持。

### 什么是 ARB？

ARB（Application Resource Bundle）是一种基于 JSON 的本地化资源格式，由 Flutter 项目推广使用。每个 `.arb` 文件代表一个语言的翻译资源，结构简洁易懂。

### 特性

- 从 ARB 文件自动生成 Go 本地化代码
- 类型安全的 `AppLocalizations` 接口，编译时检查翻译键是否正确
- 按语言拆分为独立文件（`app_localizations.go` + `app_localizations_zh.go` + ...）
- 支持 `-lang` 参数切换程序自身的输出语言（生成代码的注释也会一起翻译）
- 命令行参数可配置输出目录与包名

## 快速开始

### 步骤 1：创建 ARB 翻译文件

在你的 Go 项目目录中，新建一个文件夹（如 `l10n`），并在其中创建 ARB 文件。

ARB 文件的命名规则：`app_<语言代码>.arb`

例如，创建英文和中文的翻译文件：

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

> **注意：** 程序的语系识别码由文件名决定（`app_zh.arb` → `zh`），而非 `@@locale` 字段的值。`@@locale` 是 ARB 格式的惯例声明，本工具会将其与其它 `@` 开头的元数据键一并忽略。

### 步骤 2：运行代码生成器

```bash
go-gen-l10n -dir ./l10n -pkg l10n -lang zh
```

参数说明：

| 参数    | 默认值   | 说明                                          |
| ------- | -------- | --------------------------------------------- |
| `-dir`  | `./l10n` | ARB 文件所在目录                              |
| `-pkg`  | `l10n`   | 生成代码的 Go 包名                            |
| `-lang` | `en`     | 程序输出信息与生成代码注释的语言（en/zh/...） |

### 步骤 2.5：生成后的文件结构

运行上面命令后，`-dir` 目录中会生成如下文件：

```text
l10n/
├── app_en.arb                  # 英文翻译源文件（你手动创建的）
├── app_zh.arb                  # 中文翻译源文件（你手动创建的）
├── app_localizations.go        # 基础文件：接口定义 + GetLocalizations 函数
├── app_localizations_en.go     # 英文实现（程序自动生成）
└── app_localizations_zh.go     # 中文实现（程序自动生成）
```

> **建议：** 将自动生成的文件添加到 `.gitignore` 中，避免提交编译产物：
>
> ```gitignore
> l10n/app_localizations*.go
> *.syso
> go-gen-l10n
> go-gen-l10n.exe
> ```

### 步骤 3：在代码中使用

```go
package main

import (
    "fmt"
    "yourproject/l10n"
)

func main() {
    // 获取中文翻译实例
    l := l10n.GetLocalizations("zh")
    fmt.Println(l.Hello())   // 输出：你好
    fmt.Println(l.Welcome()) // 输出：欢迎！
}
```

`GetLocalizations()` 会根据传入的语言代码返回对应的翻译实例。如果传入未支持的语言代码，会返回默认语言（即第一个加载的 ARB 文件对应的语言）的翻译。

> **提示：** 你可以通过 `//go:generate` 指令在项目构建前自动运行代码生成器。
>
> 在项目的任意 `.go` 文件（如 `main.go`）顶部添加：
>
> ```go
> //go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang zh
> ```
>
> 如果你将可执行文件放在了项目目录而非 `$PATH` 中，请使用相对路径：
>
> ```go
> // Windows
> //go:generate .\go-gen-l10n.exe -dir ./l10n -pkg l10n -lang zh
>
> // macOS / Linux
> //go:generate ./go-gen-l10n -dir ./l10n -pkg l10n -lang zh
> ```
>
> 之后运行 `go generate ./...` 即可自动调用代码生成器。

### 常用命令速查

```bash
# 查看帮助
go-gen-l10n -help

# 使用默认参数生成（英文输出）
go-gen-l10n

# 指定目录和包名
go-gen-l10n -dir ./translations -pkg i18n

# 使用日文输出
go-gen-l10n -lang ja

# 使用繁体中文输出
go-gen-l10n -lang zh_Hant
```

## 部署和集成方式

### 方式一：下载 Release（推荐新手使用）

这是最简单的使用方式，无需安装 Go 编译环境。

1. 打开 [Releases](https://github.com/kagurazakayashi/go-gen-l10n/releases) 页面
2. 根据你的操作系统下载对应的压缩包。
3. 解压后得到可执行文件 `go-gen-l10n`（Windows 平台为 `go-gen-l10n.exe`）
4. （可选）将可执行文件所在目录添加到系统 PATH，方便全局调用

### 方式二：手动编译

如果你已经安装了 Go（**需要 1.24.4 或更高版本**），可以直接从源码编译：

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

> **说明：**
>
> - `go mod tidy` 整理依赖关系，确保 `go.mod` 文件是最新的
> - `go generate` 会运行创建 Windows 程序所需的图标等资源
> - `go build -o <文件名>` 编译并输出可执行文件到你指定的路径

#### 将编译好的程序安装到 `$GOPATH/bin`

```bash
# macOS / Linux / Windows 通用
go install
```

之后可以在任意位置直接运行 `go-gen-l10n`。

### 方式三：作为 Git 子模块引入项目

这种方式适合将代码生成器集成到你的 Go 项目中，让团队其他成员无需单独安装。

#### 1. 添加为 Git 子模块

```bash
# 在你的项目根目录下执行
git submodule add https://github.com/kagurazakayashi/go-gen-l10n.git tools/go-gen-l10n
```

#### 2. 在 go.mod 中引用本地路径

```bash
go mod edit -replace github.com/kagurazakayashi/go-gen-l10n=./tools/go-gen-l10n
```

#### 3. 安装到本地

```bash
go mod tidy
go install ./tools/go-gen-l10n
```

#### 4. 在代码中使用 `//go:generate` 自动运行

在你的项目入口文件（如 `main.go`）顶部添加：

```go
//go:generate go-gen-l10n -dir ./l10n -pkg l10n -lang zh
```

之后运行 `go generate ./...` 即可自动调用代码生成器。

## 单元测试

```bash
go test ./...
```

| 测试范围   | 测试名称                                                   | 演示内容                                     |
| ---------- | ---------------------------------------------------------- | -------------------------------------------- |
| ARB 解析   | `TestArbMap`                                               | `@` 前缀键过滤、非字符串值跳过               |
| 字符串转换 | `TestToCamelCase`                                          | `snake_case` 和 `kebab-case` 转 `PascalCase` |
| 特殊字符   | `TestArbMapTranslationQuoting`                             | 引号、换行、反斜杠的 Go 字符串转义           |
| 模板渲染   | `TestBaseTemplateRendering`、`TestLocaleTemplateRendering` | 模板输出结构验证                             |
| 缺失翻译   | `TestLocaleTemplateMissingTranslation`                     | 翻译缺失时回退到原始 key 名称                |
| 端到端     | `TestFullGeneration`                                       | 创建临时 ARB 文件、生成代码、验证输出        |
| 生成的包   | `TestGetLocalizationsKnownLocales`                         | `GetLocalizations()` 按语言返回正确类型      |
| 语言回退   | `TestGetLocalizationsFallback`                             | 未知语言代码回退到默认语言                   |
| 翻译值     | `TestAppLocalizationsEnMethods` 等                         | 每种语言的实际翻译文本                       |

## 许可证

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

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func generateGoCode(dir, pkgName string, data TemplateData) {
	tmpl, err := template.New("l10n").Parse(CodeTemplate)
	if err != nil {
		log.Fatalf("解析模板失败: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("执行模板失败: %v", err)
	}
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("格式化生成的代码失败: %v\n代码内容:\n%s", err, buf.String())
	}
	outPath := filepath.Join(dir, "app_localizations.go")
	if err := os.WriteFile(outPath, formattedCode, fs.ModePerm); err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}

	fmt.Printf("成功生成本地化代码: %s\n", outPath)
}

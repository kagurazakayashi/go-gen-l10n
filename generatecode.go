package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func generateGoCode(dir, pkgName string, data TemplateData) {
	// 生成基础文件（接口 + GetLocalizations）
	baseTmpl, err := template.New("base").Parse(BaseTemplate)
	if err != nil {
		log.Fatalf("解析基础模板失败: %v", err)
	}

	var baseBuf bytes.Buffer
	if err := baseTmpl.Execute(&baseBuf, data); err != nil {
		log.Fatalf("执行基础模板失败: %v", err)
	}

	formattedBase, err := format.Source(baseBuf.Bytes())
	if err != nil {
		log.Fatalf("格式化基础代码失败: %v\n代码内容:\n%s", err, baseBuf.String())
	}

	basePath := filepath.Join(dir, "app_localizations.go")
	if err := os.WriteFile(basePath, formattedBase, fs.ModePerm); err != nil {
		log.Fatalf("写入基础文件失败: %v", err)
	}
	fmt.Printf("成功生成本地化代码: %s\n", basePath)

	// 为每个语言生成单独的实现文件
	localeTmpl, err := template.New("locale").Parse(LocaleTemplate)
	if err != nil {
		log.Fatalf("解析语言模板失败: %v", err)
	}

	for _, locale := range data.Locales {
		var localeBuf bytes.Buffer
		localeData := map[string]interface{}{
			"PackageName": pkgName,
			"Keys":        data.Keys,
			"Locale":      locale,
		}
		if err := localeTmpl.Execute(&localeBuf, localeData); err != nil {
			log.Fatalf("执行语言模板失败 (%s): %v", locale.ID, err)
		}

		formattedLocale, err := format.Source(localeBuf.Bytes())
		if err != nil {
			log.Fatalf("格式化语言代码失败 (%s): %v\n代码内容:\n%s", locale.ID, err, localeBuf.String())
		}

		// 文件名使用小写语言代码（如 app_localizations_en.go）
		localeFileName := fmt.Sprintf("app_localizations_%s.go", strings.ToLower(locale.ID))
		localePath := filepath.Join(dir, localeFileName)
		if err := os.WriteFile(localePath, formattedLocale, fs.ModePerm); err != nil {
			log.Fatalf("写入语言文件失败 (%s): %v", locale.ID, err)
		}
		fmt.Printf("成功生成本地化代码: %s\n", localePath)
	}
}

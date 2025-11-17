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

	"github.com/kagurazakayashi/go-gen-l10n/l10n"
)

func generateGoCode(dir, pkgName string, data TemplateData, L l10n.AppLocalizations) {
	// 生成基礎文件（接口 + GetLocalizations）
	baseTemplateData := map[string]interface{}{
		"PackageName":                    pkgName,
		"Keys":                           data.Keys,
		"Locales":                        data.Locales,
		"DefaultStructSuffix":            data.DefaultStructSuffix,
		"CommentAppLocalizationsInterface": L.CommentAppLocalizationsInterface(),
		"CommentGetLocalizations":        L.CommentGetLocalizations(),
	}

	baseTmpl, err := template.New("base").Parse(BaseTemplate)
	if err != nil {
		log.Fatalf(L.ErrorParseBaseTemplate(), err)
	}

	var baseBuf bytes.Buffer
	if err := baseTmpl.Execute(&baseBuf, baseTemplateData); err != nil {
		log.Fatalf(L.ErrorExecuteBaseTemplate(), err)
	}

	formattedBase, err := format.Source(baseBuf.Bytes())
	if err != nil {
		log.Fatalf(L.ErrorFormatBaseCode(), err, baseBuf.String())
	}

	basePath := filepath.Join(dir, "app_localizations.go")
	if err := os.WriteFile(basePath, formattedBase, fs.ModePerm); err != nil {
		log.Fatalf(L.ErrorWriteBaseFile(), err)
	}
	fmt.Printf(L.SuccessGeneratedCode(), basePath)

	// 為每個語言生成單獨的實現文件
	localeTmpl, err := template.New("locale").Parse(LocaleTemplate)
	if err != nil {
		log.Fatalf(L.ErrorParseLocaleTemplate(), err)
	}

	for _, locale := range data.Locales {
		var localeBuf bytes.Buffer
		localeData := map[string]interface{}{
			"PackageName":              pkgName,
			"Keys":                     data.Keys,
			"Locale":                   locale,
			"GeneratedLocale":          data.GeneratedLocale,
			"CommentLocaleImplementation": fmt.Sprintf(L.CommentLocaleImplementation(), locale.ID),
		}
		if err := localeTmpl.Execute(&localeBuf, localeData); err != nil {
			log.Fatalf(L.ErrorExecuteLocaleTemplate(), locale.ID, err)
		}

		formattedLocale, err := format.Source(localeBuf.Bytes())
		if err != nil {
			log.Fatalf(L.ErrorFormatLocaleCode(), locale.ID, err, localeBuf.String())
		}

		// 文件名使用小寫語言代碼（如 app_localizations_en.go）
		localeFileName := fmt.Sprintf("app_localizations_%s.go", strings.ToLower(locale.ID))
		localePath := filepath.Join(dir, localeFileName)
		if err := os.WriteFile(localePath, formattedLocale, fs.ModePerm); err != nil {
			log.Fatalf(L.ErrorWriteLocaleFile(), locale.ID, err)
		}
		fmt.Printf(L.SuccessGeneratedCode(), localePath)
	}
}

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

// generateGoCode 依據提供的範本資料產生多語系 Go 程式碼。
// 會產生一個基礎檔案（介面 + GetLocalizations 函式），
// 以及每個語言一個實作檔案。
func generateGoCode(dir, pkgName string, data TemplateData, L l10n.AppLocalizations) {
	// 組合基礎範本的資料——用於產生 app_localizations.go
	baseTemplateData := map[string]interface{}{
		"PackageName":                    pkgName,
		"Keys":                           data.Keys,
		"Locales":                        data.Locales,
		"DefaultStructSuffix":            data.DefaultStructSuffix,
		"CommentAppLocalizationsInterface": L.CommentAppLocalizationsInterface(),
		"CommentGetLocalizations":        L.CommentGetLocalizations(),
	}

	// 解析並執行基礎範本，產生 AppLocalizations 介面與 GetLocalizations 函式
	baseTmpl, err := template.New("base").Parse(BaseTemplate)
	if err != nil {
		log.Fatalf(L.ErrorParseBaseTemplate(), err)
	}

	var baseBuf bytes.Buffer
	if err := baseTmpl.Execute(&baseBuf, baseTemplateData); err != nil {
		log.Fatalf(L.ErrorExecuteBaseTemplate(), err)
	}

	// 使用 go/format 格式化產生的程式碼
	formattedBase, err := format.Source(baseBuf.Bytes())
	if err != nil {
		log.Fatalf(L.ErrorFormatBaseCode(), err, baseBuf.String())
	}

	// 寫入基礎檔案
	basePath := filepath.Join(dir, "app_localizations.go")
	if err := os.WriteFile(basePath, formattedBase, fs.ModePerm); err != nil {
		log.Fatalf(L.ErrorWriteBaseFile(), err)
	}
	fmt.Printf(L.SuccessGeneratedCode(), basePath)

	// 解析語言範本——用於產生每個語言的實作檔案
	localeTmpl, err := template.New("locale").Parse(LocaleTemplate)
	if err != nil {
		log.Fatalf(L.ErrorParseLocaleTemplate(), err)
	}

	// 為每個語言產生單獨的實作檔案
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

		// 檔案名稱使用小寫語言代碼（如 app_localizations_en.go）
		localeFileName := fmt.Sprintf("app_localizations_%s.go", strings.ToLower(locale.ID))
		localePath := filepath.Join(dir, localeFileName)
		if err := os.WriteFile(localePath, formattedLocale, fs.ModePerm); err != nil {
			log.Fatalf(L.ErrorWriteLocaleFile(), locale.ID, err)
		}
		fmt.Printf(L.SuccessGeneratedCode(), localePath)
	}
}

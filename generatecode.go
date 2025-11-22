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

	// 寫入基礎檔案並輸出差異資訊
	basePath := filepath.Join(dir, "app_localizations.go")
	writeFileWithDiff(basePath, formattedBase, L)

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
		writeFileWithDiff(localePath, formattedLocale, L)
	}
}

// writeFileWithDiff 寫入檔案並輸出與既有檔案的差異統計。
// 若檔案為新建立則顯示 (new)；若為更新則顯示 (+N -M) 行數變化。
func writeFileWithDiff(path string, content []byte, L l10n.AppLocalizations) {
	oldContent, _ := os.ReadFile(path)

	if err := os.WriteFile(path, content, fs.ModePerm); err != nil {
		log.Fatalf(L.ErrorWriteLocaleFile(), path, err)
	}

	msg := fmt.Sprintf(L.SuccessGeneratedCode(), path)
	if oldContent == nil {
		fmt.Printf("%s (new)\n", msg)
	} else {
		added, removed := lineDiff(oldContent, content)
		if added == 0 && removed == 0 {
			fmt.Printf("%s (unchanged)\n", msg)
		} else {
			fmt.Printf("%s (+%d -%d)\n", msg, added, removed)
		}
	}
}

// lineDiff 計算兩份內容的行級差異。
// 回傳新增行數與刪除行數。
func lineDiff(old, new []byte) (added, removed int) {
	oldLines := strings.Split(string(old), "\n")
	newLines := strings.Split(string(new), "\n")

	// 使用簡單的 LCS 演算法找出共同行數
	common := longestCommonSubsequence(oldLines, newLines)
	added = len(newLines) - common
	removed = len(oldLines) - common
	return
}

// longestCommonSubsequence 計算兩份行陣列的最長共同子序列長度。
func longestCommonSubsequence(a, b []string) int {
	m, n := len(a), len(b)
	if m == 0 || n == 0 {
		return 0
	}
	// 使用滾動陣列節省記憶體
	prev := make([]int, n+1)
	curr := make([]int, n+1)
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				left := curr[j-1]
				up := prev[j]
				if left > up {
					curr[j] = left
				} else {
					curr[j] = up
				}
			}
		}
		prev, curr = curr, prev
	}
	return prev[n]
}

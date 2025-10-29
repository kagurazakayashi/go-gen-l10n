// gen-l10n/main.go
package main

import (
	"flag"
	"log"
	"path/filepath"
)

// main 為程式進入點，負責解析命令列參數、掃描指定目錄中的 ARB 語系檔案，並產生對應的 Go 程式碼。
func main() {
	// dirPtr 為 l10n 目錄路徑參數，用於指定 ARB 語系檔所在目錄。
	dirPtr := flag.String("dir", "./l10n", "l10n dir")
	// pkgPtr 為輸出 Go 檔案所使用的套件名稱參數。
	pkgPtr := flag.String("pkg", "l10n", "pkg name")
	flag.Parse()

	// 取得命令列解析後的目錄路徑與套件名稱。
	dir := *dirPtr
	pkgName := *pkgPtr

	// 搜尋指定目錄下所有符合 app_*.arb 命名規則的檔案。
	files, err := filepath.Glob(filepath.Join(dir, "app_*.arb"))
	if err != nil {
		log.Fatalf("[main] 查找 ARB 檔案失敗：%v", err)
	}
	if len(files) == 0 {
		log.Fatalf("[main] 在目錄 %s 中找不到 app_*.arb 檔案", dir)
	}

	// 輸出本次執行的主要參數，便於除錯與確認輸入設定是否正確。
	log.Printf("[main] 執行參數：dir=%s, pkg=%s", dir, pkgName)

	var locales []LocaleData
	var keys []KeyMeta
	defaultLocale := ""

	// 逐一解析每個 ARB 檔案，收集語系資料與對應的鍵值中繼資料。
	for _, file := range files {
		localeId, localeData, keysAdd := loadArbFile(file)
		locales = append(locales, localeData)
		keys = append(keys, keysAdd...)

		// 將第一個成功載入的語系列為預設語系，用於後續產生預設結構體名稱後綴。
		if defaultLocale == "" {
			defaultLocale = localeId
		}
	}

	// 組合程式碼產生器所需的樣板資料。
	tmplData := TemplateData{
		PackageName:         pkgName,
		Keys:                keys,
		Locales:             locales,
		DefaultStructSuffix: toCamelCase(defaultLocale),
	}

	// 輸出最終樣板資料，便於檢查產生器輸入內容。
	log.Printf("[main] 樣板資料內容：%+v", tmplData)

	// 依據整理後的樣板資料產生對應的 Go 程式碼。
	generateGoCode(dir, pkgName, tmplData)
}

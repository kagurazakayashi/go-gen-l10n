// gen-l10n/main.go
//
//go:generate goversioninfo -o=resource_windows_386.syso -64=false -icon=ico/icon.ico -manifest=main.exe.manifest
//go:generate goversioninfo -o=resource_windows_amd64.syso -64=true -icon=ico/icon.ico -manifest=main.exe.manifest
//go:generate goversioninfo -o=resource_windows_arm64.syso -arm=true -icon=ico/icon.ico -manifest=main.exe.manifest
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/kagurazakayashi/go-gen-l10n/l10n"
)

// main 為程式進入點，負責解析命令列參數、掃描指定目錄中的 ARB 語系檔案，
// 並產生對應的 Go 程式碼。
func main() {
	// dirPtr 為 l10n 目錄路徑參數，用於指定 ARB 語系檔所在目錄。
	dirPtr := flag.String("dir", "./l10n", "l10n dir")
	// pkgPtr 為輸出 Go 檔案所使用的套件名稱參數。
	pkgPtr := flag.String("pkg", "l10n", "pkg name")
	// langPtr 為程式自身輸出訊息所使用的語言。
	langPtr := flag.String("lang", "en", "program output language")
	flag.Parse()

	// 取得命令列解析後的目錄路徑、套件名稱與語言。
	dir := *dirPtr
	pkgName := *pkgPtr
	lang := *langPtr

	L := l10n.GetLocalizations(lang)

	// 搜尋指定目錄下所有符合 app_*.arb 命名規則的檔案。
	files, err := filepath.Glob(filepath.Join(dir, "app_*.arb"))
	if err != nil {
		log.Fatalf(L.ErrorFindArbFiles(), err)
	}
	if len(files) == 0 {
		log.Fatalf(L.ErrorNoArbFiles(), dir)
	}

	// 輸出本次執行的主要參數，便於除錯與確認輸入設定是否正確。
	log.Printf(L.InfoExecutionParams(), dir, pkgName)

	var locales []LocaleData
	var keys []KeyMeta
	keySet := make(map[string]bool)
	defaultLocale := ""

	// 逐一解析每個 ARB 檔案，收集語系資料與對應的鍵值中繼資料。
	for _, file := range files {
		localeId, localeData, keysAdd := loadArbFile(file, L)
		locales = append(locales, localeData)

		// 使用全域 keySet 去重，避免相同鍵在多個 ARB 檔案間重複。
		for _, k := range keysAdd {
			if !keySet[k.Key] {
				keySet[k.Key] = true
				keys = append(keys, k)
			}
		}

		// 將第一個成功載入的語系列為預設語系，用於後續產生預設結構體名稱後綴。
		if defaultLocale == "" {
			defaultLocale = localeId
		}
	}

	// 組合程式碼產生器所需的範本資料。
	tmplData := TemplateData{
		PackageName:         pkgName,
		Keys:                keys,
		Locales:             locales,
		DefaultStructSuffix: toCamelCase(defaultLocale),
		GeneratedLocale:     lang,
	}

	// 依據整理後的範本資料產生對應的 Go 程式碼。
	generateGoCode(dir, pkgName, tmplData, L)

	// 成功結束程式，回傳 0 給作業系統。
	os.Exit(0)
}

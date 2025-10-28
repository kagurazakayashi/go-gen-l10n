// gen-l10n/main.go
package main

import (
	"flag"
	"log"
	"path/filepath"
)

// main 為程式進入點，負責解析命令列參數並檢查指定目錄中的 ARB 語系檔案。
func main() {
	// dirPtr 為 l10n 目錄路徑參數，用於指定 ARB 檔案所在目錄。
	dirPtr := flag.String("dir", "./l10n", "l10n dir")
	// pkgPtr 為產生之 Go 檔案的套件名稱參數。
	pkgPtr := flag.String("pkg", "l10n", "pkg name")
	flag.Parse()

	// 取得解析後的目錄路徑與套件名稱。
	dir := *dirPtr
	pkgName := *pkgPtr

	// 比對指定目錄下所有符合 app_*.arb 命名規則的檔案。
	files, err := filepath.Glob(filepath.Join(dir, "app_*.arb"))
	if err != nil {
		log.Fatalf("[main] 查找 ARB 檔案失敗：%v", err)
	}
	if len(files) == 0 {
		log.Fatalf("[main] 在目錄 %s 中找不到 app_*.arb 檔案", dir)
	}

	// 輸出目前執行參數，便於除錯與確認輸入設定。
	log.Printf("[main] dir=%s, pkg=%s", dir, pkgName)

	var locales []LocaleData
	// 解析每一個 arb 檔案
	for _, file := range files {
		locales = append(locales, loadArbFile(file))
	}

	log.Printf("%v", locales)
}

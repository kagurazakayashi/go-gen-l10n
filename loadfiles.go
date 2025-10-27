package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadFiles(files []string) {
	// keySet := make(map[string]bool)
	// var keys []KeyMeta
	// var locales []LocaleData
	defaultLocale := ""

	// 解析每一個 arb 檔案
	for _, file := range files {
		// 取得檔案名稱（不含路徑），例如：app_zh_TW.arb
		baseName := filepath.Base(file)

		log.Println(baseName)

		// 從檔名中擷取語系識別碼，會將前綴 "app_" 與副檔名 ".arb" 去除。
		// 例如：app_zh_TW.arb -> zh_TW
		localeId := strings.TrimSuffix(strings.TrimPrefix(baseName, "app_"), ".arb")

		// 若尚未指定預設語系，則以第一個讀取到的語系作為預設值。
		if defaultLocale == "" {
			defaultLocale = localeId
		}

		// 讀取目前 ARB 檔案內容。
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("[I18nLoader] 讀取檔案 %s 失敗: %v", file, err)
		}

		// 將 JSON 內容反序列化為通用 map，供後續欄位解析使用。
		var rawMap map[string]interface{}
		if err := json.Unmarshal(content, &rawMap); err != nil {
			log.Fatalf("[I18nLoader] 解析 %s 失敗: %v", file, err)
		}
		log.Println(rawMap)

		// 初始化目前語系對應的翻譯字典。
		translations := make(map[string]string)

		// 輸出目前處理的語系與初始化後的翻譯筆數，方便追蹤載入流程。
		log.Println(translations)
	}
}

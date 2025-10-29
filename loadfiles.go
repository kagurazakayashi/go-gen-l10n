package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// loadArbFile 讀取指定的 ARB 檔案，解析語系識別碼、翻譯內容與鍵值中繼資料。
func loadArbFile(file string) (string, LocaleData, []KeyMeta) {
	// 取得檔案名稱（不含路徑），例如：app_zh_TW.arb。
	baseName := filepath.Base(file)

	// 從檔名中擷取語系識別碼，會去除前綴 "app_" 與副檔名 ".arb"。
	// 例如：app_zh_TW.arb -> zh_TW。
	localeId := strings.TrimSuffix(strings.TrimPrefix(baseName, "app_"), ".arb")

	// 讀取目前 ARB 檔案內容。
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("[I18nLoader] 讀取 ARB 檔案失敗，檔案路徑=%s，錯誤=%v", file, err)
	}

	// 將 JSON 內容反序列化為通用 map，供後續欄位解析使用。
	var rawMap map[string]interface{}
	if err := json.Unmarshal(content, &rawMap); err != nil {
		log.Fatalf("[I18nLoader] 解析 ARB 檔案失敗，檔案路徑=%s，錯誤=%v", file, err)
	}

	// 將原始 ARB 資料轉換為鍵值中繼資料與語系資料。
	keys, localeData := arbMap(localeId, rawMap)
	return localeId, localeData, keys
}

// arbMap 將 ARB 原始資料轉換為鍵值清單與語系資料結構。
func arbMap(localeId string, rawMap map[string]interface{}) ([]KeyMeta, LocaleData) {
	// keySet 用於避免重複加入相同的翻譯鍵。
	keySet := make(map[string]bool)
	var keys []KeyMeta

	// translations 儲存實際可用的翻譯內容，值會先轉為帶引號的字串格式。
	translations := make(map[string]string)
	for k, v := range rawMap {
		// 略過 ARB 中以 @ 開頭的中繼資料欄位。
		if strings.HasPrefix(k, "@") {
			continue
		}

		// 僅處理字串型別的翻譯值。
		if strVal, ok := v.(string); ok {
			translations[k] = strconv.Quote(strVal)

			// 若尚未記錄該鍵，則建立對應的中繼資料。
			if !keySet[k] {
				keySet[k] = true
				keys = append(keys, KeyMeta{
					Key:        k,
					MethodName: toCamelCase(k),
				})
			}
		}
	}

	// 組合語系資料，供後續程式碼產生器或存取邏輯使用。
	localeData := LocaleData{
		ID:           localeId,
		StructSuffix: toCamelCase(localeId),
		Translations: translations,
	}
	return keys, localeData
}

// toCamelCase 將底線或連字號分隔的字串轉為首字母大寫的 CamelCase 格式。
func toCamelCase(s string) string {
	// 依底線或連字號切分字串片段。
	parts := regexp.MustCompile(`[_\-]+`).Split(s, -1)
	for i, part := range parts {
		if len(part) > 0 {
			// 將每個片段的首個字元轉為大寫。
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}

	// 將所有片段串接為 CamelCase 字串。
	return strings.Join(parts, "")
}

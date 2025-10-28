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

func loadArbFile(file string) LocaleData {
	defaultLocale := ""

	// 取得檔案名稱（不含路徑），例如：app_zh_TW.arb
	baseName := filepath.Base(file)

	log.Println(baseName)

	// 從檔名中擷取語系識別碼，會將前綴 "app_" 與副檔名 ".arb" 去除。
	// 例如：app_zh_TW.arb -> zh_TW
	var localeId string = strings.TrimSuffix(strings.TrimPrefix(baseName, "app_"), ".arb")

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

	return arbMap(localeId, rawMap)
}

func arbMap(localeId string, rawMap map[string]interface{}) LocaleData {
	keySet := make(map[string]bool)
	var keys []KeyMeta

	translations := make(map[string]string)
	for k, v := range rawMap {
		if strings.HasPrefix(k, "@") {
			continue
		}
		if strVal, ok := v.(string); ok {
			translations[k] = strconv.Quote(strVal)
			if !keySet[k] {
				keySet[k] = true
				keys = append(keys, KeyMeta{
					Key:        k,
					MethodName: toCamelCase(k),
				})
			}
		}
	}
	return LocaleData{
		ID:           localeId,
		StructSuffix: toCamelCase(localeId),
		Translations: translations,
	}
}

func toCamelCase(s string) string {
	parts := regexp.MustCompile(`[_\-]+`).Split(s, -1)
	for i, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}

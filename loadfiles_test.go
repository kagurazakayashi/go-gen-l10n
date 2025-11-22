package main

import (
	"strconv"
	"testing"

	"github.com/kagurazakayashi/go-gen-l10n/l10n"
)

// mockL10n 提供測試用的 AppLocalizations 實作。
// 所有方法回傳英文格式字串，確保測試輸出一致。
type mockL10n struct{}

func (m *mockL10n) CommentAppLocalizationsInterface() string { return "AppLocalizations interface" }
func (m *mockL10n) CommentGetLocalizations() string          { return "Get localizations" }
func (m *mockL10n) CommentLocaleImplementation() string      { return "Locale implementation: %s" }
func (m *mockL10n) CliFlagDir() string                       { return "l10n dir" }
func (m *mockL10n) CliFlagPkg() string                       { return "package name" }
func (m *mockL10n) ErrorFindArbFiles() string                { return "Failed to find ARB files: %v" }
func (m *mockL10n) ErrorNoArbFiles() string                  { return "No app_*.arb files found in directory %s" }
func (m *mockL10n) ErrorReadArbFile() string                 { return "Failed to read ARB file: path=%s, error=%v" }
func (m *mockL10n) ErrorParseArbFile() string                { return "Failed to parse ARB file: path=%s, error=%v" }
func (m *mockL10n) ErrorParseBaseTemplate() string           { return "Failed to parse base template: %v" }
func (m *mockL10n) ErrorExecuteBaseTemplate() string         { return "Failed to execute base template: %v" }
func (m *mockL10n) ErrorFormatBaseCode() string              { return "Failed to format base code: %v\nCode content:\n%s" }
func (m *mockL10n) ErrorWriteBaseFile() string               { return "Failed to write base file: %v" }
func (m *mockL10n) ErrorParseLocaleTemplate() string         { return "Failed to parse locale template: %v" }
func (m *mockL10n) ErrorExecuteLocaleTemplate() string       { return "Failed to execute locale template (%s): %v" }
func (m *mockL10n) ErrorFormatLocaleCode() string            { return "Failed to format locale code (%s): %v\nCode content:\n%s" }
func (m *mockL10n) ErrorWriteLocaleFile() string             { return "Failed to write locale file (%s): %v" }
func (m *mockL10n) InfoExecutionParams() string              { return "Execution parameters: dir=%s, pkg=%s" }
func (m *mockL10n) InfoTemplateData() string                 { return "Template data: %+v" }
func (m *mockL10n) SuccessGeneratedCode() string             { return "Generated: %s" }

// 確認 mockL10n 實作了 AppLocalizations 介面
var _ l10n.AppLocalizations = (*mockL10n)(nil)

// TestToCamelCase 測試 toCamelCase 對各種輸入格式的轉換結果。
// 此測試示範了如何將 ARB 中的 snake_case 鍵轉換為 PascalCase 方法名稱。
func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple key",
			input: "hello",
			want:  "Hello",
		},
		{
			name:  "snake_case key",
			input: "hello_world",
			want:  "HelloWorld",
		},
		{
			name:  "kebab-case locale",
			input: "zh-TW",
			want:  "ZhTW",
		},
		{
			name:  "underscore locale",
			input: "zh_Hant",
			want:  "ZhHant",
		},
		{
			name:  "dot not treated as separator",
			input: "user.login.title",
			want:  "User.login.title",
		},
		{
			name:  "mixed separators",
			input: "some_key-with-mixed_separators",
			want:  "SomeKeyWithMixedSeparators",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "already PascalCase",
			input: "HelloWorld",
			want:  "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestArbMap 測試 ARB 原始資料到 LocaleData 的轉換。
// 此測試示範了 @ 前綴的中繼資料鍵如何被忽略，以及數字等非字串值如何處理。
func TestArbMap(t *testing.T) {
	tests := []struct {
		name           string
		localeId       string
		rawMap         map[string]interface{}
		wantKeys       int
		checkCamelCase []struct {
			key        string
			methodName string
		}
	}{
		{
			name:     "basic ARB with translation keys",
			localeId: "en",
			rawMap: map[string]interface{}{
				"hello":   "Hello",
				"welcome": "Welcome!",
			},
			wantKeys: 2,
			checkCamelCase: []struct {
				key        string
				methodName string
			}{
				{"hello", "Hello"},
				{"welcome", "Welcome"},
			},
		},
		{
			name:     "ARB with @@locale and @description metadata that should be skipped",
			localeId: "zh",
			rawMap: map[string]interface{}{
				"@@locale": "zh",
				"hello":    "你好",
				"@hello":   map[string]interface{}{"description": "greeting"},
				"welcome":  "欢迎",
				"@welcome": map[string]interface{}{"description": "welcome message"},
			},
			wantKeys: 2,
			checkCamelCase: []struct {
				key        string
				methodName string
			}{
				{"hello", "Hello"},
				{"welcome", "Welcome"},
			},
		},
		{
			name:     "ARB with non-string values are skipped",
			localeId: "en",
			rawMap: map[string]interface{}{
				"hello":   "Hello",
				"count":   42,
				"enabled": true,
				"welcome": "Welcome!",
				"price":   9.99,
			},
			wantKeys: 2,
		},
		{
			name:     "locale with underscore creates correct StructSuffix",
			localeId: "zh_Hant",
			rawMap: map[string]interface{}{
				"hello": "你好",
			},
			wantKeys: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, localeData := arbMap(tt.localeId, tt.rawMap)

			if localeData.ID != tt.localeId {
				t.Errorf("LocaleData.ID = %q, want %q", localeData.ID, tt.localeId)
			}
			if localeData.StructSuffix != toCamelCase(tt.localeId) {
				t.Errorf("LocaleData.StructSuffix = %q, want %q",
					localeData.StructSuffix, toCamelCase(tt.localeId))
			}
			if len(keys) != tt.wantKeys {
				t.Errorf("got %d keys, want %d: keys=%v", len(keys), tt.wantKeys, keys)
			}
			if len(localeData.Translations) != tt.wantKeys {
				t.Errorf("got %d translations, want %d", len(localeData.Translations), tt.wantKeys)
			}

			// 驗證翻譯值經過 strconv.Quote 轉義
			for _, k := range keys {
				trans, ok := localeData.Translations[k.Key]
				if !ok {
					t.Errorf("missing translation for key %q", k.Key)
					continue
				}
				expectedTrans := strconv.Quote(tt.rawMap[k.Key].(string))
				if trans != expectedTrans {
					t.Errorf("Translation[%q] = %s, want %s", k.Key, trans, expectedTrans)
				}
			}

			// 驗證 PascalCase 轉換
			for _, cc := range tt.checkCamelCase {
				found := false
				for _, k := range keys {
					if k.Key == cc.key {
						found = true
						if k.MethodName != cc.methodName {
							t.Errorf("MethodName for %q = %q, want %q",
								k.Key, k.MethodName, cc.methodName)
						}
					}
				}
				if !found {
					t.Errorf("key %q not found in keys", cc.key)
				}
			}

			// 驗證 @ 前綴的鍵不會被包含
			for _, k := range keys {
				if k.Key[0] == '@' {
					t.Errorf("key with @ prefix should be filtered: %q", k.Key)
				}
			}
		})
	}
}

// TestArbMapTranslationQuoting 驗證 strconv.Quote 對特殊字元的轉義，
// 確保產生的 Go 程式碼語法正確。
func TestArbMapTranslationQuoting(t *testing.T) {
	rawMap := map[string]interface{}{
		"quoted":    `He said "hello"`,
		"newline":   "line1\nline2",
		"backslash": `C:\path\to\file`,
		"tab":       "col1\tcol2",
	}
	keys, localeData := arbMap("en", rawMap)

	for _, k := range keys {
		trans := localeData.Translations[k.Key]
		// 透過 strconv.Unquote 驗證 Go 字串的有效性
		unquoted, err := strconv.Unquote(trans)
		if err != nil {
			t.Errorf("Translation for %q is not a valid Go string: %v\nRaw: %s", k.Key, err, trans)
			continue
		}
		// 驗證 unquote 後與原始值一致
		original := rawMap[k.Key].(string)
		if unquoted != original {
			t.Errorf("Unquote mismatch for %q: got %q, want %q", k.Key, unquoted, original)
		}
	}
}

// TestArbMapDeduplication 驗證 arbMap 在處理 rawMap 時不會產生重複的鍵。
func TestArbMapDeduplication(t *testing.T) {
	rawMap := map[string]interface{}{
		"hello":   "Hello",
		"welcome": "Welcome",
	}
	keys, _ := arbMap("en", rawMap)
	if len(keys) != 2 {
		t.Errorf("expected 2 unique keys, got %d", len(keys))
	}
}

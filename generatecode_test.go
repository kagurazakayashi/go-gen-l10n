package main

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// TestBaseTemplateRendering 驗證基礎範本的渲染結果。
// 此測試示範了 AppLocalizations 介面和 GetLocalizations 函式是如何產生的。
func TestBaseTemplateRendering(t *testing.T) {
	data := map[string]interface{}{
		"PackageName":                    "testpkg",
		"Keys":                           []KeyMeta{{Key: "hello", MethodName: "Hello"}, {Key: "bye", MethodName: "Bye"}},
		"Locales":                        []LocaleData{{ID: "en", StructSuffix: "En"}, {ID: "ja", StructSuffix: "Ja"}},
		"DefaultStructSuffix":            "En",
		"CommentAppLocalizationsInterface": "Demo interface",
		"CommentGetLocalizations":        "Demo getter",
	}

	tmpl, err := template.New("base").Parse(BaseTemplate)
	if err != nil {
		t.Fatalf("parse base template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("execute base template: %v", err)
	}

	output := buf.String()

	// 驗證關鍵內容是否正確輸出
	checks := []string{
		"package testpkg",
		"// Demo interface",
		"type AppLocalizations interface {",
		"Hello() string",
		"Bye() string",
		"// Demo getter",
		"func GetLocalizations(",
		`case "en":`,
		`case "ja":`,
		"&appLocalizationsEn{}",
		"&appLocalizationsJa{}",
		"// Fallback",
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Errorf("base template output missing %q", want)
		}
	}

	// 驗證產生的程式碼可以通過 go/format 格式化
	_, err = format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("format.Source on base template output: %v\n=== output ===\n%s", err, output)
	}
}

// TestLocaleTemplateRendering 驗證單一語言實作檔案的範本渲染結果。
// 此測試示範了每個語言的 struct 和方法是如何產生的。
func TestLocaleTemplateRendering(t *testing.T) {
	data := map[string]interface{}{
		"PackageName":              "testpkg",
		"Keys":                     []KeyMeta{{Key: "hello", MethodName: "Hello"}, {Key: "bye", MethodName: "Bye"}},
		"Locale":                   LocaleData{ID: "ja", StructSuffix: "Ja", Translations: map[string]string{"hello": `"こんにちは"`, "bye": `"さようなら"`}},
		"CommentLocaleImplementation": "Locale: ja",
	}

	tmpl, err := template.New("locale").Parse(LocaleTemplate)
	if err != nil {
		t.Fatalf("parse locale template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("execute locale template: %v", err)
	}

	output := buf.String()

	checks := []string{
		"package testpkg",
		"// Locale: ja",
		"type appLocalizationsJa struct{}",
		"func (l *appLocalizationsJa) Hello() string {",
		`return "こんにちは"`,
		"func (l *appLocalizationsJa) Bye() string {",
		`return "さようなら"`,
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Errorf("locale template output missing %q", want)
		}
	}

	_, err = format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("format.Source on locale template output: %v\n=== output ===\n%s", err, output)
	}
}

// TestLocaleTemplateMissingTranslation 驗證當某個語言缺少翻譯鍵時的回退行為。
// 缺少的鍵會回退到回傳原始 key 名稱。
func TestLocaleTemplateMissingTranslation(t *testing.T) {
	keys := []KeyMeta{
		{Key: "hello", MethodName: "Hello"},
		{Key: "missing_key", MethodName: "MissingKey"},
	}
	translations := map[string]string{
		"hello": `"Hello"`,
		// 故意不包含 missing_key 的翻譯
	}
	data := map[string]interface{}{
		"PackageName":              "testpkg",
		"Keys":                     keys,
		"Locale":                   LocaleData{ID: "en", StructSuffix: "En", Translations: translations},
		"CommentLocaleImplementation": "Locale: en",
	}

	tmpl, err := template.New("locale").Parse(LocaleTemplate)
	if err != nil {
		t.Fatalf("parse locale template: %v", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("execute locale template: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, `return "Hello"`) {
		t.Error("missing translation for hello key")
	}
	if !strings.Contains(output, `return "missing_key"`) {
		t.Error("missing translation should fallback to returning the key name")
	}

	_, err = format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("format.Source on locale template output: %v\n=== output ===\n%s", err, output)
	}
}

// TestFullGeneration 端到端測試：建立暫存 ARB 檔案，執行完整產生流程，驗證輸出。
// 此測試模擬了使用者在專案中執行 go-gen-l10n 的完整使用情境。
func TestFullGeneration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-gen-l10n-test")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 建立英文 ARB 檔案
	enARB := `{
	"@@locale": "en",
	"hello": "Hello",
	"greeting": "Welcome, friend!"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "app_en.arb"), []byte(enARB), 0644); err != nil {
		t.Fatalf("write app_en.arb: %v", err)
	}

	// 建立中文 ARB 檔案
	zhARB := `{
	"@@locale": "zh",
	"hello": "你好",
	"greeting": "欢迎，朋友！"
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "app_zh.arb"), []byte(zhARB), 0644); err != nil {
		t.Fatalf("write app_zh.arb: %v", err)
	}

	// 驗證 ARB 檔案能被正確載入
	L := &mockL10n{}
	_, localeDataEn, keysEn := loadArbFile(filepath.Join(tmpDir, "app_en.arb"), L)
	_, localeDataZh, keysZh := loadArbFile(filepath.Join(tmpDir, "app_zh.arb"), L)

	if localeDataEn.ID != "en" {
		t.Errorf("English locale ID = %q, want %q", localeDataEn.ID, "en")
	}
	if localeDataZh.ID != "zh" {
		t.Errorf("Chinese locale ID = %q, want %q", localeDataZh.ID, "zh")
	}
	if len(keysEn) != 2 {
		t.Errorf("English keys count = %d, want 2", len(keysEn))
	}
	if len(keysZh) != 2 {
		t.Errorf("Chinese keys count = %d, want 2", len(keysZh))
	}

	// 組合範本資料
	tmplData := TemplateData{
		PackageName:         "testing",
		Keys:                keysEn,
		Locales:             []LocaleData{localeDataEn, localeDataZh},
		DefaultStructSuffix: "En",
		GeneratedLocale:     "en",
	}

	// 產生程式碼到暫存目錄
	generateGoCode(tmpDir, "testing", tmplData, L)

	// 驗證產生的檔案存在
	expectedFiles := []string{
		"app_localizations.go",
		"app_localizations_en.go",
		"app_localizations_zh.go",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file not generated: %s", f)
		}
	}

	// 驗證產生程式碼可以編譯
	baseContent, err := os.ReadFile(filepath.Join(tmpDir, "app_localizations.go"))
	if err != nil {
		t.Fatalf("read app_localizations.go: %v", err)
	}
	validateGeneratedCode(t, baseContent, "app_localizations.go")

	enContent, err := os.ReadFile(filepath.Join(tmpDir, "app_localizations_en.go"))
	if err != nil {
		t.Fatalf("read app_localizations_en.go: %v", err)
	}
	validateGeneratedCode(t, enContent, "app_localizations_en.go")
	if !strings.Contains(string(enContent), `"Hello"`) {
		t.Error("English file missing 'Hello' translation")
	}

	zhContent, err := os.ReadFile(filepath.Join(tmpDir, "app_localizations_zh.go"))
	if err != nil {
		t.Fatalf("read app_localizations_zh.go: %v", err)
	}
	validateGeneratedCode(t, zhContent, "app_localizations_zh.go")
	if !strings.Contains(string(zhContent), `"你好"`) {
		t.Error("Chinese file missing '你好' translation")
	}
}

// validateGeneratedCode 驗證產生的程式碼是否為合法的 Go 程式碼。
func validateGeneratedCode(t *testing.T, code []byte, filename string) {
	t.Helper()
	_, err := format.Source(code)
	if err != nil {
		t.Errorf("%s: generated code is not valid Go:\n%v\n=== code ===\n%s", filename, err, string(code))
	}
}

package main

import (
	"fmt"
	"testing"

	"github.com/kagurazakayashi/go-gen-l10n/l10n"
)

// TestGetLocalizationsKnownLocales 驗證 GetLocalizations 對已知語言程式碼的處理。
// 此測試示範了語言程式碼與具體實作型別的對應關係。
func TestGetLocalizationsKnownLocales(t *testing.T) {
	tests := []struct {
		locale string
		want   string
	}{
		{"en", "*l10n.appLocalizationsEn"},
		{"ja", "*l10n.appLocalizationsJa"},
		{"zh", "*l10n.appLocalizationsZh"},
		{"zh_Hant", "*l10n.appLocalizationsZhHant"},
	}
	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			got := l10n.GetLocalizations(tt.locale)
			typeStr := fmt.Sprintf("%T", got)
			if typeStr != tt.want {
				t.Errorf("GetLocalizations(%q) type = %s, want %s", tt.locale, typeStr, tt.want)
			}
		})
	}
}

// TestGetLocalizationsFallback 驗證未知語言程式碼回退到預設語言的行為。
// 當傳入未支援的語言程式碼時，應回傳第一個載入的 ARB 檔案對應的語言（此處為 en）。
func TestGetLocalizationsFallback(t *testing.T) {
	got := l10n.GetLocalizations("fr")
	typeStr := fmt.Sprintf("%T", got)
	if typeStr != "*l10n.appLocalizationsEn" {
		t.Errorf("GetLocalizations(\"fr\") type = %s, want *l10n.appLocalizationsEn", typeStr)
	}
}

// TestAppLocalizationsEnMethods 驗證英文翻譯實例回傳的實際值。
// 此測試示範了如何在程式中呼叫翻譯方法取得英文字串。
func TestAppLocalizationsEnMethods(t *testing.T) {
	l := l10n.GetLocalizations("en")

	if got := l.CliFlagDir(); got != "l10n dir" {
		t.Errorf("CliFlagDir() = %q, want %q", got, "l10n dir")
	}
	if got := l.CliFlagPkg(); got != "package name" {
		t.Errorf("CliFlagPkg() = %q, want %q", got, "package name")
	}
	if got := l.SuccessGeneratedCode(); got != "Generated localization code: %s" {
		t.Errorf("SuccessGeneratedCode() = %q", got)
	}
}

// TestAppLocalizationsZhMethods 驗證中文翻譯實例回傳中文翻譯。
// 此測試示範了多語系切換的實際效果。
func TestAppLocalizationsZhMethods(t *testing.T) {
	l := l10n.GetLocalizations("zh")

	if got := l.CliFlagDir(); got != "l10n 目录" {
		t.Errorf("CliFlagDir() = %q, want %q", got, "l10n 目录")
	}
	if got := l.SuccessGeneratedCode(); got != "成功生成本地化代码：%s" {
		t.Errorf("SuccessGeneratedCode() = %q", got)
	}
}

// TestAppLocalizationsJaMethods 驗證日文翻譯實例回傳日文翻譯。
// 此測試示範了日文語系的翻譯內容。
func TestAppLocalizationsJaMethods(t *testing.T) {
	l := l10n.GetLocalizations("ja")

	if got := l.CliFlagDir(); got != "l10n ディレクトリ" {
		t.Errorf("CliFlagDir() = %q, want %q", got, "l10n ディレクトリ")
	}
	if got := l.CommentAppLocalizationsInterface(); got != "サポートされるすべてのローカライズされた文字列のインターフェースを定義します" {
		t.Errorf("CommentAppLocalizationsInterface() = %q", got)
	}
}

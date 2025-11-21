package main

// TemplateData 定義傳遞給 text/template 使用的資料結構。
// 此結構主要用於範本渲染階段，提供完整的國際化（i18n）相關資訊，
// 包含套件名稱、鍵值對應資訊、多語系資料以及預設結構後綴。
// 各欄位需在進入範本前完成初始化，以確保範本渲染正確無誤。
type TemplateData struct {
	PackageName         string       // 目標產生程式碼所屬的套件名稱（package name）
	Keys                []KeyMeta    // 所有翻譯鍵（key）對應的方法資訊清單
	Locales             []LocaleData // 多語系資料集合，每個元素代表一種語系
	DefaultStructSuffix string       // 預設結構名稱後綴（例如：En、ZhTw），用於產生型別名稱
	GeneratedLocale     string       // 產生程式碼時使用的程式語言
}

// KeyMeta 描述單一翻譯鍵（key）與其對應的方法名稱。
// 通常用於在範本中產生對應的存取方法（getter），
// 以提供型別安全的方式取得翻譯內容。
type KeyMeta struct {
	Key        string // 原始翻譯鍵（例如："user.login.title"）
	MethodName string // 對應產生的方法名稱（例如："UserLoginTitle"）
}

// LocaleData 表示單一語系的翻譯資料。
// 每個語系會對應一組 key-value 的翻譯內容，
// 並可透過 StructSuffix 區分不同語系產生的結構名稱。
type LocaleData struct {
	ID           string            // 語系識別碼（例如："en", "zh-TW"）
	StructSuffix string            // 該語系對應的結構名稱後綴（例如："En", "ZhTw"）
	Translations map[string]string // 翻譯內容對應表（key -> 已經過 strconv.Quote 轉義的 Go 字串）
}

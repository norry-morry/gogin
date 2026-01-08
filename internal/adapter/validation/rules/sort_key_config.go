// Package rules はカスタムバリデーションルールの登録・共通設定を提供します。
package rules

// SortKeyAllowed は sort_key=xxx ごとに許容されるソートキー一覧を定義します。
// xxx はバリデータタグの param（例: "identities"）です。
var SortKeyAllowed = map[string][]string{
	"identities": {
		"created_at",
		"provider",
		"uid",
		"email_at_signup",
	},

	// 将来の拡張イメージ:
	"addresses": {
		"purpose_id",
		"created_at",
		"postal_code",
		"country_code",
		"administrative_area",
		"locality",
		"address_line1",
	},
}

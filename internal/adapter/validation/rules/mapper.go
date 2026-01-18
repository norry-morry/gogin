// Package rules は、カスタムバリデーションルールとそのエラー変換ユーティリティを提供します。
// ここでは validator.ValidationErrors を i18n 対応のエラーマップ形式に変換します。
package rules

import (
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"

	"resume/internal/shared/util"
)

// ErrorItem は、1 件のバリデーションエラー項目を表す構造体です。
// 各エラーの翻訳キー（Code）と、埋め込みパラメータ（Params）を保持します。
type ErrorItem struct {
	Code   string         `json:"code"`
	Params map[string]any `json:"params"`
}

// normalizeRuleTagLowerCamel は、validator のルールタグを lowerCamel に正規化します。
// - snake/kebab/Pascal/camel/UPPER いずれの入力でも受ける
// - go-playground/validator の一部タグ表記（"startswith"/"endswith"）も camel に寄せる
func normalizeRuleTagLowerCamel(tag string) string {
	t := strings.TrimSpace(tag)
	// まず kebab を snake に
	t = strings.ReplaceAll(t, "-", "_")

	// validator の素のタグが "startswith"/"endswith" のため補正
	switch t {
	case "startswith":
		t = "starts_with"
	case "endswith":
		t = "ends_with"
	}

	// 何が来ても lowerCamel に統一
	return util.ToCamel(t)
}

// region KeyStore

// MessageKeyStore は「キーが辞書に存在するかだけ」を見るためのインターフェース。
type MessageKeyStore interface {
	HasKey(key string) bool
}

var msgKeyStore MessageKeyStore

// SetMessageKeyStore は、バリデーションメッセージ用のキー存在チェックに使用する
// MessageKeyStore 実装をグローバルに設定します。
func SetMessageKeyStore(s MessageKeyStore) {
	msgKeyStore = s
}

func hasValidationKey(key string) bool {
	if msgKeyStore == nil {
		slog.Debug("hasValidationKey: msgKeyStore is nil", "key", key)
		return false
	}
	ok := msgKeyStore.HasKey(key)
	slog.Debug("hasValidationKey", "key", key, "ok", ok)
	return ok
}

// endregion

// region Key Builders

// rulesKey は、通常のバリデーションルール用のメッセージキーを生成します。
func rulesKey(tagLowerCamel string) string { return "validation.rules." + tagLowerCamel }

// customKey は、カスタムバリデーションルール用のメッセージキーを生成します。
func customKey(tagLowerCamel string) string { return "validation.custom." + tagLowerCamel }

// fieldSpecificKey は、フィールド固有のバリデーションメッセージキーを生成します。
func fieldSpecificKey(scope, fieldKeyLowerCamel, tagLowerCamel string) string {
	return "validation.fieldspecific." + scope + "." + fieldKeyLowerCamel + "." + tagLowerCamel
}

// endregion

// MapValidationErrors は validator.ValidationErrors を i18n 対応の
// map[string][]ErrorItem 形式に変換します。
// scope にはドメインスコープ（例: "domain.address"）を指定します。
func MapValidationErrors(
	verrs validator.ValidationErrors,
	scope string, // 例: "domain.address"
) map[string][]ErrorItem {

	out := map[string][]ErrorItem{}

	for _, fe := range verrs {
		fieldLowerCamel := util.ToCamel(fe.Field()) // ex) address_line1 -> addressLine1
		base := scope + "." + fieldLowerCamel       // ex) domain.address.addressLine1
		//legacyKey := "validation.fields." + base    // 旧互換
		labelKey := base + ".label" // 推奨: domain.*.label（存在判定しないでキーをそのまま渡す）
		//fieldParam := labelKey                                // デフォは labelKey
		tagLowerCamel := normalizeRuleTagLowerCamel(fe.Tag()) // ex) startswith -> startsWith

		// ★ sort_key は oneof として扱う（コードとテンプレ両方を oneof に統一）
		if tagLowerCamel == "sortKey" {
			tagLowerCamel = "oneof"
		}

		tagSnake := util.ToSnake(tagLowerCamel)

		// 直接キーを渡す（フロント側で解決）。未定義時はそのままキー表示になるが、Humanizeより整合的。
		//fieldParam := labelKey

		// パラメータ（辞書側で {field}, {param}, {min}, {max}, {len}, {value} などを利用）
		p := map[string]any{"field": labelKey}

		if param := strings.TrimSpace(fe.Param()); param != "" {
			p["param"] = param

			// ★ sort_key = identities の時、allowedSortKeys["identities"] を values に展開する
			if tagLowerCamel == "oneof" { // sort_key → oneof に統一されている
				if keys, ok := SortKeyAllowed[param]; ok && len(keys) > 0 {
					// "created_at provider uid email_at_signup"
					joined := strings.Join(keys, " ")
					p["values"] = joined // ★ フロントが欲しかった最終形
					p["oneof"] = joined  // 従来 param の代わりに oneof をより正確に
				} else {
					// フォールバック（パラメータ名だけ）
					p["values"] = param
					p["oneof"] = param
				}
			}

			switch tagLowerCamel {
			case "min":
				p["min"] = param
			case "max":
				p["max"] = param
			case "len":
				p["len"] = param
			case "oneof":
				p["oneof"] = param
			case "gte", "lte", "gt", "lt", "gtfield", "gtefield", "ltfield", "ltefield":
				p["value"] = param
				// ToDo: gt,gte,lt,lteの場合は、paramにドメインキーも付与して辞書に差し込みたいかも。(フロントと接続後検証してチケット化)
			}
		}

		// 3種類の候補キーを生成
		fsKey := fieldSpecificKey(scope, fieldLowerCamel, tagLowerCamel)
		cKey := customKey(tagLowerCamel)
		rKey := rulesKey(tagLowerCamel)
		//key := rulesKey(tagLowerCamel)

		// --- 検索用（辞書内に存在するかを見る）キー: snake ベース ---
		//     ※ builder は文字列を連結しているだけなので、最後の引数に snake を渡しても問題なし。
		fsKeyLookup := fieldSpecificKey(scope, fieldLowerCamel, tagSnake)
		cKeyLookup := customKey(tagSnake)
		// rKeyLookup := rulesKey(tagSnake) // 今回は rules は存在チェックしなくてもOKなら未使用でよい

		// 優先度:フィールド固有 > custom > rules
		key := rKey
		if hasValidationKey(cKeyLookup) {
			key = cKey
		}
		if hasValidationKey(fsKeyLookup) {
			key = fsKey
		}

		out[base] = append(
			out[base],
			ErrorItem{
				Code:   key,
				Params: p,
			},
		)
	}
	return out
}

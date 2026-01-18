「辞書で多言語化（コード＋パラメータ方式）」に切り替えるなら、translator は adapter 側に“残しても良いが基本は不要”、多言語化の本体は shared に置くのが正解 です。

どう分けるか（指針）

shared（横断関心・純粋ロジック）

i18n ストア/ロケール判定/ETag：

validator → i18n コード化マッパー：internal/shared/validate/*

例：MapValidationErrors()（validator.ValidationErrors → {code, params}）

ここでタグ→validation.rules.* のキー決定と params(min/max/field) を整形

adapter（外部ライブラリ結線・入出力境界）

既存の internal/adapter/validation/translator は go-playground の “内蔵翻訳” を使う層
→ 辞書ベースに移行するなら原則使わない（二重管理になるため）。
→ 互換のため当面は残してもOK（err.Translate(trans) を呼ぶ旧コードがあるなら）。

なぜ shared/validate に置くの？

バックのどの層（REST/GRPC/CLI）でも 同じ i18n キー を使える（adapterに引きずられない）

Clean Architecture 的にも「ドメインに依らない横断ロジック」は shared が収まり良い

後で カスタムルール を足すときも shared/validate に1か所追加で済む

具体アクション（最小差分）

internal/shared/validate/map.go（新規）

さっき提案した MapValidationErrors() をここに置く

toRuleCodeAndParams() で標準タグ（required/min/max/len/email…）を validation.rules.* にマップ

params.field は i18n.ResolveFieldName() でラベル名に解決

既存の internal/adapter/validation/translator は…

新規実装では呼ばない（jatrans.RegisterDefaultTranslations も不要）

旧コードで err.Translate(Trans) を使っている箇所は、MapValidationErrors() に置換

例：if ve, ok := err.(validator.ValidationErrors) { writeJSON(MapValidationErrors(ve, locale, store, "domain.user")) }

DI/初期化

validator.New() は従来通り

翻訳登録（jatrans.RegisterDefaultTranslations）を呼ばない → バリデーションの日本語化は 辞書で 行う

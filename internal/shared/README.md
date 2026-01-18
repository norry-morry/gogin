# shared

## 用途ディレクトリの命名ガイド
| 用途        | 内容                              | 代表例                          |
|-----------|---------------------------------|------------------------------|
| util      | 軽い便利関数群（ポインタ、スライス、マップなど）        | pointer.go, slice.go         |
| conv      | 型変換関連（string⇔timeなど）            | timeconv.go, strconv.go      |
| fmt       | フォーマットやログ文字列整形                  | stringfmt.go, durationfmt.go |
| fn        | 汎用関数ライブラリ（エラーハンドリングやTry/Catch風） | try.go, recover.go           |
| validator | 共通バリデーション                       | email.go, enum.go            |

# 必ずお読み下さい

* [ローカル環境構築方法](https://rm.sstinc.jp/projects/resume_service/wiki/%E3%83%AD%E3%83%BC%E3%82%AB%E3%83%AB%E7%92%B0%E5%A2%83%E6%A7%8B%E7%AF%89)
* [Git運用ルール](https://rm.sstinc.jp/projects/open/wiki/Develop_guide_gitflow)
* [コーディング規約](https://rm.sstinc.jp/projects/resume_service/wiki/%E3%82%B3%E3%83%BC%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0%E8%A6%8F%E7%B4%84)
* [フロントエンドコーディング規約](https://rm.sstinc.jp/projects/resume_service/wiki/%E3%83%95%E3%83%AD%E3%83%B3%E3%83%88%E3%82%A8%E3%83%B3%E3%83%89%E8%A6%8F%E7%B4%84)
* [バックエンドコーディング規約](https://rm.sstinc.jp/projects/resume_service/wiki/%E3%83%90%E3%83%83%E3%82%AF%E3%82%A8%E3%83%B3%E3%83%89%E8%A6%8F%E7%B4%84)

## お時間のあるときにお読み下さい
* [Wiki運用ガイド](https://rm.sstinc.jp/projects/open/wiki/Guide)

# バックエンド

## ディレクトリツリーと役割
**要修正**
```shell
.
.
├── cmd/
│   └── api/
│       └── main.go               # エントリポイント（DI起動・HTTPサーバ起動）
├── configs/                      # 環境別設定テンプレ（yaml/env）
├── internal/
│   ├── app/
│   │   └── app.go                # App ランナー（Wireで組んだ依存を束ねて起動）
│   ├── config/
│   │   └── config.go             # 構成体 + 読み込み（viper/env）
│   ├── shared/                   # 横断ユーティリティ（純粋な依存逆転の対象）
│   │   ├── clock/clock.go        # 時刻抽象
│   │   ├── id/uuid.go            # ID生成抽象
│   │   ├── log/logger.go         # Logger IF（infra 実装は adapter/gateway 側）
│   │   └── errs/                 # ドメイン横断のエラー型・コード（i18nキーと結合）
│   ├── domain/                   # ドメイン層（純粋なビジネス）
│   │   ├── entity/
│   │   ├── valueobject/
│   │   ├── service/              # ドメインサービス
│   │   └── repository/           # Repository IF（ポート）
│   ├── usecase/                  # アプリケーション層（ユースケース/入力出力ポート）
│   │   ├── port/                 # 入力(Input)・出力(Output) Port IF->request,responseに分けた方が良いかも
│   │   └── interactor/           # 実装（ユースケースごと）
│   ├── adapter/
│   │   ├── http/                 # インターフェースアダプタ（REST）
│   │   │   ├── router/           # ルーティング定義（chi/ginなど）
│   │   │   ├── middleware/       # Auth, Recover, Locale, RequestID 等
│   │   │   ├── controller/       # ハンドラ（= 旧 handler）
│   │   │   ├── dto/              # Request/Response DTO（Validate対象）
│   │   │   └── presenter/        # 出力整形（エラー/ページング/リンク等）
│   │   ├── gateway/              # 永続化や外部I/Oの実装（= 旧 infrastructure）
│   │   │   ├── db/               # GORM 実装（repository IF を満たす）
│   │   │   ├── firebase/         # Firebase 実装
│   │   │   └── log/              # slog 実装（shared/log IF を満たす）
│   │   └── validation/           # 入力検証（go-playground/validator v10）
│   │       ├── rules/            # 再利用可能なルール（enum, 相関チェック等）
│   │       ├── translator/       # validator のメッセージ→i18n 連携
│   │       └── errors/           # Validation エラー→Presenter 用の変換
│   ├── i18n/                     # ☆ 多言語の中枢（ここが要件の肝）
│   │   ├── bundle.go             # バンドル作成（go-i18n 等でロード）
│   │   ├── loader.go             # /locales 読み込み、ホットリロード可
│   │   ├── locale.go             # 言語判定（Accept-Language, query, cookie）
│   │   ├── middleware.go         # Locale を Context に差し込む
│   │   └── locales/
│   │       ├── en.yaml
│   │       └── ja.yaml
│   ├── di/
│   │   └── wire.go               # Wire セット（http, gateway, usecase, i18n, config）
│   └── test/
│       ├── it/                   # 結合/IT（Testcontainers 等）
│       └── e2e/                  # API E2E（必要に応じて）
├── migrations/
│   └── *.sql
├── Makefile
├── go.mod
└── go.sum
```

## DIの注入
```shell
cd di && wire
```
## マイグレーション
```shell
# 新規作成
cd /go/src && go run cmd/migrate/main.go new {作成したいテーブル名}
# アップ(実行)
cd /go/src && go run cmd/migrate/main.go up
# ロールバック
cd /go/src && go run cmd/migrate/main.go down
# リセット
cd /go/src && go run cmd/migrate/main.go drop
# 現在のマイグレーション進捗状況確認
cd /go/src && go run cmd/migrate/main.go version
```
## テストデータ作成
```shell
# 新規作成
cd /go/src && go run cmd/migrate/main.go newseed {作成したいテーブル名}
# アップ(実行)
cd /go/src && go run cmd/migrate/main.go seedup
# 再実行
cd /go/src && go run cmd/migrate/main.go seedreset
# テストデータを空にする
cd /go/src && go run cmd/migrate/main.go seeddrop
```

## lintの実行
```shell
# lintの実行は
golangci-lint run --config .golangci.yml --out-format junit-xml ./... > ./artifacts/lint-report.xml
## で行えます。テスト結果はバックエンドのルートディレクトリに lint-report.xml が生成されます
```

# unittestの実行は
```shell
## UT、integrationテストも行う場合
gotestsum --junitfile ./artifacts/unit-test-report.xml -- -coverprofile=./artifacts/coverage.out ./...
## UTのみ行う場合
gotestsum --junitfile ./artifacts/unit-test-report.xml -- -coverprofile=./artifacts/coverage.out ./tests/unit/...
## featureのみ行う場合
gotestsum --junitfile ./artifacts/unit-test-report.xml -- -coverprofile=./artifacts/coverage.out ./tests/feature/...
```

## ビルドが通らなかったときのコマンド
```shell
go build -v ./... 2>&1 | tee /tmp/build.log
```

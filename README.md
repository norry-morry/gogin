# バックエンド

## ディレクトリツリーと役割

```shell
.
├── main.go                        # ✅ アプリエントリポイント
├── config/
│   ├── config.go                  # .env / yaml読み込み
│   ├── db.go                      # NewDB, DSN生成
│   ├── legger.go                  # slog + lumberjack ロガー設定
│   └── provider.go                #
├── di/
│   ├── wire.go                    # Wire定義（InitApp）
│   └── wire_gen.go                # 自動生成された依存解決コード
├── domain/
│   ├── user/
│   │   ├── entity.go              # Userエンティティ定義
│   │   └── repository.go          # UserRepositoryインターフェース定義
│   └── valueobject/
├── infrastructure/
│   ├── gorm/
│   │   └── repository.go          # GORMによるUserRepository実装
│   └── logger/
│       └── sql_logger.go          # slog連携したGORMロガー
├── interface/
│   └── handler/
│       └── user_handler.go        # GinのHandler（Controller）
├── middleware/                    # ミドルウェア
│       └── db.go
├── migration/
│       ├── migrations/            # マイグレーションファイル(up,down共に書く、autoマイグレーションは使わない)
│       └── main.go                # マイグレーションの実行コマンド
├── router/
│       └── router.go              # Ginエンジン・ルーティング
├── storage/                       # ログなどアプリケーション本体から生成されるけど、dockerで管理したくないディレクトリ(docker-compose.ymlでマウント)
│   └── logs/
├── usecase/                       # ユースケース
│   └── user/
│       ├── user_usecase.go        # ユースケース実装
│       └── user_usecase_test.go   # 手書きモックによるユニットテスト
├── Utility/                       # ユーティリティ
│   └── log.go
├── go.mod
└── go.sum
```

## DIの注入
```shell
cd di && wire
```
## マイグレーションの実行
```shell
cd migration && go run main.go up
```
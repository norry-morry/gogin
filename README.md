# バックエンド

## ディレクトリツリーと役割

```shell
.
├── main.go                    # ✅ エントリーポイントをここに
├── config/
│   ├── config.go              # 環境変数、設定 [x]
│   ├── db.go                  # 環境変数、設定 [x]
│   ├── legger.go              # 環境変数、設定 [x]
│   └── provider.go            # 環境変数、設定 [x]
├── di/
│   ├── wire.go                # Wireプロバイダ定義ファイル
│   └── wire_gen.go            # 自動生成されるDIコード
├── domain/                    # ドメイン（エンティティ・インタフェース）
│   ├── user/
│   │   ├── entity.go          # [x]
│   │   └── repository.go      # [x]
│   └── valueobject/
├── infrastructure/
│       └── db.go              # GORM接続 [x]
├── interface/                 # ハンドラとリポジトリ実装
│   ├── handler/
│   │   └── user_handler.go    # [x]
│   └── repository/
│       └── user_repository.go # [x]
├── middleware/
│       └── db.go              # GORM接続 [x]
├── migration/
│       ├── migrations/
│       └── main.go            # GORM接続 [x]
├── router/
│       └── router.go            # GORM接続 [x]
├── storage/                    # ドメイン（エンティティ・インタフェース）
│   └── logs/
├── usecase/                   # ユースケース
│   └── user/
│       └── user_usecase.go        # [x]
├── Utility/                   # ユースケース
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
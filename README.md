# バックエンド

## ディレクトリツリーと役割

```shell
.
├── main.go                    # ✅ エントリーポイントをここに
├── wire.go                    # Wireプロバイダ定義ファイル
├── wire_gen.go                # 自動生成されるDIコード
├── config/
│   └── config.go              # 環境変数、設定 [x]
├── domain/                    # ドメイン（エンティティ・インタフェース）
│   └── user/
│       ├── entity.go          # [x]
│       └── repository.go      # [x]
├── usecase/                   # ユースケース
│   └── user_usecase.go        # [x]
├── interface/                 # ハンドラとリポジトリ実装
│   ├── handler/
│   │   └── user_handler.go    # [x]
│   └── repository/
│       └── user_repository.go # [x]
├── infrastructure/
│       └── db.go              # GORM接続 [x]
├── go.mod
└── go.sum
```

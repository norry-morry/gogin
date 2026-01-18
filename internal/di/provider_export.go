// Package di は依存性注入 (Wire) のセット定義と Injector を提供します。
//
// Export* で始まるシンボルは、本来 internal/di パッケージ内に閉じている
// private な provider / factory 関数を、tests/unit/di など別パッケージから
// ユニットテストするためのエイリアスです。
package di

// ExportBackground は、アプリケーション全体で共有する background コンテキストを
// 生成する provider のテスト用エイリアスです。
var ExportBackground = background

// ExportNewAppLogger は、アプリケーションログ用 slog.Logger を生成する
// newAppLogger のテスト用エイリアスです。
var ExportNewAppLogger = newAppLogger

// ExportNewSQLLogger は、GORM の SQL ログ出力に使用する slog.Logger を生成する
// newSQLLogger のテスト用エイリアスです。
var ExportNewSQLLogger = newSQLLogger

// ExportProvideGormLogger は、GORM に渡す gorm.Logger 実装を生成する
// provideGormLogger のテスト用エイリアスです。
var ExportProvideGormLogger = provideGormLogger

// ExportProvideFBCreds は、Firebase Admin SDK で利用する認証情報を生成する
// provideFBCreds のテスト用エイリアスです。
var ExportProvideFBCreds = provideFBCreds

// ExportProvideFBOptions は、Firebase App の初期化に用いる設定オプションを生成する
// provideFBOptions のテスト用エイリアスです。
var ExportProvideFBOptions = provideFBOptions

// ExportProvideSharedTxRunner は、アプリケーション共通で利用する
// トランザクションランナーを生成する provideSharedTxRunner のテスト用エイリアスです。
var ExportProvideSharedTxRunner = provideSharedTxRunner

// ExportProvideAuthClock は、認証まわりで利用する Clock 実装を生成する
// provideAuthClock のテスト用エイリアスです。
var ExportProvideAuthClock = provideAuthClock

// ExportProvideDictionaryRepository は、辞書 (i18n) データを扱う Repository 実装を生成する
// provideDictionaryRepository のテスト用エイリアスです。
var ExportProvideDictionaryRepository = provideDictionaryRepository

// ExportProvideCacheStoreImpl は、翻訳メッセージ用の in-memory キャッシュストア実装を生成する
// provideCacheStoreImpl のテスト用エイリアスです。
var ExportProvideCacheStoreImpl = provideCacheStoreImpl

// ExportProvideTranslator は、i18n.Translator 実装を生成する
// provideTranslator のテスト用エイリアスです。
var ExportProvideTranslator = provideTranslator

// ExportProvideCacheStore は、i18n.CacheStore インターフェースに対する実装を提供する
// provideCacheStore のテスト用エイリアスです。
var ExportProvideCacheStore = provideCacheStore

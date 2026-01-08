//go:build wireinject
// +build wireinject

// Package di は依存性注入(Wire)のセット定義と Injector を提供します。
package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	gatewaydb "resume/internal/adapter/gateway/db"
	fb "resume/internal/adapter/gateway/firebase"
	gatewaymem "resume/internal/adapter/gateway/inmemory"
	gwlog "resume/internal/adapter/gateway/log"
	"resume/internal/adapter/http/controller"
	"resume/internal/adapter/http/router"
	validationadapter "resume/internal/adapter/validation"
	rules "resume/internal/adapter/validation/rules"
	"resume/internal/config"
	repo "resume/internal/domain/repository"
	vo "resume/internal/domain/valueobject/i18n"
	infradb "resume/internal/infra/db"
	stx "resume/internal/shared/tx"
	ucauth "resume/internal/usecase/auth"
	uceducation "resume/internal/usecase/education"
	uci18n "resume/internal/usecase/i18n"
	ucmeta "resume/internal/usecase/meta"
	ucprofile "resume/internal/usecase/profile"
	"resume/internal/usecase/user"
)

//
// ---- Helpers (wire が解決できるように包む) ----
//

// 背景コンテキスト
func background() context.Context { return context.Background() }

func newAppLogger(cfg config.Config) *slog.Logger {
	// gwlog.InitRotatingLogger に引数が要るならここで呼ぶ
	// 引数が不要なら return gwlog.InitRotatingLogger()
	return gwlog.InitRotatingLogger(cfg.AppLogPath)
}

// SQLLog は SQL ロガーを区別するためのマーカー型です。
// Wire の依存解決でログ種別を明確にするために使用します。
type SQLLog struct{ L *slog.Logger }

func newSQLLogger(cfg config.Config) *SQLLog {
	return &SQLLog{L: gwlog.InitRotatingLogger(cfg.SQLLogPath)}
}

func provideGormLogger(sql *SQLLog) *gatewaydb.GormLogger {
	return gatewaydb.NewGormLogger(sql.L, glogger.Info, 500*time.Millisecond, true)
}

// Firebase 用の Credentials/Options を Config から作る
func provideFBCreds(cfg config.Config) fb.Credentials {
	return fb.Credentials{
		CredsFile: cfg.Firebase.CredsFile,
		ProjectID: cfg.Firebase.ProjectID,
		// CredsJSON を使う場合はここで詰める
	}
}

func provideFBOptions(cfg config.Config) fb.Options {
	return fb.Options{
		UseAuthEmulator:  cfg.Firebase.UseAuthEmulator,
		AuthEmulatorHost: cfg.Firebase.AuthEmulatorHost,
	}
}

func provideSharedTxRunner(db *gorm.DB) stx.Runner {
	return infradb.NewTxRunner(db)
}

// func provideAuthTxRunner(db *gorm.DB) ucauth.TxRunner           { return infradb.NewTxRunner(db) }
func provideAuthClock() ucauth.Clock { return nil }

// --- i18n providers（wire.go内の非export） ---
func provideDictionaryRepository(root string) repo.DictionaryRepository {
	return gatewaydb.NewI18nRepositoryImpl(root)
}
func provideCacheStoreImpl() *gatewaymem.CacheStoreImpl { // infraの具象
	return gatewaymem.NewCacheStoreImpl()
}
func provideTranslator(cs *gatewaymem.CacheStoreImpl) uci18n.Translator { return cs }
func provideCacheStore(cs *gatewaymem.CacheStoreImpl) uci18n.CacheStore { return cs }

// Usecase生成＋Reloadを1本化（重複プロバイダ回避）
func provideI18nUsecase(
	r repo.DictionaryRepository,
	tr uci18n.Translator,
	cs uci18n.CacheStore,
) (uci18n.Usecase, error) {
	// ① デフォルトロケールを決める（例）
	//    実際の Config 構造に合わせて修正してね。
	//    もしなければ一旦 "ja" ベタ書きでも可。
	//   locale := vo.NewLocale(cfg.I18n.DefaultLocale)
	//   or locale := vo.MustParseLocale(cfg.I18n.DefaultLocale)
	//   みたいなイメージ。
	locale := vo.NewLocale("ja")

	// ② Translator を使った MessageKeyStore を作って rules に登録
	keyStore := validationadapter.NewI18nMessageKeyStore(tr, locale)
	rules.SetMessageKeyStore(keyStore)

	// ③ I18n Usecase 本体の初期化 + Reload
	uc := uci18n.New(tr, r, cs)
	if err := uc.Reload(context.Background()); err != nil {
		return nil, err
	}
	return uc, nil
}

//
// ---- Sets ----
//

// ConfigSet は アプリ設定を提供する Wire セット。
var ConfigSet = wire.NewSet(config.Load)

// AppLoggerSet は アプリログを提供する Wire セット。
var AppLoggerSet = wire.NewSet(
	newAppLogger, // ← 直接 gwlog.InitRotatingLogger を渡さず、引数解決込みのラッパを渡す
)

// SQLLoggerSet は クエリログを提供する Wire セット。
var SQLLoggerSet = wire.NewSet(
	newSQLLogger,
)

// GormLoggerSet は クエリログを提供する Wire セット。
var GormLoggerSet = wire.NewSet(
	SQLLoggerSet,
	wire.Bind(new(glogger.Interface), new(*gatewaydb.GormLogger)),
	// GORM ロガーは SQLLog だけを受け取るので曖昧性無し
	provideGormLogger,
)

// DBSet は MySQLのdsn情報を提供する Wire セット。
var DBSet = wire.NewSet(infradb.NewGorm)

// GatewaySet は　Repositoryを提供するセット
var GatewaySet = wire.NewSet(
	gatewaydb.NewUserRepository, // func(*gorm.db) repository.UserRepository
	gatewaydb.NewIdentityRepository,
	gatewaydb.NewUserAddressRepository,
	gatewaydb.NewAddressPurposeRepository,
	gatewaydb.NewCountryRepository,
	gatewaydb.NewGenderRepository,
	provideDictionaryRepository,
	provideCacheStoreImpl,
	gatewaydb.NewUserProfileRepository,
	gatewaydb.NewEducationStatusRepository,
	gatewaydb.NewDegreeTypeRepository,
	gatewaydb.NewUserEducationRepository,
)

// UsecaseSet は　usecaseを提供するセット
var UsecaseSet = wire.NewSet(
	provideSharedTxRunner,
	user.NewUsecase, // func(repository.UserRepository, *slog.Logger) user.Usecase
	ucauth.New,
	provideAuthClock,
	ucprofile.New,
	provideTranslator,
	provideCacheStore,
	provideI18nUsecase,
	ucmeta.New,
	uceducation.New,
)

// ControllerSet は ハンドラを提供するセット
var ControllerSet = wire.NewSet(
	controller.NewSampleHandler,
	controller.NewTestHandler,
	controller.NewUserHandler,
	controller.NewAuthHandler,
	controller.NewProfileHandler,
	controller.NewI18nHandler,
	controller.NewMetaHandler,
	controller.NewEducationHandler,
)

// FirebaseSet は Firebaseの認証情報を提供するセット
var FirebaseSet = wire.NewSet(
	background,
	provideFBCreds,
	provideFBOptions,
	fb.NewAuthClient, // func(ctx, creds, opt) (*auth.Client, error)
	fb.NewAuth,       // func(*auth.Client) fb.Auth
)

// HTTPSet は Gin ルータを構成する Wire セット。
var HTTPSet = wire.NewSet(
	router.New, // func(fb fb.Auth) *gin.Engine
)

//
// ---- Injector ----
//

// InitAPI は、Wire によって依存関係を解決し、アプリケーション全体の API エンジンを初期化して返します。
// 戻り値の *gin.Engine は、ルーティングやミドルウェアを含む HTTP サーバーのエントリポイントです。
func InitAPI(i18nRoot string) (*gin.Engine, error) {
	wire.Build(
		ConfigSet,
		AppLoggerSet,
		GormLoggerSet,
		DBSet,
		GatewaySet,
		UsecaseSet,
		FirebaseSet,
		ControllerSet,
		HTTPSet,
	)
	return nil, nil
}

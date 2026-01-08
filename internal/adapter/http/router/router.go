// Package router はアプリケーションのHTTPルーティングを定義します。
package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	fba "resume/internal/adapter/gateway/firebase"
	"resume/internal/adapter/http/controller"
	"resume/internal/adapter/http/middleware"
	"resume/internal/domain/repository"
	uci18n "resume/internal/usecase/i18n"
)

// New はアプリケーションのルーティングを初期化し、*gin.Engine を返します。
func New(
	fb fba.Auth,
	hSample *controller.SampleHandler,
	hTest *controller.TestHandler,
	hUser *controller.UserHandler,
	hAuth *controller.AuthHandler,
	hProfile *controller.ProfileHandler,
	hMeta *controller.MetaHandler,
	hI18n *controller.I18nHandler,
	appLog *slog.Logger,
	userRepo repository.UserRepository,
	ucI18n uci18n.Usecase,
	hEducation *controller.EducationHandler,
) *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = false // 末尾スラッシュ補正をやめる
	r.RedirectFixedPath = false     // 大文字小文字などの補正もやめる
	r.Use(middleware.CamelSnakeCodec())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors(middleware.DefaultCORSConfig()))
	r.Use(middleware.InjectAuth(fb))

	r.Use(middleware.RequestLogger(appLog, &middleware.RequestLoggerOptions{
		Skipper:        func(c *gin.Context) bool { return c.Request.URL.Path == "/health" },
		IncludeHeaders: []string{"X-Request-ID", "X-Trace-Id"},
	}))

	// region 全ルーター対策
	// 404が返らなかった対策
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
	})
	// メソッド違いに対応出来てなかった対策
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "method not allowed",
		})
	})
	// cors対策 任意：OPTIONS を一括ハンドル（プリフライトを確実に 204 で返す）
	//r.OPTIONS("/*path", func(c *gin.Context) {
	//	c.Status(204)
	//})
	// endregion 全ルーター対策

	// region i18n

	// (b) ロケール交渉ミドルウェア（?lang, cookie=lang, Accept-Language）
	src := NewUsecaseLocaleSource(ucI18n)
	r.Use(middleware.Middleware(middleware.Opts{
		Default:   "ja",
		Allow:     nil, // 増やしたらここに追記
		Query:     "lang",
		Header:    "Accept-Language",
		Cookie:    "", // Cookieを使わないなら空
		SetHeader: true,
		Store:     src,
	}))

	i18nG := r.Group("/i18n")
	{
		// (c) バンドル配信用ハンドラ（ETag/304 対応）
		i18nG.GET("/bundle", hI18n.GetBundle)
		i18nG.POST("/reload", hI18n.Reload)
		i18nG.GET("/list", hI18n.ListLocales)
	}

	// endregion i18n

	// region テスト
	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// テスト
	r.GET("ping", hSample.Ping)
	testG := r.Group("test")
	testG.Use(middleware.RequireAuth(fb))
	{
		testG.GET("1", hTest.Ping)
	}
	// endregion テスト

	// region 認証
	authG := r.Group("/auth")
	{
		authG.POST("/fl/login", hAuth.UserLogin)
	}
	// endregion 認証

	// region firebase認可
	resolveUserMw := middleware.NewResolveUser(userRepo)
	// フリーランサー情報
	freelanceG := r.Group("/fl")
	freelanceG.Use(middleware.RequireAuth(fb))
	freelanceG.Use(resolveUserMw.RequireResolveUser())
	{
		// ログインしている自身の情報
		profileG := freelanceG.Group("/profile")
		profileG.GET("/", hProfile.Me)
		profileG.GET("/fba", hProfile.Fba)
		profileG.GET("/identities", hProfile.Identities)

		// 履歴書掲載情報の氏名など
		personalG := freelanceG.Group("/personal")
		personalG.GET("/", hProfile.GetPersonalInfo)
		personalG.PATCH("/", hProfile.PatchPersonalInfo)
		personalG.GET("/exists", hProfile.HasPersonalInfo)
		//PATCH  /profile/account             # メールやパスワードなどアカウント設定（任意）

		addressG := freelanceG.Group("/address")
		addressG.GET("/exists", hProfile.HasUserAddress)
		addressG.GET("/", hProfile.ListUserAddress)
		addressG.GET("/:address_id", hProfile.DetailUserAddress)
		addressG.POST("/", hProfile.CreateUserAddress)
		addressG.PATCH("/:address_id", hProfile.UpdateUserAddress)
		addressG.DELETE("/:address_id", hProfile.DeleteUserAddress)

		// 学歴関連情報
		educationG := freelanceG.Group("/education")
		educationG.GET("/", hEducation.ListEducation)
		educationG.GET("/exists", hEducation.ExistsEducation)
		educationG.POST("/", hEducation.AddEducation)
		educationG.PATCH("/:education_id", hEducation.UpdateEducation)
		educationG.DELETE("/:education_id", hEducation.DeleteEducation)
		educationG.PATCH("/order", hEducation.OrderEducation)
	}
	// endregion firebase認可

	// region agent、corporate、adminが使うユーザー情報
	userGroup := r.Group("/users")
	{
		userGroup.GET("/:id", hUser.GetByID)
	}
	// endregion agent、corporate、adminが使うユーザー情報

	// region 認証とかロールとか無関係に誰でも使うやつを置く

	// selectとかcheckboxとかradioのオプションを取得するAPIを格納する
	optionG := r.Group("/option")
	{
		optionAddressG := optionG.Group("/address")
		{
			// 目的
			optionAddressG.GET("/purpose", hMeta.ListAddressPurpose)
			// 国コード
			optionAddressG.GET("/countries", hMeta.ListCountry)
		}
		// 性別
		optionG.GET("/gender", hMeta.ListGender)

		optionEducationG := optionG.Group("/education")
		{
			// 学歴状態
			optionEducationG.GET("/status", hMeta.ListEducationStatus)
			// 学位種別
			optionEducationG.GET("/degree", hMeta.ListDegreeType)
		}
	}

	// endregion 認証とかロールとか無関係に誰でも使うやつを置く

	return r
}

//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"log/slog"
	"resume/config"
	infra "resume/infrastructure/gorm"
	"resume/interface/handler"
	"resume/router"
	"resume/usecase/user"
)

func InitApp(cfg config.Config, sqlLogger *slog.Logger) (*gin.Engine, error) {
	wire.Build(
		config.ProvideMySQLSettings,
		config.NewDB,

		// User関連
		infra.NewUserRepository,
		user.NewUserUsecase,
		handler.NewUserHandler,
		// note コントローラーが増えたら下記に追加して、wireを実行すると動的にDIの注入は行ってくれます
		//		ルーターの引数は、手書きで適宜追加しないといけないっぽいです

		router.SetupRouter,
	)
	return &gin.Engine{}, nil
}

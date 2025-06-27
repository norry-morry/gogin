package Routes

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"resume/Handlers"
	"resume/middlewares"
)

func SetupRouter(logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.RequestLoggerMiddleware(logger))

	group := router.Group("/api")
	{
		group.GET("/ping", Handlers.GetPing)
		group.GET("/sample", Handlers.GetSample)
		test := group.Group("/test")
		{
			test.GET("1", Handlers.GetTest)
			test.GET("2", Handlers.GetTest2)
		}
	}
	return router
}

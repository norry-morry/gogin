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
		//user := group.Group("/user")
		//{
		//	user.GET("/", Handlers.GetUsers)
		//	user.GET("/:id", Handlers.GetUserById)
		//	user.POST("/", Handlers.CreateUser)
		//	user.PUT("/:id", Handlers.UpdateUser)
		//	user.DELETE("/:id", Handlers.DeleteUser)
		//}
	}
	return router
}

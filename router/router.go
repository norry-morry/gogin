package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"resume/interface/handler"
	"resume/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		slog.Debug("PingHandler called") // default logger 使用
		c.JSON(200, gin.H{
			"message": "Hello air! pong2",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/ping", handler.GetPing)
		api.GET("/sample", handler.GetSample)
		test := api.Group("/test")
		{
			test.GET("1", handler.GetTest)
			test.GET("2", handler.GetTest2)
		}
		user := api.Group("/users")
		{
			user.GET("/", userHandler.GetAllUsers)
			user.GET("/:id", userHandler.GetUserById)
			user.POST("/", userHandler.CreateUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}
	}
	return r
}

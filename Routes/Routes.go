package Routes

import (
	"github.com/gin-gonic/gin"
	"resume/Controllers"
	"resume/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	group := r.Group("/api")
	{
		group.GET("/ping", Controllers.GetPing)
		group.GET("/sample", Controllers.GetSample)
		test := group.Group("/test")
		{
			test.GET("1", handler.GetTest)
			test.GET("2", handler.GetTest2)
		}
		user := group.Group("/user")
		{
			user.GET("/", handler.GetUsers)
			user.GET("/:id", handler.GetUserById)
			user.POST("/", handler.CreateUser)
			user.PUT("/:id", handler.UpdateUser)
			user.DELETE("/:id", handler.DeleteUser)
		}
	}
	return r
}

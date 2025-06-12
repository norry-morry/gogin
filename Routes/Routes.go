package Routes

import (
	"github.com/gin-gonic/gin"
	"resume/Controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	grp1 := r.Group("/api")
	{
		grp1.GET("/ping", Controllers.GetPing)
		grp1.GET("/sample", Controllers.GetSample)
	}
	return r
}

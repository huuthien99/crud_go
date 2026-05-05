package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//  prefix
	api := r.Group("/api/v1")

	// Đăng ký các routes theo module
	SetupAuthRoutes(api)
	SetupTaskRoutes(api)

	return r
}

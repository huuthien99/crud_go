package routes

import (
	"auth_crud/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//  prefix
	api := r.Group("/api/v1")

	// auth routes
	SetupAuthRoutes(api)

	// middleware
	api.Use(middleware.AuthMiddleware())

	// task
	SetupTaskRoutes(api)

	return r
}

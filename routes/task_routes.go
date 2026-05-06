package routes

import (
	"auth_crud/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTaskRoutes(router *gin.RouterGroup) {
	protected := router.Group("")

	task := protected.Group("/tasks")
	{
		task.POST("/", controllers.CreateTask)
		task.GET("/", controllers.GetTaskList)
		task.GET("/:id", controllers.GetTaskById)
		task.DELETE("/:id", controllers.DeleteTask)
		task.PATCH("/:id", controllers.UpdateTask)
		task.PATCH("/change-status/:id", controllers.ChangeStatus)
	}
}

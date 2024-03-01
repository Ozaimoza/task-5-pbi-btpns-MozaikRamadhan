package routes

import (
	"task-5-pbi-btpns/controllers"
	"task-5-pbi-btpns/middleware"

	"github.com/gin-gonic/gin"
)

func PhotoRouter(r *gin.Engine) {
	userGroup := r.Group("/photo")
	{
		userGroup.GET("/", controllers.GetAllPhoto)
		userGroup.POST("/", middleware.AuthMiddleware, controllers.PostPhoto)
		userGroup.PUT("/:PhotoId", middleware.AuthMiddleware, controllers.UpdatePhoto)
		userGroup.DELETE("/:PhotoId", middleware.AuthMiddleware, controllers.DeletePhoto)

	}
}

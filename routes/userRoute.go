package routes

import (
	"task-5-pbi-btpns/controllers"
	"task-5-pbi-btpns/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", controllers.RegisterUser)
		userGroup.POST("/login", controllers.LoginUser)
		userGroup.POST("/logout", middleware.AuthMiddleware, controllers.LogoutUser)
		// userGroup.GET("/validate", middleware.AuthMiddleware, controllers.Validate)
		userGroup.PUT("/:userId", middleware.AuthMiddleware, controllers.UpdateUser)
		userGroup.DELETE("/:userId", middleware.AuthMiddleware, controllers.DeleteUser)
	}
}

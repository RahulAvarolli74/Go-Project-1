package routes

import (
	"recipe-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", controllers.RegisterUser)
		users.GET("/:id", controllers.GetUserByID)
	}
}

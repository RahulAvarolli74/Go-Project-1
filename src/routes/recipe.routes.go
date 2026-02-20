package routes

import (
	"recipe-api/src/controllers"
	"recipe-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRecipeRoutes(rg *gin.RouterGroup) {
	recipes := rg.Group("/recipes")
	{
		recipes.GET("/search", controllers.SearchByIngredients)
		recipes.GET("", controllers.GetAllRecipes)
		recipes.GET("/:id", controllers.GetRecipeByID)
		recipes.POST("", middlewares.UploadImage(), controllers.CreateRecipe)
		recipes.PUT("/:id", controllers.UpdateRecipe)
		recipes.DELETE("/:id", controllers.DeleteRecipe)
	}
}

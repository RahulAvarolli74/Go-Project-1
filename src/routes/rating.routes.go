package routes

import (
	"recipe-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRatingRoutes(rg *gin.RouterGroup) {
	ratings := rg.Group("/recipes")
	{
		ratings.POST("/:id/ratings", controllers.AddRating)
		ratings.GET("/:id/ratings", controllers.GetRatings)
	}
}

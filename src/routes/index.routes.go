package routes

import (
	"net/http"

	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/api/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, "🚀 Recipe API is running!", gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	api := router.Group("/api")

	RegisterRecipeRoutes(api)
	RegisterRatingRoutes(api)
	RegisterUserRoutes(api)

	router.NoRoute(func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusNotFound,
			"Route not found. Check the API documentation at /api/health")
	})
}

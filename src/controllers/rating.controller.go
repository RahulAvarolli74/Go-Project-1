package controllers

import (
	"net/http"

	"recipe-api/src/db"
	"recipe-api/src/models"
	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
)

func AddRating(c *gin.Context) {
	recipeID := c.Param("id")

	var recipe models.Recipe
	if err := db.DB.First(&recipe, "id = ?", recipeID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Recipe not found")
		return
	}

	var rating models.Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			"Invalid rating data. Score must be 1-5 and user_name is required: "+err.Error())
		return
	}

	rating.RecipeID = recipeID

	result := db.DB.Create(&rating)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to add rating: "+result.Error.Error())
		return
	}

	updateAverageRating(recipeID)
	db.DB.First(&recipe, "id = ?", recipeID)

	utils.SuccessResponse(c, http.StatusCreated, "Rating added successfully! ⭐", gin.H{
		"rating": rating,
		"recipe": recipe,
	})
}

func GetRatings(c *gin.Context) {
	recipeID := c.Param("id")

	var recipe models.Recipe
	if err := db.DB.First(&recipe, "id = ?", recipeID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Recipe not found")
		return
	}

	var ratings []models.Rating
	result := db.DB.Where("recipe_id = ?", recipeID).
		Order("created_at DESC").
		Find(&ratings)

	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to fetch ratings: "+result.Error.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Ratings fetched successfully", gin.H{
		"recipe_id":      recipeID,
		"average_rating": recipe.AverageRating,
		"count":          len(ratings),
		"ratings":        ratings,
	})
}

func updateAverageRating(recipeID string) {
	var avgResult struct {
		Average float64
	}

	db.DB.Model(&models.Rating{}).
		Select("COALESCE(AVG(score), 0) as average").
		Where("recipe_id = ?", recipeID).
		Scan(&avgResult)

	db.DB.Model(&models.Recipe{}).
		Where("id = ?", recipeID).
		Update("average_rating", avgResult.Average)
}

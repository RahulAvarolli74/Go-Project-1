package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"recipe-api/src/db"
	"recipe-api/src/models"
	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
)

func CreateRecipe(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	ingredients := c.PostForm("ingredients")
	userID := c.PostForm("user_id")
	prepTime, _ := strconv.Atoi(c.DefaultPostForm("prep_time", "0"))
	cookTime, _ := strconv.Atoi(c.DefaultPostForm("cook_time", "0"))
	servings, _ := strconv.Atoi(c.DefaultPostForm("servings", "1"))

	if title == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Title is required")
		return
	}
	if ingredients == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ingredients are required (JSON array string)")
		return
	}

	var ingredientList []string
	if err := json.Unmarshal([]byte(ingredients), &ingredientList); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			"Ingredients must be a valid JSON array. Example: [\"tomato\",\"onion\"]")
		return
	}

	var imageURL string
	if filePath, exists := c.Get("uploadedFilePath"); exists {
		processedFilename, err := utils.ProcessImage(filePath.(string))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError,
				"Failed to process uploaded image: "+err.Error())
			return
		}
		imageURL = "/uploads/" + processedFilename
	}

	recipe := models.Recipe{
		Title:       title,
		Description: description,
		Ingredients: ingredients,
		ImageURL:    imageURL,
		PrepTime:    prepTime,
		CookTime:    cookTime,
		Servings:    servings,
		UserID:      userID,
	}

	result := db.DB.Create(&recipe)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to create recipe: "+result.Error.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Recipe created successfully! 🎉", recipe)
}

func GetAllRecipes(c *gin.Context) {
	var recipes []models.Recipe
	var totalCount int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	db.DB.Model(&models.Recipe{}).Count(&totalCount)

	result := db.DB.Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&recipes)

	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to fetch recipes: "+result.Error.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK,
		"Recipes fetched successfully", recipes, page, perPage, totalCount)
}

func GetRecipeByID(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	result := db.DB.Preload("Ratings").First(&recipe, "id = ?", id)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Recipe not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recipe fetched successfully", recipe)
}

func SearchByIngredients(c *gin.Context) {
	ingredientsParam := c.Query("ingredients")
	if ingredientsParam == "" {
		utils.ErrorResponse(c, http.StatusBadRequest,
			"Please provide ingredients to search. Example: ?ingredients=tomato,onion")
		return
	}

	searchTerms := strings.Split(ingredientsParam, ",")

	var recipes []models.Recipe
	query := db.DB

	for i, term := range searchTerms {
		term = strings.TrimSpace(strings.ToLower(term))
		if term == "" {
			continue
		}
		if i == 0 {
			query = query.Where("LOWER(ingredients) LIKE ?", "%"+term+"%")
		} else {
			query = query.Or("LOWER(ingredients) LIKE ?", "%"+term+"%")
		}
	}

	result := query.Find(&recipes)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Search failed: "+result.Error.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK,
		"Search results fetched successfully", gin.H{
			"query":   ingredientsParam,
			"count":   len(recipes),
			"recipes": recipes,
		})
}

func UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := db.DB.First(&recipe, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Recipe not found")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	delete(updateData, "id")
	delete(updateData, "created_at")
	delete(updateData, "average_rating")

	result := db.DB.Model(&recipe).Updates(updateData)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to update recipe: "+result.Error.Error())
		return
	}

	db.DB.First(&recipe, "id = ?", id)
	utils.SuccessResponse(c, http.StatusOK, "Recipe updated successfully", recipe)
}

func DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := db.DB.First(&recipe, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Recipe not found")
		return
	}

	db.DB.Where("recipe_id = ?", id).Delete(&models.Rating{})

	result := db.DB.Delete(&recipe)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to delete recipe: "+result.Error.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Recipe deleted successfully", nil)
}

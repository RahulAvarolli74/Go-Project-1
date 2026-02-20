package controllers

import (
	"net/http"

	"recipe-api/src/db"
	"recipe-api/src/models"
	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			"Invalid user data. Username and email are required: "+err.Error())
		return
	}

	var existing models.User
	if err := db.DB.Where("username = ? OR email = ?", user.Username, user.Email).
		First(&existing).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username or email already exists")
		return
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			"Failed to register user: "+result.Error.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully! 👤", user)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	result := db.DB.Preload("Recipes").First(&user, "id = ?", id)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User fetched successfully", user)
}

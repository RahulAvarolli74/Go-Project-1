package main

import (
	"log"
	"os"

	"recipe-api/src/db"
	"recipe-api/src/middlewares"
	"recipe-api/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	db.ConnectDatabase()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.Use(middlewares.ErrorHandler())

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./public/temp"
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("⚠️  Could not create upload directory: %v", err)
	}
	router.Static("/uploads", uploadDir)

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("==============================================")
	log.Println("  🍳 Recipe Sharing API")
	log.Println("  📍 Running on: http://localhost:" + port)
	log.Println("  📦 Database:   SQLite (./recipe.db)")
	log.Println("  📁 Uploads:    " + uploadDir)
	log.Println("==============================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

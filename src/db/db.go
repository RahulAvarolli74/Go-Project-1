package db

import (
	"log"
	"os"

	"recipe-api/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./recipe.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully (SQLite)")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Rating{},
	)
	if err != nil {
		log.Fatalf("❌ Auto-migration failed: %v", err)
	}

	log.Println("✅ Database tables migrated successfully")
}

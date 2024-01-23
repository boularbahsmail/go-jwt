package initializers

import "go-jwt/models"

func SyncDatabase() {
	// Migrate the User schema
	DB.AutoMigrate(&models.User{})
}

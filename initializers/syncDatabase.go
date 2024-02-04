package initializers

import "github.com/jwt-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}

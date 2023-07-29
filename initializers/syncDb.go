package initializers

import "backend-technoscape/models"

func SyncDb() {
	DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		// &models.Order{},
	)

}

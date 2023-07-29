package main

import (
	"backend-technoscape/models"
	"fmt"
	"log"

	"github.com/wpcodevo/golang-gorm-postgres/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("================= Could not load environment variables ", err, " ================¸")
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)

	fmt.Println("================¸ Migration complete ================¸")
}

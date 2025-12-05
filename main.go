package main

import (
	"core/config"
	"core/models"
	"core/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.DatabaseSetup()
	migrate := config.GetDB().AutoMigrate(&models.User{},
		&models.FormToken{},
		&models.AllowedDomain{})
	if migrate != nil {
		log.Fatalf("Error on migrations")
	}

	router := routes.SetupRoutes()
	if err := router.Run(":8030"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

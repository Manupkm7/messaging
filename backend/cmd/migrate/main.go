package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error loading config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer database.Close()

	if err := database.RunMigrations(); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations completed successfully!")
}

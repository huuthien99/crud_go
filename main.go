package main

import (
	"auth_crud/config"
	"auth_crud/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	// lấy PORT từ env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	config.ConnectDatabase()
	r := routes.SetupRouter()

	r.Run(":" + port)
}

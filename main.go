package main

import (
	"log"
	"os"

	"go-rest-api/app"
	"go-rest-api/models"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB, err := models.NewDB(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)
	if err != nil {
		log.Panic(err)
	}

	a := &app.App{
		User: DB,
		Auth: DB,
	}

	a.Routes()
}

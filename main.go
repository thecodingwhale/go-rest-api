package main

import (
  "os"
  "log"
  "github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  app := App{}
  app.database(
    os.Getenv("APP_DB_DRIVER"),
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"),
  )
  app.routes()
  app.run(":8080")
}

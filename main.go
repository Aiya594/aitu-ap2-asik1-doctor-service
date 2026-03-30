package main

import (
	"log"
	"os"

	"github.com/Aiya594/doctor-service/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	port := os.Getenv("PORT")
	app := app.NewApp()
	app.RunServer(port)
}

package main

import (
	"golang_template/db"
	"golang_template/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	log.Println("Starting ...")

	// initial route
	r := routes.InitialRoutes()

	// Database connect
	if err := db.ConnectMongoDb(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Start the server
	log.Println("Started server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

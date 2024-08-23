package main

import (
	"log"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/routes"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	database.ConnectDB()
}

func main() {
	// get underlying sql connection
	SQLdb, err := database.DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	SQLdb.SetMaxOpenConns(25)
	SQLdb.SetMaxIdleConns(25)
	SQLdb.SetConnMaxLifetime(5 * time.Minute)

	defer SQLdb.Close()

	r := routes.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"log"

	"github.com/Alan-Luc/clog/backend/database"
	"github.com/Alan-Luc/clog/backend/routes"
	"github.com/Alan-Luc/clog/backend/utils/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	// connect db
	database.ConnectDB()
}

func main() {
	defer database.CloseDB() // defer db close
	// router setup
	srv := routes.StartServer()

	server.GracefulShutdown(srv)
}

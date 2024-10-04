package main

import (
	"log"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/routes"
	"github.com/Alan-Luc/VertiLog/backend/utils/logger"
	"github.com/Alan-Luc/VertiLog/backend/utils/server"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	// connect db
	database.ConnectDB()
	// initialize logging
	logger.InitLogger()
}

func main() {
	defer database.CloseDB()   // defer db close
	defer logger.Logger.Sync() // flush logger
	// router setup
	srv := routes.StartServer()

	server.GracefulShutdown(srv, logger.Logger)
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	// get underlying sql connection and defer close
	defer func() {
		if err := database.SQLdb.Close(); err != nil {
			log.Fatal("Error closing database connection... ", err)
		} else {
			log.Println("Closing database connection...")
		}
	}()

	// router setup
	r := routes.SetupRouter()
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Stop channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// block until signal
	<-sc
	log.Println("Gracefully shutting down...")

}

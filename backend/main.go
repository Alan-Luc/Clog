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
	// defer close
	defer database.CloseDB()

	// router setup
	go routes.SetupRouter()

	// Stop channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// block until signal
	<-sc
	log.Println("Gracefully shutting down...")

}

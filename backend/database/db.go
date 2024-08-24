package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SQLdb *sql.DB

func ConnectDB() {
	var err error
	connStr := os.Getenv("DB_URI")

	// open GORM connection
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	SQLdb, err = DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	SQLdb.SetMaxOpenConns(25)
	SQLdb.SetMaxIdleConns(25)
	SQLdb.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Successfully connected to the database!")
}

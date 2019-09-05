package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	// For postgres connection
	_ "github.com/lib/pq"
)

// DB exposes a global variable referencing the database connection
var DB *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
	}
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbuser, dbpassword, dbname)
	OpenDB(dbinfo)
}

// OpenDB initializes the database connection
func OpenDB(dataSourceName string) {
	DB, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	// For postgres connection
	_ "github.com/lib/pq"
)

// Conn exposes a global variable referencing the database connection
var Conn *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbname := os.Getenv("DB_NAME")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbuser, dbpassword, dbname)
	Conn = OpenDB(dbinfo)
}

// OpenDB initializes the database connection
func OpenDB(dataSourceName string) *sql.DB {
	conn, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Panic(err)
	}

	if err = conn.Ping(); err != nil {
		log.Panic(err)
	}

	return conn
}

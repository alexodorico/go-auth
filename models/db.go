package models

import (
	"database/sql"
	"log"

	// For postgres connection
	_ "github.com/lib/pq"
)

// DB is the global variable referencing the database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

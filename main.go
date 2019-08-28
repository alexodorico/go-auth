package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname     = os.Getenv("DB_NAME")
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var lastInsertID int
	err = db.QueryRow("INSERT INTO users(username,password) VALUES($1,$2) RETURNING id", "oxideorcoal", "wow").Scan(&lastInsertID)
	checkErr(err)
	fmt.Println("last inserted ID =", lastInsertID)

	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {

	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

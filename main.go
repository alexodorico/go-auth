package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname     = os.Getenv("DB_NAME")
)

func main() {
	http.HandleFunc("/api/user", handleUser)
	err := http.ListenAndServe(":9000", nil)
	checkErr(err)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUser(w, r)
	case "POST":
		postUser(w, r)
	default:
		fmt.Printf("Only GET and POST methods are allowed on this endpoint")
	}

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

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET on /api/user")
}

func postUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("POST on /api/user")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

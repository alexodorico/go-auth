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
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", dbinfo)

	switch r.Method {
	case "GET":
		getUser(w, r, db)
	case "POST":
		postUser(w, r, db)
	default:
		fmt.Printf("Only GET and POST methods are allowed on this endpoint")
	}

	checkErr(err)
	defer db.Close()
}

func getUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Printf("GET on /api/user")
	rows, err := db.Query("SELECT id, username, password FROM users")
	checkErr(err)

	for rows.Next() {
		var id int
		var username string
		var password string
		err = rows.Scan(&id, &username, &password)
		checkErr(err)
		fmt.Printf("%3v | %s | %s ", id, username, password)
	}
}

func postUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Printf("POST on /api/user")
	var lastInsertID int
	err := db.QueryRow("INSERT INTO users(username,password) VALUES($1,$2) RETURNING id", "oxideorcoal", "wow").Scan(&lastInsertID)
	checkErr(err)
	fmt.Println("last inserted ID =", lastInsertID)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

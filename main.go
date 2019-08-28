package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname     = os.Getenv("DB_NAME")
)

type user struct {
	UID      int    `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/api/user/", handleUser)
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
	var u user
	params := strings.Split(r.URL.Path, "/")
	uid := params[3]
	rows, err := db.Query("SELECT id, username, password FROM users WHERE id = $1", uid)
	checkErr(err)

	for rows.Next() {
		var id int
		var uname string
		var pw string
		err = rows.Scan(&id, &uname, &pw)
		u = user{UID: id, Username: uname, Password: pw}
		checkErr(err)
	}

	j, err := json.Marshal(u)
	checkErr(err)
	fmt.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
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

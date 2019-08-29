package main

import (
	"database/sql"
	"encoding/json"
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

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	err := http.ListenAndServe(":9000", nil)
	checkErr(err)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var sStmt = "INSERT INTO users(username,password,email) VALUES($1,$2,$3)"
	var u user
	var res response

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	checkErr(err)

	stmt, err := db.Prepare(sStmt)
	checkErr(err)

	_, err = stmt.Exec(u.Username, u.Password, u.Email)
	checkErr(err)

	res = response{Message: "successful registration", Success: true}
	j, err := json.Marshal(res)
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	stmt.Close()
}

// func getUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	var u user
// 	params := strings.Split(r.URL.Path, "/")
// 	uid := params[3]
// 	rows, err := db.Query("SELECT username, password FROM users WHERE id = $1", uid)
// 	checkErr(err)

// 	for rows.Next() {
// 		var uname string
// 		var pw string
// 		err = rows.Scan(&uname, &pw)
// 		u = user{Username: uname, Password: pw}
// 		checkErr(err)
// 	}

// 	j, err := json.Marshal(u)
// 	checkErr(err)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(j)
// }

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

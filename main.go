package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/alexodorico/goserver/models"

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

// registration response
type regRes struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbuser, dbpassword, dbname)
	models.InitDB(dbinfo)

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	err := http.ListenAndServe(":9000", nil)
	checkErr(err)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var sStmt = "INSERT INTO users(username,password,email) VALUES($1,$2,$3) RETURNING id"
	var u user
	var res regRes
	var userID int

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	checkErr(err)

	stmt, err := models.DB.Prepare(sStmt)
	checkErr(err)

	err = stmt.QueryRow(u.Username, u.Password, u.Email).Scan(&userID)
	checkErr(err)

	token := createToken(userID)

	res = regRes{Message: "successful registration", Success: true, Token: token}
	j, err := json.Marshal(res)
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func createToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userID,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	checkErr(err)

	return tokenString
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

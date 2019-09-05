package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexodorico/goserver/models"
	"github.com/alexodorico/goserver/utils"

	_ "github.com/lib/pq"
)

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/protected", verifyAuth(http.HandlerFunc(handleProtected)))
	err := http.ListenAndServe(":9000", nil)
	utils.CheckErr(err)
}

func verifyAuth(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header["Authorization"][0]
		split := strings.Split(header, " ")
		if len(split) != 2 {
			sendJSON(w, response{Message: "Access denied", Success: false, Token: ""})
			return
		}
		token := split[1]
		if token == "" {
			sendJSON(w, response{Message: "Access denied", Success: false, Token: ""})
			return
		}
		id, err := utils.ParseToken(token)
		if err != nil {
			sendJSON(w, response{Message: "Access denied", Success: false, Token: ""})
			return
		}
		exists := utils.CheckDB("id", id)
		if !exists {
			sendJSON(w, response{Message: "Access denied", Success: false, Token: ""})
			return
		}
		h.ServeHTTP(w, r)
	})
}

func handleProtected(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, response{Message: "Success", Success: true, Token: ""})
}

// handleLogin decodes JSON sent in the body,
// then checks by email to see if the user exists.
// If the email exists, it compares the given password to the
// hashed password stored in the database.
// If the passwords match, a token containing the users id is sent in the response
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var u user
	var hashpw string
	var uid int

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	utils.CheckErr(err)
	exists := utils.CheckDB("email", u.Email)
	if !exists {
		sendJSON(w, response{Message: "Incorrect email or password", Success: false, Token: ""})
		return
	}

	err = models.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", u.Email).Scan(&uid, &hashpw)
	valid := utils.ComparePasswords(hashpw, u.Password)
	if !valid {
		sendJSON(w, response{Message: "Incorrect email or password", Success: false, Token: ""})
		return
	}
	token := utils.CreateToken(uid)
	sendJSON(w, response{Message: "Login successful", Success: true, Token: token})
	return
}

// handleRegister check to see if a user already exists in the database.
// If they don't, inserts hashed password and email into users table
// and sends a JWT containing the user's id in the response
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var sStmt = "INSERT INTO users(password,email) VALUES($1,$2) RETURNING id"
	var u user
	var userID int

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	utils.CheckErr(err)

	exists := utils.CheckDB("email", u.Email)
	if exists {
		sendJSON(w, response{Message: "User already exists", Success: false, Token: ""})
		return
	}

	stmt, err := models.DB.Prepare(sStmt)
	utils.CheckErr(err)
	hash := utils.HashAndSalt(u.Password)
	err = stmt.QueryRow(hash, u.Email).Scan(&userID)
	utils.CheckErr(err)
	token := utils.CreateToken(userID)
	sendJSON(w, response{Message: "Registration successful", Success: true, Token: token})
}

// sendJSON sends a JSON response to client
func sendJSON(w http.ResponseWriter, res response) {
	j, err := json.Marshal(res)
	utils.CheckErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

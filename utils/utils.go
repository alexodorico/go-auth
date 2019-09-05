package utils

import (
	"database/sql"
	"fmt"
	"github.com/alexodorico/goserver/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt generates a hashed password from a plain string
func HashAndSalt(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	CheckErr(err)
	return string(hash)
}

// ComparePasswords takes two strings, one hashed and one plain,
// and checked to see if they are equal
func ComparePasswords(hashed string, plain string) bool {
	hashedByte := []byte(hashed)
	plainByte := []byte(plain)
	err := bcrypt.CompareHashAndPassword(hashedByte, plainByte)
	if err != nil {
		return false
	}
	return true
}

// GetID parses the JWT and returns the user ID as a string
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", fmt.Errorf("Invalid token")
	}
	id := token.Claims.(jwt.MapClaims)["id"]
	str := fmt.Sprintf("%v", id)
	return str, nil
}

// CheckDB checks for the existance of a row in the databaase
func CheckDB(field, value string) bool {
	query := fmt.Sprintf("SELECT id FROM users WHERE %s = $1", field)
	err := models.DB.QueryRow(query, value).Scan()
	if err != sql.ErrNoRows {
		return true
	}
	return false
}

// CreateToken creates a JWT that includes the users ID in the body
func CreateToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userID,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	CheckErr(err)
	return tokenString
}

// CheckErr checks to see if there is an error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

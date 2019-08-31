package utils

import (
	"github.com/dgrijalva/jwt-go"
	"database/sql"
	"github.com/alexodorico/goserver/models"
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

// CheckUserExists tries to find if a user is present in the database
func CheckUserExists(email string) bool {
	err := models.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan()
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth", Auth).Methods("PUT")
	r.HandleFunc("/hello", GetUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", r))
}

// User details for authenticating
type User struct {
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auth controller(handler?) to authenticate user with
// username/password and return token
func Auth(w http.ResponseWriter, r *http.Request) {
	admin := User{UUID: "1", Username: "admin", Password: "pass"}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		log.Print(err)
	}

	// Validate user/passwd
	if user.Username == admin.Username && user.Password == admin.Password {
		// Create Token
		fmt.Println("Match!")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

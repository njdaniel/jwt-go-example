package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth", Auth).Methods("PUt")

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

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

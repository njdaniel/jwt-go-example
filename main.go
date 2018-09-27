package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// RSA keys
const (
	privKeyPath = "./app.rsa"
	pubKeyPath  = "./app.rsa.pub"
)

var VerifyKey, SignKey []byte

func init() {
	var err error

	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index)
	r.HandleFunc("/auth", Auth).Methods("PUT")

	log.Fatal(http.ListenAndServe(":5000", nil))
}

// Index is the root of the url tree to test if up
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Inside the index")
	return
}

//------------------ STRUCT DEFINITIONS-------------------

// UserCredentials details for authenticating
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User with a unique id
type User struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

// Auth controller(handler?) to authenticate user with
// username/password and return token
func Auth(w http.ResponseWriter, r *http.Request) {
	admin := User{Username: "admin", Password: "pass"}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		log.Print(err)
	}

	// Validate user/passwd
	if user.Username == admin.Username && user.Password == admin.Password {
		// Start to create Token
		fmt.Println("Match!")
		// create rsa256 signer
		signer := jwt.New(jwt.GetSigningMethod("RS256"))

		// set claims
		claims := make(jwt.MapClaims)
		claims["iss"] = "admin"
		claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
		claims["CustomUserInfo"] = struct {
			Name string
			Role string
		}{user.Username, "Member"}

		tokenString, err := signer.SignedString(SignKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while signing the token")
			log.Printf("Error signing token: %v\n", err)
		}

		//create a token instance using the token string
		response := Token{tokenString}
		JSONResponse(response, w)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

// JSONResponse helper
func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

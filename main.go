package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// RSA keys
const (
	privKeyPath = "app.rsa"
	pubKeyPath  = "app.rsa.pub"
)

// VerifyKey are used to  verify the token
var VerifyKey []byte

// SignKey used to create token
var SignKey []byte

func init() {
	fmt.Println("init")
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
	a := App{}
	a.Run(":5000")
	// r := mux.NewRouter()

	// r.HandleFunc("/", Index)
	// r.HandleFunc("/auth", Auth).Methods("POST")
	// r.HandleFunc("/public", GetPublic).Methods("GET")

	// //PROTECTED ENDPOINTS
	// r.Handle("/resource/", negroni.New(
	// 	negroni.HandlerFunc(ValidateTokenMiddleware),
	// 	negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	// ))

	// fmt.Println("Listening...")
	// http.ListenAndServe(":5000", r)
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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response JSON
type Response struct {
	Data string `json:"data"`
}

// Token is struct to return token in JSON
type Token struct {
	Token string `json:"token"`
}

func GetPublic(w http.ResponseWriter, r *http.Request) {
	response := Response{"Public Access"}
	JSONResponse(response, w)
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
		signer.Claims = claims

		fmt.Println(string(SignKey))
		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(SignKey)
		if err != nil {
			log.Println(err)
		}
		tokenString, err := signer.SignedString(signKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while signing the token")
			log.Printf("Error signing token: %v\n", err)
			return
		}

		//create a token instance using the token string
		response := Token{tokenString}
		JSONResponse(response, w)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

// ProtectedHandler is example of endpoint that needs token to access
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JSONResponse(response, w)
}

//ValidateTokenMiddleware AUTH TOKEN VALIDATION
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(VerifyKey)
	if err != nil {
		log.Println(err)
	}
	//validate token
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
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

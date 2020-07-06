//Handles user autentication, and authroization

package goChess

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"reflect"
	"time"
)

//This seems bad.  Why do I have to export these sensitive fields?
type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//This is temporary.  Set ENV variables!
var store = sessions.NewCookieStore([]byte("topsecret"))

func handleAuth(w http.ResponseWriter, r *http.Request) {
	//Decode data from client, route to register/login
	w.Header().Set("Content-Type", "application/json")
	var regData authData
	json.NewDecoder(r.Body).Decode(&regData)

	if regData.Email != "" {
		regData.register()
	} else {
		regData.login()
	}
}

func (regData authData) register() (user, email bool) {
	//Query DB to ensure unique username/email
	var userExists, emailExists bool
	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);", regData.Username).Scan(&userExists)
	err = dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);", regData.Email).Scan(&emailExists)
	//Create user
	if !(userExists || emailExists) {
		hashpass, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.MinCost)
		createTime := time.Now()
		createTime.Format("01-02-2000")
		res, err := dbconn.Exec(context.Background(),
			"INSERT INTO users(id,username,password,email,createdate) VALUES (uuid_generate_v4(),$1,$2,$3,$4);", regData.Username, hashpass, regData.Email, createTime)
		log.Println(reflect.TypeOf(res))

		if err != nil {
			log.Println(err)
		}
	}

	if err != nil {
		log.Println(err)
	}
	return user, email
}

func (loginData authData) login() (validUser bool, validPass bool, token bool) { //Your gonna have to return cookie/session info from here, or bad user or bad password.
	//Fetch and compare password hashes
	var passHash string

	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1;", loginData.Username).Scan(&passHash)
	log.Println(passHash)

	if passHash == "" {
		return false, false, false
	}

	storedHash := []byte(passHash)
	inputText := []byte(loginData.Password)
	err = bcrypt.CompareHashAndPassword(storedHash, inputText)

	if err == nil {
		//Do login authorization stuff here
		//Get ID from database and store to Player{id:responsevalue}
		return true, true, true
	} else {
		return true, false, false
		//The passwords didn't match, tell the client
	}

}

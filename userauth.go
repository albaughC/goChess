//Handles user autentication, and authroization

package goChess

import (
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//This seems bad.  Why do I have to export these sensitive fields?
//Also can I combine these two structs into one?
//Add Json fields
type registerData struct {
	Username string
	Password string
	Email    string
}

type loginData struct {
	Username string
	Password string
}

func register(w http.ResponseWriter, r *http.Request) {
	//Decode data from client
	w.Header().Set("Content-Type", "application/json")
	var temp registerData
	json.NewDecoder(r.Body).Decode(&temp)

	//Query DB to ensure unique username/email
	var userExists, emailExists bool
	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);", temp.Username).Scan(&userExists)
	err = dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);", temp.Email).Scan(&emailExists)
	//Create user
	if !(userExists || emailExists) {
		hashpass, err := bcrypt.GenerateFromPassword([]byte(temp.Password), bcrypt.MinCost)
		createTime := time.Now()
		createTime.Format("01-02-2000")
		log.Println(createTime)
		res, err := dbconn.Exec(context.Background(),
			"INSERT INTO users(id,username,password,email,createdate) VALUES (uuid_generate_v4(),$1,$2,$3,$4);", temp.Username, hashpass, temp.Email, createTime)

		if err != nil {
			log.Println(err)
		}
		log.Println(res)
	}
	if err != nil {
		log.Println(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//Decode data from client
	w.Header().Set("Content-Type", "application/json")
	var temp loginData
	json.NewDecoder(r.Body).Decode(&temp)
	//Fetch and compare password hashes
	var userExists bool
	var passHash string
	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);", temp.Username).Scan(&userExists)
	err = dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE password= $1);", passHash).Scan(&passHash)
	storedHash := []byte(passHash)
	inputText := []byte(temp.Password)
	err = bcrypt.CompareHashAndPassword(storedHash, inputText)

	if err == nil {
		//Do login authorization stuff here
		//Get ID from database and store to Player{id:responsevalue}
	} else {
		log.Println("Do stuff that says the password didn't exits, username didn't exist or didn't match")
	}
}

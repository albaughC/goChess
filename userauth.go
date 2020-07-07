//Handles user autentication, and authorization

package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//This is temporary.  Set ENV variables!
var store = sessions.NewCookieStore([]byte("topsecret"))
var userMap map[string]*Player

//This seems bad.  Why do I have to export these sensitive fields?
type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	//Decode data from client, route to register/login
	w.Header().Set("Content-Type", "application/json")
	var regData authData
	json.NewDecoder(r.Body).Decode(&regData)
	if regData.Username == "" {
		json.NewEncoder(w).Encode(struct{ Err string }{Err: "Username invalid"})
	} else if regData.Password == "" {
		json.NewEncoder(w).Encode(struct{ Err string }{Err: "Password invalid"})
	} else if regData.Email != "" {
		user, email := regData.register()
		json.NewEncoder(w).Encode(struct {
			User  bool
			Email bool
		}{User: user, Email: email})
		//This logic is backwards.  If they are both false, a new user was created
	} else {
		//Should test whether the user has an active session
		log.Println("your in LogData")
		validUser, validPass := regData.login()
		if validUser && validPass {
			log.Println("Your init session")
			session, _ := store.Get(r, "session-name")
			session.Values["username"] = regData.Username
			session.Values["authenicated"] = true
			session.Save(r, w)
		} else {
			json.NewEncoder(w).Encode(struct {
				User bool
				Pass bool
			}{User: validUser, Pass: validPass})
		}
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

		if err != nil {
			log.Println(err)
			log.Println(res)
		}
	}
	if err != nil {
		log.Println(err)
	}
	return userExists, emailExists
}

func (loginData authData) login() (validUser bool, validPass bool) { //Your gonna have to return cookie/session info from here, or bad user or bad password.
	//Fetch and compare password hashes
	var passHash string

	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1;", loginData.Username).Scan(&passHash)

	if passHash == "" {
		return false, false
	}

	storedHash := []byte(passHash)
	inputText := []byte(loginData.Password)
	err = bcrypt.CompareHashAndPassword(storedHash, inputText)

	if err == nil {
		return true, true
	} else {
		return true, false
	}

}

func logout() {
	log.Println("You can leave any time you want, but you can never logout")
}

func SessionMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		if session.IsNew {
			log.Println("The session is new, there wasn't one before")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		val := session.Values["username"]
		user, ok := val.(string)
		i, ok := userMap[user]
		if !ok {
			log.Println("your intantinating a player")
			log.Println(i)
			//Instantiate Player, store in map
		}
		next.ServeHTTP(w, r)
	})
}
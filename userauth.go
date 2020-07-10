//Handles user autentication, and authorization

package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

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
	} else {
		//Should test whether the user has an active session
		validUser, validPass := regData.login()
		if validUser && validPass {
			session, _ := store.Get(r, "mySession")
			session.Values["username"] = regData.Username
			session.Values["authenicated"] = true
			session.Save(r, w)
			json.NewEncoder(w).Encode(struct{ Login bool }{Login: true})
		} else {
			json.NewEncoder(w).Encode(struct {
				User  bool
				Pass  bool
				Login bool
			}{User: validUser, Pass: validPass, Login: false})
		}
	}
}

func (regData authData) register() (user, email bool) {
	//Query DB to ensure unique username/email
	var userExists, emailExists bool
	dbconn := connectToDb()
	err := dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);", regData.Username).Scan(&userExists)
	err = dbconn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);", regData.Email).Scan(&emailExists)
	//This inverts the queries.  We want to know if the user/email DO NOT exist, so we can create them
	userExists, emailExists = !userExists, !emailExists
	defer dbconn.Close()

	if err != nil {
		log.Println(err)
	}
	//Create user
	if userExists && emailExists {
		hashpass, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.MinCost)
		createTime := time.Now()
		createTime.Format("01-02-2000")
		res, err := dbconn.Exec(context.Background(),
			"INSERT INTO users(id,username,password,email,createdate) VALUES (uuid_generate_v4(),$1,$2,$3,$4);", regData.Username, hashpass, regData.Email, createTime)
		userExists, emailExists = true, true
		if err != nil {
			log.Println(res, err)
		}
	}
	return userExists, emailExists
}

func (loginData authData) login() (validUser bool, validPass bool) {
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
		session, err := store.Get(r, "mySession")
		if err != nil {
			log.Println(err)
		}

		if session.IsNew {
			http.Redirect(w, r, "/public/html/login.html", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

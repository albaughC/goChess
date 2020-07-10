package main

import (
	//	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type chatUsers struct {
	userList []string `json:"User List"`
}

func (list chatUsers) openUserpageSocket(w http.ResponseWriter, r *http.Request) {
	//Establish socket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()
	//Collect session data, or redirect to login
	session, err := store.Get(r, "mySession")
	val := session.Values["username"]
	user, ok := val.(string)
	if !ok {
		log.Println("Session not defined")
		//This should never happen, but...
		//I could throw in an http redirect here to login
	}
	list.userList = append(list.userList, user)
	//Do stuff with the socket....
	//Send list.userList to the client, maintain the socket
}

//Logic
//On userpage.html load, fetch calls api endpoint to activate socket
//Server adds username to list of online users
//Updates online users in client with new user list

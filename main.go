//ChessLife design
//Features:
//Chess clock, with intervals
//Pieces Captured
//AI?
//Move highlight
//Clock inits AFTER you move, so that the first move is not timed.

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Player struct {
	Id       string `json:"playerId"`
	Username string `json:"username"`
	//	Color rune `json:"playerColor"` // 'w' || 'b'
	//	clock    *clock
	//captured [15]rune
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Your in test func")
}

//Set server timeouts
func main() {

	var userList chatUsers

	port := ":" + os.Getenv("PORT")
	log.Println("Bringing up server on port:" + port)

	r := mux.NewRouter()
	publicRoute := r.PathPrefix("/public").Subrouter()
	privateRoute := r.PathPrefix("/private").Subrouter()
	privateRoute.Use(SessionMid)

	publicfs := http.FileServer(http.Dir("./public/html/"))
	privatefs := http.FileServer(http.Dir("./private/html/"))
	publicRoute.PathPrefix("/html").Handler(http.StripPrefix("/public/html/", publicfs))
	privateRoute.PathPrefix("/html").Handler(http.StripPrefix("/private/html/", privatefs))

	publicRoute.HandleFunc("/api/login", handleAuth).Methods("POST")
	publicRoute.HandleFunc("/api/userwebsocket", userList.openUserpageSocket)

	log.Fatal(http.ListenAndServe(port, r))
}

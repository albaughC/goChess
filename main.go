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

func main() {

	var userList chatUsers

	port := os.Getenv("PORT")
	log.Println("Brining up server on port:" + port)
	port = ":" + port

	publicRoute := mux.NewRouter()
	publicRoute.PathPrefix("/public")
	privateRoute := publicRoute.PathPrefix("/private").Subrouter()
	privateRoute.Use(SessionMid)

	publicfs := http.FileServer(http.Dir("./public/html/"))
	privatefs := http.FileServer(http.Dir("./private/html/"))
	publicRoute.PathPrefix("/html").Handler(http.StripPrefix("/html/", publicfs))
	privateRoute.PathPrefix("/html").Handler(http.StripPrefix("/html/", privatefs))

	publicRoute.HandleFunc("/api/login", handleAuth).Methods("POST")
	privateRoute.HandleFunc("/api/userwebsocket", userList.openUserpageSocket)

	log.Fatal(http.ListenAndServe(port, publicRoute))
}

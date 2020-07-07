//ChessLife design
//Features:
//Chess clock, with intervals
//Pieces Captured
//AI?
//Move highlight
//Clock inits AFTER you move, so that the first move is not timed.

package main

import (
	"github.com/gorilla/handlers"
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

func main() {
	port := os.Getenv("PORT")
	port = ":" + port
	route := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	fs := http.FileServer(http.Dir("./html/"))

	route.HandleFunc("/api/login", handleAuth).Methods("POST")
	route.PathPrefix("/html").Handler(http.StripPrefix("/html/", fs))

	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(route)))
}

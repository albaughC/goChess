//ChessLife design
//Features:
//Chess clock, with intervals
//Pieces Captured
//AI?
//Move highlight
//Clock inits AFTER you move, so that the first move is not timed.

package goChess

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

type Player struct {
	//	Id    int  `json:"playerId"`
	Color rune `json:"playerColor"` // 'w' || 'b'
	//	clock    *clock
	//captured [15]rune
}

func connectToDb() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(),
		"postgresql://topher@127.0.0.1:5432/chesslife")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func main() {
	//Dummy player inits, to be implemented later
	player1 := Player{'w'}
	player2 := Player{'b'}
	boardState := Board{BoardState: standardChessInit(&player1, &player2), Id: 1}

	route := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	route.HandleFunc("/api/boardstate/{id}", boardState.getBoardState).Methods("GET")
	route.HandleFunc("/api/login", login).Methods("POST")
	route.HandleFunc("/api/register", register).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(route)))
}

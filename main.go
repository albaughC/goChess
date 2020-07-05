//ChessLife design
//Features:
//Chess clock, with intervals
//Pieces Captured
//AI?
//Move highlight
//Clock inits AFTER you move, so that the first move is not timed.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type Board struct {
	BoardState [8][8]Tile `json:"BoardState"`
	Id         int        `json:"Id"`
	dbconn     *pgxpool.Pool
}

type Tile struct {
	TileColor  string `json:"tileColor"`
	IsOccupied bool   `json:"isOccupied"`
	Piece      *Piece `json:"piece"`
}

type Piece struct {
	Color     rune    `json:"color"`
	PieceType rune    `json:"type"`
	LocationX int     `json:"x-cord"`
	LocationY int     `json:"y-cord"`
	Owner     *Player `json:"owner"`
}

type Player struct {
	//	Id    int  `json:"playerId"`
	Color rune `json:"playerColor"` // 'w' || 'b'
	//	clock    *clock
	//captured [15]rune
}

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

func generateBoard() [8][8]Tile {
	var boardArray [8][8]Tile
	for i := range boardArray {
		for k := range boardArray[i] {
			if (i+k)%2 == 0 {
				boardArray[i][k] = Tile{TileColor: "white", IsOccupied: false}
			} else {
				boardArray[i][k] = Tile{TileColor: "black", IsOccupied: false}
			}
		}
	}
	return boardArray
}

func standardChessInit(playerOne *Player, playerTwo *Player) [8][8]Tile { //Blit pieces, set white to active, start clock
	chessBoard := generateBoard()
	chessPieces := [8]rune{'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r'}

	for i := range chessBoard[0] {
		chessBoard[0][i].Piece = &Piece{Color: 'w', PieceType: chessPieces[i], LocationX: i, LocationY: 0}
		chessBoard[0][i].IsOccupied = true
	}
	for k := range chessBoard[1] {
		chessBoard[1][k].Piece = &Piece{Color: 'w', PieceType: 'p', LocationX: k, LocationY: 1}
		chessBoard[1][k].IsOccupied = true
	}
	for j := range chessBoard[6] {
		chessBoard[6][j].Piece = &Piece{Color: 'b', PieceType: 'p', LocationX: j, LocationY: 6}
		chessBoard[6][j].IsOccupied = true
	}
	for l := range chessBoard[7] {
		chessBoard[7][l].Piece = &Piece{Color: 'b', PieceType: chessPieces[l], LocationX: l, LocationY: 7}
		chessBoard[7][l].IsOccupied = true
	}
	if playerOne.Color == 'w' {
		for m := range chessBoard[0] {
			chessBoard[0][m].Piece.Owner = playerOne
			chessBoard[1][m].Piece.Owner = playerOne
			chessBoard[6][m].Piece.Owner = playerTwo
			chessBoard[7][m].Piece.Owner = playerTwo
		}
	} else {
		for m := range chessBoard[0] {
			chessBoard[0][m].Piece.Owner = playerTwo
			chessBoard[1][m].Piece.Owner = playerTwo
			chessBoard[6][m].Piece.Owner = playerOne
			chessBoard[7][m].Piece.Owner = playerOne
		}
	}
	return chessBoard
}

func (state Board) getBoardState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params)
	json.NewEncoder(w).Encode(state)
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

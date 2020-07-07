//Handles generation, manipulation, and storage of the board state

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
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
	//Set the user ID for black/white in the board state instead of the player state, since its a function of a game
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

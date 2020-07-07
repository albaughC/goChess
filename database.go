package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func connectToDb() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(),
		"postgresql://topher@127.0.0.1:5432/chesslife")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

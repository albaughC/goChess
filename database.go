package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func connectToDb() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(),
		os.Getevn("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

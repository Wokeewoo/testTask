package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
)

func connectToDB() (*pgx.Conn, error) {
	db, err := pgx.Connect(configs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		return nil, err
	}
	return db, nil
}

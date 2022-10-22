package main

import (
	"github.com/jackc/pgx"
)

var configs = pgx.ConnConfig{Host: "localhost", Port: 5432, Database: "postgres", User: "aba", Password: "123321"}

func main() {
	SetCacheFromDB()
	GetChanMsgs()
}

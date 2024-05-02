package db

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

var db *sql.DB
var sqlMigrations embed.FS

func dbConnect() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = connection.Ping()
	if err != nil {
		log.Fatal("Error pinging database")
	}
	defer connection.Close()
	db = connection
	goose.SetBaseFS(sqlMigrations)
	goose.SetDialect("postgres")
	log.Println("Database connected")
	dir := "migrations"
	err = goose.Up(db, dir)
	log.Println(err)
	log.Println("Migrations executed")
}

func Connect() {
	dbConnect()
}

func GetDB() *sql.DB {
	return db
}

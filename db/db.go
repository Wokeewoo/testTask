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
	"testTask/logger"
)

var db *sql.DB
var sqlMigrations embed.FS

func dbConnect() {
	logger := logger.GetLogger()
	e := godotenv.Load()

	if e != nil {
		log.Fatal("Error loading .env file")
	}
	logger.Debug("Loaded .env file")

	username := os.Getenv("db_user")
	logger.WithField("username", username).Debugln("Loaded username from .env file")
	password := os.Getenv("db_pass")
	logger.Debugln("Loaded password from .env file")
	dbName := os.Getenv("db_name")
	logger.WithField("dbName", dbName).Debugln("Loaded dbName from .env file")
	dbHost := os.Getenv("db_host")
	logger.WithField("dbHost", dbHost).Debugln("Loaded dbHost from .env file")
	dbPort := os.Getenv("db_port")
	logger.WithField("dbPort", dbPort).Debugln("Loaded dbPort from .env file")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infoln("Connected to database")
	err = connection.Ping()
	if err != nil {
		logger.Fatal("Error pinging database")
	}
	logger.Debugln("Pinged database")

	db = connection
	err = goose.SetDialect("postgres")
	logger.Infoln("migrations started")
	err = goose.Up(db, "migrations")
	if err != nil {
		logger.WithError(err).Errorln("Error executing migrations")
	}
	logger.Infoln("migrations finished")
}

func Connect() {
	dbConnect()
}

func GetDB() *sql.DB {
	return db
}

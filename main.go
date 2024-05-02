package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"testTask/db"
)

func init() {
	log.Println("Starting server...")
	db.Connect()
	log.Println("Database connected")
}

func main() {
	router := mux.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}

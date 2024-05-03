package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"testTask/controllers"
	"testTask/db"
	"time"
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
		port = "8000"
	}
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	router.HandleFunc("api/cars/{id}", controllers.GetCar).Methods("GET")
	router.HandleFunc("api/cars/", controllers.GetCars).Methods("Get")
	router.HandleFunc("api/cars/", controllers.CreateCar).Methods("POST")
	http.Handle("/", router)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

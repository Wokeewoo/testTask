package main

import (
	"log"
	"net/http"
	"os"
	"testTask/controllers"
	"testTask/db"
)

func main() {
	log.Println("Starting server...")
	db.Connect()
	log.Println("Database connected")
	mux := http.NewServeMux()
	port := os.Getenv("PORT")
	if port == "" {
		port = "545"
	}

	host := os.Getenv("HOST")

	mux.HandleFunc("GET /api/cars/{id}", controllers.GetCar)
	mux.HandleFunc("GET /api/cars/", controllers.GetCars)
	mux.HandleFunc("POST /api/cars/", controllers.CreateCar)
	mux.HandleFunc("DELETE /api/cars/{id}", controllers.DeleteCar)
	mux.HandleFunc("PUT /api/cars/{id}", controllers.UpdateCar)
	err := http.ListenAndServe(host+":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()
}

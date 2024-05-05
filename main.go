package main

import (
	"log"
	"net/http"
	"os"
	"testTask/controllers"
	"testTask/db"
)

//	@title			Test Task
//	@version		1.0
//	@description	This is an service for car catalog.

//	@contact.name	Abdallah Izaripov
//	@contact.email	abazerov@yandex.ru

//	@host		localhost:8000
//	@BasePath	/api

func main() {
	log.Println("Starting server...")
	db.Connect()
	log.Println("Database connected")
	mux := http.NewServeMux()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	host := os.Getenv("HOST")

	mux.HandleFunc("GET /api/cars/{id}", controllers.GetCar)
	mux.HandleFunc("GET /api/cars/", controllers.GetCars)
	mux.HandleFunc("POST /api/cars/", controllers.CreateCar)
	mux.HandleFunc("DELETE /api/cars/{id}", controllers.DeleteCar)
	mux.HandleFunc("PATCH /api/cars/{id}", controllers.UpdateCar)
	err := http.ListenAndServe(host+":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()
}

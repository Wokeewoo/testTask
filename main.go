package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"testTask/controllers"
	"testTask/db"
	"testTask/logger"
)

//	@title			Test Task
//	@version		1.0
//	@description	This is an service for car catalog.

//	@contact.name	Abdallah Izaripov
//	@contact.email	abazerov@yandex.ru

//	@host		localhost:8000
//	@BasePath	/api

func main() {
	logger := logger.GetLogger()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	db.Connect()
	mux := http.NewServeMux()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	logger.WithField("port", port).Debugln("getting port")
	host := os.Getenv("HOST")
	logger.WithField("host", host).Debugln("getting host")
	mux.HandleFunc("GET /api/cars/{id}", controllers.GetCar)
	mux.HandleFunc("GET /api/cars/", controllers.GetCars)
	mux.HandleFunc("POST /api/cars/", controllers.CreateCar)
	mux.HandleFunc("DELETE /api/cars/{id}", controllers.DeleteCar)
	mux.HandleFunc("PATCH /api/cars/{id}", controllers.UpdateCar)
	logger.Infoln("Starting server...")
	err := http.ListenAndServe(host+":"+port, mux)

	if err != nil {
		logger.WithError(err).Errorln("Failed to start server")
	}

	defer db.GetDB().Close()
	defer logger.Infoln("Server stopped")
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testTask/models"
)

func GetCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("id is not a number")
		return
	}
	car, err := models.GetCar(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(car)
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response models.CreateCarResponse
	response.Cars = make([]models.Car, 0)
	response.Errors = make([]string, 0)
	var req models.CreateCarRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	url := os.Getenv("car_info_api_url")
	car := models.Car{}

	extReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	rq := extReq.URL.Query()

	for i := range req.RegNums {
		regNum := req.RegNums[i]
		rq.Set("regNum", regNum)
		extReq.URL.RawQuery = rq.Encode()
		resp, err := http.DefaultClient.Do(extReq)

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error getting car info.\n status code: %d \n status message: %s", resp.StatusCode, resp.Status)
			response.Errors = append(response.Errors, fmt.Sprintf("Error getting car info for regNum: %s \n status code: %d \n status message: %s", regNum, resp.StatusCode, resp.Status))
			continue
		}
		data, err := io.ReadAll(resp.Body)
		err = json.Unmarshal(data, &car)
		err = car.ValidateCar()
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Error getting car info for regNum: %s \n error: %s", regNum, err.Error()))
			continue
		}
		err = models.CreateCar(&car)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Error when creating car for regNum: %s \n error: %s", regNum, err.Error()))
			continue
		}
		response.Cars = append(response.Cars, car)
	}

	json.NewEncoder(w).Encode(response)
}

func GetCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	values := r.URL.Query()
	year, _ := strconv.Atoi(values.Get("year"))
	page, _ := strconv.Atoi(values.Get("page"))
	limit, _ := strconv.Atoi(values.Get("limit"))
	filter := models.CarFilter{
		RegNum:          values.Get("regNum"),
		Mark:            values.Get("mark"),
		Model:           values.Get("model"),
		Year:            year,
		OwnerName:       values.Get("owner_name"),
		OwnerSurname:    values.Get("owner_surname"),
		OwnerPatronymic: values.Get("owner_patronymic"),
		Limit:           limit,
		Page:            page,
	}
	carsList, err := models.GetCars(&filter)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	response := models.Cars{Cars: carsList}
	json.NewEncoder(w).Encode(response)
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("id is not a number")
		return
	}

	err = models.DeleteCar(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode("deleted")

}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("id is not a number")
		return
	}
	var req models.UpdateCarRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	err = models.UpdateCar(id, &req)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode("updated")
}

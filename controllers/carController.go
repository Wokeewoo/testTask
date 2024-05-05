package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	logg "testTask/logger"
	"testTask/models"
)

var logger = logg.GetLogger()

type createCarRequest struct {
	RegNums []string `json:"regNums"`
}

type updateCarRequest struct {
	Id  int                     `json:"id"`
	Car models.UpdateCarRequest `json:"car"`
}

func ValidateRequestBody(r *http.Request, req interface{}) error {

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		return err
	}
	return nil
}

// GetCar godoc
//
//	@Summary		Get a car
//	@Description	get a car by ID
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Car ID"
//	@Success		200	{object}	models.Car
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/cars/{id} [get]
func GetCar(w http.ResponseWriter, r *http.Request) {
	logger.Infoln("get request on get /cars/{id}")
	id, err := strconv.Atoi(r.PathValue("id"))
	logger.WithField("id", id).Debugln("getting id")
	if err != nil {
		logger.WithError(err).Errorln("Did not get id, bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	logger.Debugln("Set Content-Type to application/json")

	logger.Debugln("go to models.GetCar")
	car, err := models.GetCar(id)
	if err != nil {
		logger.WithError(err).Errorln("Error getting car")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(car)
	logger.Debugln("put result in json")
}

// CreateCar godoc
//
//	@Summary		Create cars
//	@Description	create new cars by list of regNums
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			regNums		body		models.CreateCarRequest	true	"List of regNums"
//	@Success		200	{object}	models.CreateCarResponse
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/cars [post]
func CreateCar(w http.ResponseWriter, r *http.Request) {
	err := ValidateRequest(r, &createCarRequest{})
	if err != nil {
		logger.WithError(err).Errorln("Bad request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}
	logger.Infoln("get request on post /cars")
	w.Header().Set("Content-Type", "application/json")
	logger.Debugln("Set Content-Type to application/json")
	var response models.CreateCarResponse
	response.Cars = make([]models.Car, 0)
	response.Errors = make([]string, 0)
	var req models.CreateCarRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	logger.Debugln("decode request")
	if err != nil {
		logger.WithError(err).Errorln("Error decoding request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := os.Getenv("car_info_api_url")
	logger.Debugln("get url from env")
	car := models.Car{}

	extReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.WithError(err).Errorln("Did not get car info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rq := extReq.URL.Query()

	for i := range req.RegNums {
		regNum := req.RegNums[i]
		rq.Set("regNum", regNum)
		extReq.URL.RawQuery = rq.Encode()
		resp, err := http.DefaultClient.Do(extReq)
		logger.Debugln("Do request in external api")
		defer resp.Body.Close()

		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Error getting car info for regNum: %s \n error: %s", regNum, err.Error()))
			continue
		}

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
			logger.Errorln(fmt.Sprintf("Error getting car info for regNum: %s \n error: %s", regNum, err.Error()))
			continue
		}
		err = models.CreateCar(&car)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Error when creating car for regNum: %s \n error: %s", regNum, err.Error()))
			logger.Errorln(fmt.Sprintf("Error when creating car for regNum: %s \n error: %s", regNum, err.Error()))
			continue
		}
		response.Cars = append(response.Cars, car)
		logger.Debugln("car for regNum: %s is created", regNum)
	}

	json.NewEncoder(w).Encode(response)
	logger.Debugln("put result in json")
}

// GetCars godoc
//
//	@Summary		get list of cars
//	@Description	get list of cars with optional filters and pagination
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			limit 	query		int	false	"Limit"
//	@Param			page 	query		int	false	"Page"
//	@Param			regNum 	query		string	false	"RegNum"
//	@Param			mark 	query		string	false	"Mark"
//	@Param			model 	query		string	false	"Model"
//	@Param			year 	query		int	false	"Year"
//	@Param			owner_name 	query		string	false	"OwnerName"
//	@Param			owner_surname 	query		string	false	"OwnerSurname"
//	@Param			owner_patronymic 	query		string	false	"OwnerPatronymic"
//	@Success		200	{object}	models.Cars
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/cars [get]
func GetCars(w http.ResponseWriter, r *http.Request) {

	logger.Infoln("get request on get /cars")
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
	logger.Debugln("Get filter and pagination parameters")
	w.Header().Set("Content-Type", "application/json")
	logger.Debugln("Set Content-Type to application/json")
	carsList, err := models.GetCars(&filter)
	logger.Debugln("Get list of cars")
	if err != nil {
		logger.WithError(err).Errorln("Error getting cars list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := models.Cars{Cars: carsList}
	json.NewEncoder(w).Encode(response)
	logger.Debugln("put result in json")
}

// DeleteCar godoc
// @Summary		delete car
// @Description	delete car by id
// @Tags			cars
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Car ID"
// @Success		200	{string}	string	"deleted"
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/cars/{id} [delete]
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	logger.WithField("id", r.PathValue("id")).Infoln("get request on delete /cars/{id}")
	w.Header().Set("Content-Type", "text/plain")
	logger.Debugln("Set Content-Type to text/plain")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Errorln("id is not a number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Debugln("Get car id")

	err = models.DeleteCar(id)
	logger.Debugln("Delete car")
	if err != nil {
		logger.WithError(err).Errorln("Error deleting car")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("deleted"))
	logger.Debugln("put result in response")

}

// UpdateCar godoc
// @Summary		update car
// @Description	update car by id
// @Tags			cars
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Car ID"
// @Param			car	body		models.UpdateCarRequest	true	"Car"
// @Success		200	{string}	string	"updated"
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Failure		500	{object}	string
// @Router			/cars/{id} [patch]
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	err := ValidateRequest(r, &updateCarRequest{})
	if err != nil {
		logger.WithError(err).Errorln("Bad request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}
	logger.WithField("id", r.PathValue("id")).Infoln("get request on patch /cars/{id}")
	w.Header().Set("Content-Type", "text/plain")
	logger.Debugln("Set Content-Type to text/plain")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		logger.Errorln("id is not a number, bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Debugln("Get car id")
	var req models.UpdateCarRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	logger.Debugln("Decode request")
	if err != nil {
		logger.WithError(err).Errorln("Error decoding request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = models.UpdateCar(id, &req)
	logger.Debugln("Update car")
	if err != nil {
		logger.WithError(err).Errorln("Error updating car")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("updated"))
	logger.Debugln("put result in response")
}

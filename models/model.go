package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"testTask/db"
	logg "testTask/logger"
	"time"
)

var logger = logg.GetLogger()

// People entity for db
// Account example
type People struct {
	Name       string `db:"owner_name" json:"name"`
	Surname    string `db:"owner_surname" json:"surname"`
	Patronymic string `db:"owner_patronymic,omitempty" json:"patronymic,omitempty"`
}

// CarFilter entity for db
// CarFilter example
type CarFilter struct {
	RegNum          string `json:"regNum,omitempty"`
	Mark            string `json:"mark,omitempty"`
	Model           string `json:"model,omitempty"`
	Year            int    `json:"year,omitempty"`
	OwnerName       string `json:"owner_name,omitempty"`
	OwnerSurname    string `json:"owner_surname,omitempty"`
	OwnerPatronymic string `json:"owner_patronymic,omitempty"`
	Limit           int    `json:"limit,omitempty"`
	Page            int    `json:"page,omitempty"`
}

// Cars example
type Cars struct {
	Cars []Car
}

// Car example
type Car struct {
	Id     int    `db:"id" json:"id"`
	RegNum string `db:"regNum" json:"regNum"`
	Mark   string `db:"mark" json:"mark"`
	Model  string `db:"model" json:"model"`
	Year   int    `db:"year,omitempty" json:"year,omitempty"`
	Owner  People `db:"owner" json:"owner"`
}

// CreateCarRequest example
type CreateCarRequest struct {
	RegNums []string `json:"regNums"`
}

// CreateCarResponse example
type CreateCarResponse struct {
	Cars   []Car    `json:"cars"`
	Errors []string `json:"errors"`
}

// UpdateCarRequest example
type UpdateCarRequest struct {
	RegNum string `json:"regNum,omitempty"`
	Mark   string `json:"mark,omitempty"`
	Model  string `json:"model,omitempty"`
	Year   int    `json:"year,omitempty"`
	Owner  struct {
		Name       string `json:"name,omitempty"`
		Surname    string `json:"surname,omitempty"`
		Patronymic string `json:"patronymic,omitempty"`
	} `json:"owner,omitempty"`
}

// ValidateCar example
func (c *Car) ValidateCar() error {
	if c.RegNum == "" {
		return fmt.Errorf("regNum is empty")
	}

	if c.Mark == "" {
		return fmt.Errorf("mark is empty")
	}

	if c.Model == "" {
		return fmt.Errorf("model is empty")
	}

	if c.Owner.Name == "" {
		return fmt.Errorf("owner name is empty")
	}

	if c.Owner.Surname == "" {
		return fmt.Errorf("owner surname is empty")
	}
	return nil
}

// GetCar example
func GetCar(id int) (Car, error) {
	car := Car{}
	var updatedDate pgtype.Timestamp
	var createdDate pgtype.Timestamp
	err := db.GetDB().QueryRow("SELECT * FROM Cars WHERE id = $1", id).Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic, &createdDate, &updatedDate)
	if err != nil {
		logger.Error(err)
		return car, errors.New("car not found")
	}
	return car, nil
}

// CreateCar example
func CreateCar(car *Car) error {
	row := db.GetDB().QueryRow("SELECT * FROM Cars WHERE regNum = $1", car.RegNum)
	logger.Debugln("getting row from db)")
	var updatedDate pgtype.Timestamp
	var createdDate pgtype.Timestamp
	var currentTime pgtype.Timestamp
	currentTime.Time = time.Now()
	err := row.Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic, &createdDate, &updatedDate)
	logger.Debugln("scanning row from db)")
	if createdDate.Time != currentTime.Time && err == nil {
		_, err := db.GetDB().Exec("UPDATE Cars SET (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) = ($2, $3, $4, $5, $6, $7, $8) WHERE id = $1",
			car.Id, car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic, createdDate, updatedDate)
		logger.Debugln("updating row in db)")
		if err != nil {
			logger.WithError(err).Error("failed to update car")
			return err
		}
	} else {
		logger.Debugln("inserting row in db)")
		_, err = db.GetDB().Exec("INSERT INTO Cars (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)

		if err != nil {
			logger.WithError(err).Error("failed to insert car")
			return err
		}
	}
	return nil
}

// GetCars example
func GetCars(filter *CarFilter) ([]Car, error) {
	var cars []Car
	logger.WithField("filter", filter).Debugln("filter: ")
	query := "SELECT * FROM Cars"
	logger.Debugln("getting rows from db)")
	if filter != nil {
		query += " WHERE 1=1"
	}
	if filter.RegNum != "" {
		query += fmt.Sprintf(" AND regnum = '%s'", filter.RegNum)
	}
	if filter.Mark != "" {
		query += fmt.Sprintf(" AND mark = '%s'", filter.Mark)
	}
	if filter.Model != "" {
		query += fmt.Sprintf(" AND model = '%s'", filter.Model)
	}
	if filter.Year != 0 {
		query += fmt.Sprintf(" AND year = '%d'", filter.Year)
	}
	if filter.OwnerName != "" {
		query += fmt.Sprintf(" AND owner_name = '%s'", filter.OwnerName)
	}
	if filter.OwnerSurname != "" {
		query += fmt.Sprintf(" AND owner_surname = '%s'", filter.OwnerSurname)
	}
	if filter.OwnerPatronymic != "" {
		query += fmt.Sprintf(" AND owner_patronymic = '%s'", filter.OwnerPatronymic)
	}
	query += " ORDER BY id"
	if filter.Page != 0 {
		query += fmt.Sprintf(" OFFSET %d", (filter.Page-1)*filter.Limit)
	}
	if filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}
	logger.WithField("query", query).Debugln("make query: ")
	rows, err := db.GetDB().Query(query)
	if err != nil {
		logger.WithError(err).Error("failed to get cars from db")
		return cars, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		car := Car{}
		var updatedDate pgtype.Timestamp
		var createdDate pgtype.Timestamp
		err := rows.Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic, &createdDate, &updatedDate)
		logger.WithField("car id", car.Id).Debugln("get car from db)")
		if err != nil {
			logger.WithError(err).Error("failed to get car from db")
			return cars, err
		}
		cars = append(cars, car)
	}
	logger.Debugln("returning cars list")
	return cars, err
}

// DeleteCar example
func DeleteCar(id int) error {
	_, err := db.GetDB().Exec("DELETE FROM Cars WHERE id = $1", id)
	if err != nil {
		logger.WithError(err).Error("failed to delete car from db")
		return err
	}
	logger.Debugln("delete car from db")
	return nil
}

// UpdateCar example
func UpdateCar(id int, car *UpdateCarRequest) error {
	if car.RegNum == "" && car.Mark == "" && car.Model == "" && car.Year == 0 && car.Owner.Name == "" && car.Owner.Surname == "" && car.Owner.Patronymic == "" {
		logger.Debugln("no data to update")
		return errors.New("no data to update")
	}
	logger.Debugln("making update query...")
	query := "UPDATE Cars SET "
	if car.RegNum != "" {
		query += fmt.Sprintf("regnum = '%s', ", car.RegNum)
	}
	if car.Mark != "" {
		query += fmt.Sprintf("mark = '%s', ", car.Mark)
	}
	if car.Model != "" {
		query += fmt.Sprintf("model = '%s', ", car.Model)
	}
	if car.Year != 0 {
		query += fmt.Sprintf("year = '%d', ", car.Year)
	}
	if car.Owner.Name != "" {
		query += fmt.Sprintf("owner_name = '%s', ", car.Owner.Name)
	}
	if car.Owner.Surname != "" {
		query += fmt.Sprintf("owner_surname = '%s', ", car.Owner.Surname)
	}
	if car.Owner.Patronymic != "" {
		query += fmt.Sprintf("owner_patronymic = '%s', ", car.Owner.Patronymic)
	}
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = %d", id)
	logger.WithField("query", query).Debugln("make query: ")
	_, err := db.GetDB().Exec(query)
	if err != nil {
		logger.WithError(err).Error("failed to update car in db")
		return err
	}
	logger.Debugln("update car in db")
	return nil
}

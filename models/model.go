package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/pgtype"
	"testTask/db"
	"time"
)

type People struct {
	Name       string `db:"owner_name" json:"name"`
	Surname    string `db:"owner_surname" json:"surname"`
	Patronymic string `db:"owner_patronymic,omitempty" json:"patronymic,omitempty"`
}

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

type Cars struct {
	Cars []Car
}
type Car struct {
	Id     int    `db:"id" json:"id"`
	RegNum string `db:"regNum" json:"regNum"`
	Mark   string `db:"mark" json:"mark"`
	Model  string `db:"model" json:"model"`
	Year   int    `db:"year,omitempty" json:"year,omitempty"`
	Owner  People `db:"owner" json:"owner"`
}

type CreateCarRequest struct {
	RegNums []string `json:"regNums"`
}

type CreateCarResponse struct {
	Cars   []Car    `json:"cars"`
	Errors []string `json:"errors"`
}

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

func GetCar(id int) (Car, error) {
	car := Car{}
	var updatedDate pgtype.Timestamp
	var createdDate pgtype.Timestamp
	err := db.GetDB().QueryRow("SELECT * FROM Cars WHERE id = $1", id).Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic, &createdDate, &updatedDate)
	return car, err
}

func CreateCar(car *Car) error {
	row := db.GetDB().QueryRow("SELECT * FROM Cars WHERE regNum = $1", car.RegNum)
	var updatedDate pgtype.Timestamp
	var createdDate pgtype.Timestamp
	var currentTime pgtype.Timestamp
	currentTime.Time = time.Now()
	err := row.Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic, &createdDate, &updatedDate)

	if createdDate.Time != currentTime.Time && err == nil {
		_, err := db.GetDB().Exec("UPDATE Cars SET (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) = ($2, $3, $4, $5, $6, $7, $8) WHERE id = $1",
			car.Id, car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic, createdDate, updatedDate)
		if err != nil {
			return err
		}
	} else {
		_, err = db.GetDB().Exec("INSERT INTO Cars (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)

		if err != nil {
			return err
		}
	}
	return nil
}

func GetCars(filter *CarFilter) ([]Car, error) {
	var cars []Car
	query := "SELECT * FROM Cars"
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
	rows, err := db.GetDB().Query(query)
	if err != nil {
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
		if err != nil {
			return cars, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func DeleteCar(id int) error {
	_, err := db.GetDB().Exec("DELETE FROM Cars WHERE id = $1", id)
	return err
}

func UpdateCar(id int, car *UpdateCarRequest) error {
	if car.RegNum == "" && car.Mark == "" && car.Model == "" && car.Year == 0 && car.Owner.Name == "" && car.Owner.Surname == "" && car.Owner.Patronymic == "" {
		return errors.New("no data to update")
	}
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

	_, err := db.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}

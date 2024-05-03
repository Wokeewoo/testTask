package models

import (
	"fmt"
	"testTask/db"
)

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Car struct {
	id     int    `json:"id"`
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year,omitempty"`
	Owner  People `json:"owner"`
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

	err := db.GetDB().QueryRow("SELECT * FROM Cars WHERE id = $1", id).Scan(&car.id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)

	return car, err
}

func CreateCar(car *Car) error {
	_, err := db.GetDB().Exec("INSERT INTO Cars (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)

	if err != nil {
		return err
	}

	return nil
}

func GetCars() ([]Car, error) {
	cars := []Car{}
	rows, err := db.GetDB().Query("SELECT * FROM Cars")
	if err != nil {
		return cars, err
	}
	defer rows.Close()
	for rows.Next() {
		car := Car{}
		err := rows.Scan(&car.id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
		if err != nil {
			return cars, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

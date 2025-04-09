package controllers

import (
	"database/sql"
	"fmt"

	"github.com/crhawkins/db-sandbox/internal/models"
)

func SelectCar(db *sql.DB, id int) (models.Car, error) {
	var car models.Car
	var companyID, colorID int

	// 1. Query the car record
	query := `SELECT id, model, company_id, color_id FROM car WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&car.ID, &car.Model, &companyID, &colorID)
	if err != nil {
		return car, fmt.Errorf("SelectCar: failed to query car: %w", err)
	}

	// 2. Query related company
	err = db.QueryRow(`SELECT id, name FROM company WHERE id = $1`, companyID).
		Scan(&car.Company.ID, &car.Company.Name)
	if err != nil {
		return car, fmt.Errorf("SelectCar: failed to load company: %w", err)
	}

	// 3. Query related color
	err = db.QueryRow(`SELECT id, name FROM color WHERE id = $1`, colorID).
		Scan(&car.Color.ID, &car.Color.Name)
	if err != nil {
		return car, fmt.Errorf("SelectCar: failed to load color: %w", err)
	}

	return car, nil
}

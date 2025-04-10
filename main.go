package main

import (
	"fmt"
	"log"

	dbhelper "github.com/crhawkins/db-sandbox/internal/db-helper"
	"github.com/crhawkins/db-sandbox/internal/models"
	_ "github.com/lib/pq"
)

const (
	host     = "10.20.30.5"
	port     = 0
	dbName   = "foo"
	username = "su"
	password = "Pa$$w0rd"
)

func main() {

	dbhelper.CreatePostgreSQL(host, port, dbName, username, password)

	db, err := dbhelper.ConnectPostgreSQL(host, port, dbName, username, password)
	if err != nil {
		log.Fatal(err)
	}

	db.Delete("foo")

	defer db.Close()

	if version, err := db.Version(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(version)
	}

	if err := db.CreateTable(models.Country{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	if err := db.CreateTable(models.Color{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	if err := db.CreateTable(models.Company{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	if err := db.CreateTable(models.Car{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	if err := db.CreateTable(models.Stock{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	if err := db.CreateTable(models.Dealership{}); err != nil {
		log.Fatal("failed to create table:", err)
	}

	canada := models.Country{Name: "Canada"}
	usa := models.Country{Name: "USA"}

	ID, _ := db.Insert(canada)
	canada.ID = ID

	ID, _ = db.Insert(usa)
	usa.ID = ID

	ford := models.Company{Name: "FORD", Country: canada}
	dodge := models.Company{Name: "Dodge", Country: usa}

	ID, _ = db.Insert(ford)
	ford.ID = ID
	ID, _ = db.Insert(dodge)
	dodge.ID = ID

	red := models.Color{Name: "Red"}
	green := models.Color{Name: "Green"}
	blue := models.Color{Name: "Blue"}

	ID, _ = db.Insert(red)
	red.ID = ID
	ID, _ = db.Insert(green)
	green.ID = ID
	ID, err = db.Insert(blue)
	if err != nil {
		log.Fatal(err)
	}
	blue.ID = ID

	f150 := models.Car{
		Model:   "F-150",
		Company: ford,
		Color:   green,
	}
	vipor := models.Car{
		Model:   "Viporasdfasdfasdfasdfasdfasdfasdfasdfasdfafds",
		Company: dodge,
		Color:   green,
	}
	ID, _ = db.Insert(f150)
	f150.ID = ID
	ID, _ = db.Insert(vipor)
	vipor.ID = ID

	stock1 := models.Stock{Count: 3, Price: 12345.6, Car: f150}
	stock2 := models.Stock{Count: 7, Price: 32345.6, Car: vipor}
	ID, _ = db.Insert(stock1)
	stock1.ID = ID
	ID, _ = db.Insert(stock2)
	stock2.ID = ID

	dealer1 := models.Dealership{Name: "Riverside FORD", Stock: stock1}
	dealer2 := models.Dealership{Name: "Midtown Dodge", Stock: stock2}

	ID, err = db.Insert(dealer1)
	if err != nil {
		log.Fatal(err)
	}
	ID, _ = db.Insert(dealer2)

	sql := dbhelper.CreateSelectSQL(
		"Dealership",
		[]string{
			"ID", "Stock.Count", "Stock.Price",
			"Car.Model",
			"Company.Name", "Dealership.Stock.Company.Country.Name", "Color.Name",
		},
		[]string{"Price>1.0"},
		[]string{"Stock.Car.Company.Country", "Stock.Car.Color"},
	)

	fmt.Println(sql)

	raw, _ := db.Raw(false)

	rows, err := raw.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		// Create a slice of empty interfaces to hold each column
		cols := make([]interface{}, len(columns))
		colPointers := make([]interface{}, len(columns))
		for i := range cols {
			colPointers[i] = &cols[i]
		}

		// Scan the row
		if err := rows.Scan(colPointers...); err != nil {
			log.Fatal(err)
		}

		// Print row
		for i, col := range cols {
			fmt.Printf("%s: %v\t", columns[i], col)
		}
		fmt.Println()
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

}

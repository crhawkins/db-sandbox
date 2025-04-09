package main

import (
	"fmt"
	"log"

	dbhelper "github.com/crhawkins/db-sandbox/internal/db-helper"
	_ "github.com/lib/pq"
)

func main() {
	db, err := dbhelper.NewPostgreSQL("10.20.30.5", 0, "su", "Pa$$w0rd")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if version, err := db.Version(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(version)
	}

	if err := db.CreateIfNotExists("foo"); err != nil {
		log.Fatal(err)
	}

}

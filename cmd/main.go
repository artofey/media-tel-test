package main

import (
	"fmt"
	"log"

	"github.com/artofey/media-tel-test/internal/app"
	"github.com/artofey/media-tel-test/internal/db/postgres"
	"github.com/artofey/media-tel-test/internal/file/csv"
)

var (
	// flags
	inputFile, dbIP, dbPort, dbName, dbLogin, dbPasswd string
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	csvFile, err := csv.New(inputFile)
	if err != nil {
		return fmt.Errorf("NewCSVFile error: %w", err)
	}
	log.Printf("File %v read successfully\n", inputFile)

	db, err := postgres.New(dbIP, dbPort, dbName, dbLogin, dbPasswd)
	if err != nil {
		return fmt.Errorf("database Error: %w", err)
	}
	log.Printf("Database %v initialized\n", dbName)

	app, err := app.NewApplication(db, csvFile)
	if err != nil {
		return fmt.Errorf("application initialized Error: %w", err)
	}
	log.Println("Application initialized")

	if err := app.Do(); err != nil {
		return fmt.Errorf("application Error: %w", err)
	}
	log.Println("Application Done")
	return nil
}

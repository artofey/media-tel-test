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

/*
- передать все необходимые параметры
- создать валидный reader csv
*/
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// прочитать и проанализировать заданный файл (задается имя файла в локальной файловой системе)
	csvFile, err := csv.New(inputFile)
	if err != nil {
		return fmt.Errorf("NewCSVFile error: %v", err)
	}
	log.Println("File readed.")
	// fmt.Printf("%v", csvFile)
	// подключиться к заданной БД (задается IP адрес, порт, имя БД, логин и пароль пользователя)
	db, err := postgres.New(dbIP, dbPort, dbName, dbLogin, dbPasswd)
	if err != nil {
		return fmt.Errorf("postgres database Error: %v", err)
	}
	log.Println("Database inited.")
	app, err := app.NewApplication(db, csvFile)
	if err != nil {
		return fmt.Errorf("application Error: %v", err)
	}
	log.Println("Application inited.")

	if err := app.Do(); err != nil {
		return fmt.Errorf("app Error: %v", err)
	}
	log.Println("Application Done.")
	return nil
}

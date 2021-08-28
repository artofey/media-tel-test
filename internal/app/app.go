package app

import (
	"fmt"

	"github.com/artofey/media-tel-test/internal/db"
	"github.com/artofey/media-tel-test/internal/file"
)

type Application struct {
	rr file.InputFile
	db db.DB
}

func NewApplication(db db.DB, rr file.InputFile) (*Application, error) {
	return &Application{
		rr: rr,
		db: db,
	}, nil
}

/*
создать в БД таблицу с количеством столбцов, равным количеству столбцов в обрабатываемом файле (имя таблицы соответствует имени файла [имя файла задается в ascii без пробелов]). Имена столбцов таблицы в БД имеют формат col1,col2,col3,... .
загрузить данные из файла в таблицу (определение типа данных не требуется, все данные можно преобразовать к типу string при загрузке)
если в CSV файле был определен заголовок, реальные имена столбцов должны сохраниться в БД в специальную таблицу формата

- создать таблицу в базе данных
- записать в таблицу все значения из csv
*/
func (a *Application) Do() error {
	// fmt.Println(a.rr)
	r, err := a.rr.GetRecords()
	// fmt.Println(r.R)
	if err != nil {
		return fmt.Errorf("get records err: %v", err)
	}

	if err := a.db.CreateTable(a.rr.FileName(), *r); err != nil {
		return fmt.Errorf("create table err: %v", err)
	}

	if err := a.db.InsertRecordsToTable(a.rr.FileName(), *r); err != nil {
		return fmt.Errorf("insert records to table err: %s", err)
	}

	return nil
}

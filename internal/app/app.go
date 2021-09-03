package app

import (
	"fmt"

	"github.com/artofey/media-tel-test/internal/db"
	"github.com/artofey/media-tel-test/internal/file"
)

type Application struct {
	f  file.InputFile
	db db.DB
}

func NewApplication(db db.DB, f file.InputFile) (*Application, error) {
	return &Application{
		f:  f,
		db: db,
	}, nil
}

func (a *Application) Do() error {
	r, err := a.f.GetRecords()
	if err != nil {
		return fmt.Errorf("get records err: %w", err)
	}

	if err := a.db.CreateTable(a.f.FileName(), *r); err != nil {
		return fmt.Errorf("create table err: %w", err)
	}

	if err := a.db.InsertRecordsToTable(a.f.FileName(), *r); err != nil {
		return fmt.Errorf("insert records to table err: %w", err)
	}

	return nil
}

package db

import rr "github.com/artofey/media-tel-test"

type DB interface {
	CreateTable(string, rr.Records) error
	InsertRecordsToTable(string, rr.Records) error
}
